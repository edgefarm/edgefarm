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
	"math/rand"
	"time"
)

//--------------------
// TICKER
//--------------------

// TickerFunc defines a function sending signals for each condition
// check when polling. The ticker can be canceled via the given
// context. Closing the returned signal channel means that the
// ticker ended. Sending ticks should be able to handle not
// received ones in case the condition check of the poller is
// working.
type TickerFunc func(ctx context.Context) <-chan struct{}

// TickChangerFunc allows to work with changing intervals. The
// current one is the argument, the next has to be returned. In
// case the bool return value is false the ticker will stop.
type TickChangerFunc func(in time.Duration) (out time.Duration, ok bool)

// MakeGenericIntervalTicker is a factory for tickers based on time
// intervals. The given changer is responsible for the intervals and
// if the ticker shall signal a stopping. The changer is called initially
// with a duration of zero to allow the changer stopping the ticker even
// before a first tick.
func MakeGenericIntervalTicker(changer TickChangerFunc) TickerFunc {
	return func(ctx context.Context) <-chan struct{} {
		tickc := make(chan struct{})
		interval := 0 * time.Millisecond
		ok := true
		go func() {
			defer close(tickc)
			// Defensive changer call.
			if interval, ok = changer(interval); !ok {
				return
			}
			// TickerFunc for the interval.
			timer := time.NewTimer(interval)
			defer timer.Stop()
			// Loop sending signals.
			for {
				select {
				case <-timer.C:
					// One interval tick. Ignore if needed.
					select {
					case tickc <- struct{}{}:
					default:
					}
				case <-ctx.Done():
					// Given context stopped.
					return
				}
				// Reset timer with next interval.
				if interval, ok = changer(interval); !ok {
					return
				}
				timer.Reset(interval)
			}
		}()
		return tickc
	}
}

// MakeIntervalTicker returns a ticker signalling in intervals.
func MakeIntervalTicker(interval time.Duration) TickerFunc {
	changer := func(_ time.Duration) (out time.Duration, ok bool) {
		return interval, true
	}
	return MakeGenericIntervalTicker(changer)
}

// MakeMaxIntervalsTicker returns a ticker signalling in intervals. It
// stops after a maximum number of signals.
func MakeMaxIntervalsTicker(interval time.Duration, max int) TickerFunc {
	count := 0
	changer := func(_ time.Duration) (out time.Duration, ok bool) {
		count++
		if count > max {
			return 0, false
		}
		return interval, true
	}
	return MakeGenericIntervalTicker(changer)
}

// MakeDeadlinedIntervalTicker returns a ticker signalling in intervals
// and stopping after a deadline.
func MakeDeadlinedIntervalTicker(interval time.Duration, deadline time.Time) TickerFunc {
	changer := func(_ time.Duration) (out time.Duration, ok bool) {
		if time.Now().After(deadline) {
			return 0, false
		}
		return interval, true
	}
	return MakeGenericIntervalTicker(changer)
}

// MakeExpiringIntervalTicker returns a ticker signalling in intervals
// and stopping after a timeout.
func MakeExpiringIntervalTicker(interval, timeout time.Duration) TickerFunc {
	deadline := time.Now().Add(timeout)
	changer := func(_ time.Duration) (out time.Duration, ok bool) {
		if time.Now().After(deadline) {
			return 0, false
		}
		return interval, true
	}
	return MakeGenericIntervalTicker(changer)
}

// MakeJitteringTicker returns a ticker signalling in jittering intervals. This
// avoids converging on periadoc behavior during condition check. The returned
// intervals jitter between the given interval and interval + factor * interval.
// The ticker stops after reaching timeout.
func MakeJitteringTicker(interval time.Duration, factor float64, timeout time.Duration) TickerFunc {
	deadline := time.Now().Add(timeout)
	changer := func(_ time.Duration) (time.Duration, bool) {
		if time.Now().After(deadline) {
			return 0, false
		}
		if factor <= 0.0 {
			factor = 1.0
		}
		return interval + time.Duration(rand.Float64()*factor*float64(interval)), true
	}
	return MakeGenericIntervalTicker(changer)
}

// EOF
