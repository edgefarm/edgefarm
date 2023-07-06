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
	"time"
)

//--------------------
// POLL
//--------------------

// ConditionFunc has to be implemented for checking the wanted condition. A positive
// condition will return true and nil, a negative false and nil. In case of failure
// during the check false and the error have to be returned. The function will
// be used by the poll functions.
type ConditionFunc func() (bool, error)

// Poll provides different ways to wait for conditions by polling. The conditions
// are checked by user defined functions with the signature
//
//     func() (ok bool, err error)
//
// Here the bool return value signals if the condition is fulfilled, e.g. a file
// you're waiting for has been written into the according directory.
//
// This signal for check a condition is returned by a ticker with the signature
//
//     func(ctx context.Context) <-chan struct{}
//
// The context is for signalling the ticker to end working, the channel for the signals.
// Pre-defined tickers support
//
// - simple constant intervals,
// - a maximum number of constant intervals,
// - a constant number of intervals with a deadline,
// - a constant number of intervals with a timeout, and
// - jittering intervals.
//
// The behaviour of changing intervals can be user-defined by
// functions with the signature
//
//     func(in time.Duration) (out time.Duration, ok bool)
//
// Here the argument is the current interval, return values are the
// wanted interval and if the polling shall continue. For the predefined
// tickers according convenience functions named With...() exist.
//
// Example (waiting for a file to exist):
//
//     // Tick every second for maximal 30 seconds.
//     ticker := wait.MakeExpiringIntervalTicker(time.Second, 30*time.Second),
//
//     // Check for existence of a file.
//     condition := func() (bool, error) {
//         _, err := os.Stat("myfile.txt")
//         if err != nil {
//             if os.IsNotExist(err) {
//                 return false, nil
//             }
//             return false, err
//         }
//         // Found file.
//         return true, nil
//     }
//
//     // And now poll.
//     wait.Poll(ctx, ticker, condition)
//
// From outside the polling can be stopped by cancelling the context.
func Poll(ctx context.Context, ticker TickerFunc, condition ConditionFunc) error {
	tickCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tickc := ticker(tickCtx)
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return fmt.Errorf("context has been cancelled with error: %v", ctx.Err())
			}
			return nil
		case _, open := <-tickc:
			// Ticker sent a signal to check for condition.
			if !open {
				// Oh, ticker tells to end.
				return fmt.Errorf("ticker exceeded while waiting for the condition")
			}
			ok, err := check(condition)
			if err != nil {
				// ConditionFunc has an error.
				return fmt.Errorf("poll condition returned error: %v", err)
			}
			if ok {
				// ConditionFunc is happy.
				return nil
			}
		}
	}
}

// WithInterval is convenience for Poll() with MakeIntervalTicker().
func WithInterval(
	ctx context.Context,
	interval time.Duration,
	condition ConditionFunc,
) error {
	return Poll(
		ctx,
		MakeIntervalTicker(interval),
		condition,
	)
}

// WithMaxIntervals is convenience for Poll() with MakeMaxIntervalsTicker().
func WithMaxIntervals(
	ctx context.Context,
	interval time.Duration,
	max int,
	condition ConditionFunc,
) error {
	return Poll(
		ctx,
		MakeMaxIntervalsTicker(interval, max),
		condition,
	)
}

// WithDeadline is convenience for Poll() with MakeDeadlinedIntervalTicker().
func WithDeadline(
	ctx context.Context,
	interval time.Duration,
	deadline time.Time,
	condition ConditionFunc,
) error {
	return Poll(
		ctx,
		MakeDeadlinedIntervalTicker(interval, deadline),
		condition,
	)
}

// WithTimeout is convenience for Poll() with MakeExpiringIntervalTicker().
func WithTimeout(
	ctx context.Context,
	interval, timeout time.Duration,
	condition ConditionFunc,
) error {
	return Poll(
		ctx,
		MakeExpiringIntervalTicker(interval, timeout),
		condition,
	)
}

// WithJitter is convenience for Poll() with MakeJitteringTicker().
func WithJitter(
	ctx context.Context,
	interval time.Duration,
	factor float64,
	timeout time.Duration,
	condition ConditionFunc,
) error {
	return Poll(
		ctx,
		MakeJitteringTicker(interval, factor, timeout),
		condition,
	)
}

//--------------------
// PRIVATE HELPER
//--------------------

// check runs the condition catching potential panics and returns
// them as failure.
func check(condition ConditionFunc) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			err = fmt.Errorf("panic during condition check: %v", r)
		}
	}()
	ok, err = condition()
	return
}

// EOF
