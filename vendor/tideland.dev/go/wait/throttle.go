// Tideland Go Wait
//
// Copyright (C) 2019-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait // import "tideland.dev/go/wait"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

const (
	InfiniteLimit    = math.MaxFloat64
	InfiniteDuration = time.Duration(math.MaxInt64)
)

//--------------------
// LIMIT
//--------------------

// Limit defines the number of allowed job starts per second by the throttle.
type Limit float64

// durationForTokens calculates the duration for a given number of tokens in
// a throttle with limit tokens per second.
func (limit Limit) tokensToDuration(tokens float64) time.Duration {
	if limit <= 0 {
		return InfiniteDuration
	}
	seconds := tokens / float64(limit)
	return time.Duration(float64(time.Second) * seconds)
}

// durationToTokens calculates the number of tokens for a given duration in
// a throttle with limit tokens per second.
func (limit Limit) durationToTokens(duration time.Duration) float64 {
	if limit <= 0 {
		return 0
	}
	return duration.Seconds() * float64(limit)
}

//--------------------
// THROTTLE
//--------------------

// Event wraps the event to be processed inside a function executed by a throttle.
type Event func() error

// Throttle controls the maximum number of events processed seconds per second. Here
// it internally uses a token bucket like described at https://en.wikipedia.org/wiki/Token_bucket.
// A throttle is created with the limit of allowed events per second and an additional
// burst size. A higher burst size allows to process more than event with the Process()
// method in one call.
type Throttle struct {
	mu         sync.RWMutex
	limit      Limit
	burst      int
	bucket     float64
	lastUpdate time.Time
	lastEvent  time.Time
}

// NewThrottle creates a new throttle with a limit of allowed events to process per
// second and a burst size for a possible number of events per call.
func NewThrottle(limit Limit, burst int) *Throttle {
	return &Throttle{
		limit: limit,
		burst: burst,
	}
}

// Limit returns the current limit of the throttle.
func (t *Throttle) Limit() Limit {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.limit
}

// Burst returns the current burst of the throttle.
func (t *Throttle) Burst() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.burst
}

// Process executes one or more events in a given context. If the limit is
// not infinite and the number of events is higher than the burst size or if
// the number of events is too high for the maximum duration, then the call
// will be declined. Also waiting may be too long if the context timeout is
// reached earlier which leads to an error.
func (t *Throttle) Process(ctx context.Context, events ...Event) error {
	// Check burst.
	t.mu.RLock()
	limit := t.limit
	burst := t.burst
	t.mu.RUnlock()
	if len(events) > burst && limit != InfiniteLimit {
		return fmt.Errorf("wait: processing %d event(s) exceeds throttle burst size %d", len(events), burst)
	}
	// Check if the context is already cancelled.
	select {
	case <-ctx.Done():
		return fmt.Errorf("wait: processing %d event(s) throttle context already done: %v", len(events), ctx.Err())
	default:
	}
	// Calculate a possible timeout.
	now := time.Now()
	maxDuration := InfiniteDuration
	if deadline, ok := ctx.Deadline(); ok {
		maxDuration = deadline.Sub(now)
	}
	// Create clock.
	c := newClock(t, len(events), now, maxDuration)
	if !c.ok {
		return fmt.Errorf("wait: processing %d event(s) would exceed throttle context deadline", len(events))
	}
	// Delay as long as needed.
	delay := c.delayFrom(now)
	if delay > 0 {
		t := time.NewTimer(delay)
		defer t.Stop()
		select {
		case <-t.C:
			// Continue.
			break
		case <-ctx.Done():
			// Context was canceled.
			c.cancelAt(time.Now())
			return fmt.Errorf("wait: processing %d event(s) throttle context timed out or cancelled: %v", len(events), ctx.Err())
		}
	}
	// Process event(s).
	for n, event := range events {
		if err := event(); err != nil {
			return fmt.Errorf("wait: processing event %d returned error: %v", n, err)
		}
	}
	return nil
}

//--------------------
// CLOCK
//--------------------

