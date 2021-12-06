package gs

import (
	"context"
	"fmt"
	"time"

	"github.com/dairaga/gs/funcs"
)

type Future[T any] interface {
	fmt.Stringer

	Done() <-chan struct{}
	Completed() bool
	Get() (Try[T], bool)

	OnSuccess(func(T))
	OnError(func(error))
	OnCompleted(func(Try[T]))

	Foreach(func(T))

	Filter(context.Context, funcs.Predict[T]) Future[T]
	FilterNot(context.Context, funcs.Predict[T]) Future[T]

	Result(context.Context, time.Duration) Try[T]
	Wait() Try[T]
}
