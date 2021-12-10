// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"context"
	"fmt"
	"time"

	"github.com/dairaga/gs/funcs"
)

// Future imitates Scala Future. It waits result from some goroutine.
// Result is Try type either Success or Failure.
// Unlike Scala, it leverages Context to break a Future.
type Future[T any] interface {
	fmt.Stringer

	// Done returns channel to check goroutine is done.
	Done() <-chan struct{}

	// Completed returns true if goroutine is completed.
	Completed() bool

	// Get returns result and true in ok if goroutine is completed.
	Get() (result Try[T], ok bool)

	// OnSuccess applies given function p when goroutine is completed,
	// and result is a Success.
	OnSuccess(p func(T))

	// OnError applies given function p when goroutine is completed,
	// and result is a Failure.
	OnError(p func(error))

	// OnCompleted always applies given function p when goroutine is completed,
	// even giving a function to OnSuccess or OnError.
	OnCompleted(p func(Try[T]))

	// Forall applies given function p to value from Success result.
	Foreach(op func(T))

	// Filter returns a new Future to wait this and apply result to given function p.
	// Returned Future contains result from this if result is Failure or satisfies given function p,
	// or contains Failure with ErrUnsatisfied.
	Filter(ctx context.Context, p funcs.Predict[T]) Future[T]

	// FilterNot returns a new Future to wait this and apply result to given function p.
	// Returned Future contains result from this if result is Failure or dose not satisfies given function p,
	// or contains Failure with ErrUnsatisfied.
	FilterNot(context.Context, funcs.Predict[T]) Future[T]

	// Result waits result at most given time t.
	Result(ctx context.Context, t time.Duration) Try[T]

	// Wait waits result forever.
	Wait() Try[T]
}