// clock is responsible for the controlling of limiit and tokens of a throttle.
type clock struct {
	throttle *Throttle
	limit    Limit
	ok       bool
	tokens   int
	act      time.Time
}

// newClock creates and initializes a clock for the given throttle and
// the processing of one or more events.
func newClock(throttle *Throttle, n int, now time.Time, maxReserve time.Duration) *clock {
	c := &clock{
		throttle: throttle,
	}

	c.throttle.mu.Lock()
	defer c.throttle.mu.Unlock()

	if c.throttle.limit == InfiniteLimit {
		c.ok = true
		c.tokens = n
		c.act = now

		return c
	} else if c.throttle.limit == 0 {
		var ok bool
		if c.throttle.burst >= n {
			ok = true
			c.throttle.burst -= n // Reduce throttle burst.
		}
		c.ok = ok
		c.tokens = c.throttle.burst
		c.act = now
		return c
	}

	now, lastUpdate, tokens := c.advanceTokens(now)

	// How many tokens remain after this request?
	tokens -= float64(n)

	// How long do we have to wait?
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = c.throttle.limit.tokensToDuration(-tokens)
	}

	// Decide if processing is okay.
	ok := n <= c.throttle.burst && waitDuration <= maxReserve

	c.ok = ok
	c.limit = c.throttle.limit

	if ok {
		c.tokens = n
		c.act = now.Add(waitDuration)
	}

	// Update throttle state.
	if ok {
		c.throttle.lastUpdate = now
		c.throttle.bucket = tokens
		c.throttle.lastEvent = c.act
	} else {
		c.throttle.lastUpdate = lastUpdate
	}

	return c
}

// advanceTokens calculates and returns an updated state for the throttle based on the given timestamp
// and without changing it. Due to access of fields the caller must have locked the throttle.
func (c *clock) advanceTokens(now time.Time) (newNow time.Time, newLastUpdate time.Time, newTokens float64) {
	// Check last update.
	lastUpdate := c.throttle.lastUpdate
	if now.Before(lastUpdate) {
		lastUpdate = now
	}
	// Calculate number of tokens.
	elapsed := now.Sub(lastUpdate)
	delta := c.throttle.limit.durationToTokens(elapsed)
	tokens := c.throttle.bucket + delta
	if burst := float64(c.throttle.burst); tokens > burst {
		tokens = burst
	}
	return now, lastUpdate, tokens
}

// delayFrom returns the duration for which the processing must wait before processing
// the event(s). When zero there's no need to wait while InfiniteDuration means that
// prossibly the tokens won't be available before timeout.
func (c *clock) delayFrom(now time.Time) time.Duration {
	if !c.ok {
		return InfiniteDuration
	}
	delay := c.act.Sub(now)
	if delay < 0 {
		return 0
	}
	return delay
}

// cancelAt is called in case of a timeout and so the event(s) are not processed. So
// the processor can reverse the effects on the throttle to allow other callers to
// process their event(s).
func (c *clock) cancelAt(now time.Time) {
	if !c.ok {
		return
	}

	c.throttle.mu.Lock()
	defer c.throttle.mu.Unlock()

	if c.throttle.limit == InfiniteLimit || c.tokens == 0 || c.act.Before(now) {
		return
	}

	// Calculate the tokens to restore in throttle.
	restoreTokens := float64(c.tokens) - c.limit.durationToTokens(c.throttle.lastEvent.Sub(c.act))
	if restoreTokens <= 0 {
		return
	}

	// Advance tokens to now.
	now, _, tokens := c.advanceTokens(now)

	// Now calculates the new number of tokens.
	tokens += restoreTokens
	if burst := float64(c.throttle.burst); tokens > burst {
		tokens = burst
	}

	// Finally update the throttle state.
	c.throttle.lastUpdate = now
	c.throttle.bucket = tokens
	if c.act == c.throttle.lastEvent {
		prevEvent := c.act.Add(c.limit.tokensToDuration(float64(-c.tokens)))
		if !prevEvent.Before(now) {
			c.throttle.lastEvent = prevEvent
		}
	}
}

// EOF
