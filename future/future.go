// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package future

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/try"
)

type F[T any] struct {
	_         struct{}
	ctx       context.Context
	cancel    context.CancelFunc
	completed bool
	result    gs.Try[T]
}

var _ gs.Future[int] = &F[int]{}

func (f *F[T]) assign(x gs.Try[T]) *F[T] {
	f.result = x
	f.completed = true
	f.cancel()
	return f
}

func (f *F[T]) String() string {
	if f.completed {
		return fmt.Sprintf(`Completed(%v)`, f.result)
	}
	return fmt.Sprintf(`Future(?)`)
}

func (f *F[T]) Completed() bool {
	return f.completed
}

func (f *F[T]) Get() (gs.Try[T], bool) {
	return f.result, f.completed
}

func (f *F[T]) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *F[T]) OnCompleted(op func(gs.Try[T])) {
	go func() {
		<-f.Done()
		if f.completed {
			op(f.result)
		}
	}()
}

func (f *F[T]) OnSuccess(op func(T)) {
	go func() {
		<-f.Done()
		if f.completed && f.result.IsSuccess() {
			op(f.result.Get())
		}
	}()
}

func (f *F[T]) OnError(op func(error)) {
	go func() {
		<-f.Done()
		if f.completed && f.result.IsFailure() {
			op(f.result.Failed())
		}
	}()
}

func (f *F[T]) Wait() gs.Try[T] {
	if f.completed {
		return f.result
	}
	<-f.Done()
	return f.result
}

func (f *F[T]) Result(ctx context.Context, atMost time.Duration) gs.Try[T] {
	wait, cancel := context.WithTimeout(ctx, atMost)
	defer cancel()

	select {
	case <-f.Done():
		if f.completed {
			return f.result
		}
		return gs.Failure[T](f.ctx.Err())
	case <-wait.Done():
		return gs.Failure[T](wait.Err())
	}
}

func (f *F[T]) Filter(ctx context.Context, p funcs.Predict[T]) gs.Future[T] {
	return TransformWith[T](ctx, f, func(x gs.Try[T]) gs.Future[T] {
		return promise[T](ctx).assign(
			funcs.Cond(
				x.IsFailure() || p(x.Success()),
				x,
				gs.Failure[T](gs.ErrUnsatisfied),
			),
		)
	})
}

func (f *F[T]) FilterNot(ctx context.Context, p funcs.Predict[T]) gs.Future[T] {
	return f.Filter(ctx, func(v T) bool { return !p(v) })
}

// Failure returns default failed result.
func Failure[T any]() gs.Try[T] {
	return gs.Failure[T](gs.ErrEmpty)
}

// promise is waiting for something in future.
func promise[T any](parent context.Context) *F[T] {
	ret := &F[T]{
		completed: false,
		result:    gs.Failure[T](gs.ErrEmpty),
	}
	ret.ctx, ret.cancel = context.WithCancel(parent)
	return ret
}

// Run returns a Future waitng for the result from given function op.
func Run[T any](parent context.Context, op func() T) gs.Future[T] {
	ret := promise[T](parent)

	go func(op func() T, f *F[T]) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					f.result = gs.Failure[T](v)
				default:
					f.result = gs.Failure[T](fmt.Errorf(`%v`, v))
				}
			}
			f.completed = true
			f.cancel()
		}()

		f.result = gs.Success(op())
	}(op, ret)

	return ret
}

// Try returns a Future waitng for the result from given function op.
func Try[T any](parent context.Context, op func() (T, error)) gs.Future[T] {
	ret := promise[T](parent)
	go func(f *F[T]) {
		f.assign(try.From(op()))
	}(ret)
	return ret
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// Transform returns a Future waiting for the result applied by given function to result of given future f.
func Transform[T, U any](ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Try[U]) gs.Future[U] {
	ret := promise[U](ctx)

	go func(ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Try[U], ret *F[U]) {
		select {
		case <-f.Done():
			result, completed := f.Get()
			if completed {
				ret.result = op(result)
				ret.completed = true
			}
		case <-ret.Done():
		}
		ret.cancel()
	}(ctx, f, op, ret)

	return ret
}

// TransformWith returns a new Future to wait another Future made by applying the given function op to the result of future f.
func TransformWith[T, U any](ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Future[U]) gs.Future[U] {
	ret := promise[U](ctx)

	go func(ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Future[U], ret *F[U]) {
		select {
		case <-f.Done():
			fresult, fcompleted := f.Get()
			if fcompleted {
				g := op(fresult)
				go func() {
					select {
					case <-g.Done():
						gresult, gcompleted := g.Get()
						if gcompleted {
							ret.result = gresult
							ret.completed = true
						}
					case <-ret.Done():
					}
					ret.cancel()
				}()
			}
		case <-ret.Done():
			ret.cancel()
		}
	}(ctx, f, op, ret)

	return ret
}

// FlatMap returns a new Future by applying given function op to the successful result of future f, and returns the result of the function as the new future.
func FlatMap[T, U any](ctx context.Context, f gs.Future[T], op func(T) gs.Future[U]) gs.Future[U] {
	return TransformWith(ctx, f, func(x gs.Try[T]) gs.Future[U] {
		if x.IsSuccess() {
			return op(x.Success())
		}
		return promise[U](ctx).assign(gs.Failure[U](x.Failed()))
	})
}

// Map returns a new Future by applying given function op to the successful result of future f.
func Map[T, U any](ctx context.Context, f gs.Future[T], op func(T) U) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.Map(x, op)
	})
}

// TryMap returns a new Future by applying given function op to the result of future f.
func TryMap[T, U any](ctx context.Context, f gs.Future[T], op funcs.Try[T, U]) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.TryMap(x, op)
	})
}

// PartialMap returns a new Future by applying given function op to the result of future f.
func PartialMap[T, U any](ctx context.Context, f gs.Future[T], op funcs.Partial[T, U]) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.PartialMap(x, op)
	})
}

// Zip zips the values of futures f and g, and creates a new future holding the tuple of their results.
func Zip[T, U any](ctx context.Context, f gs.Future[T], g gs.Future[U]) gs.Future[gs.Tuple2[gs.Try[T], gs.Try[U]]] {
	ret := promise[gs.Tuple2[gs.Try[T], gs.Try[U]]](ctx)

	go func(f gs.Future[T], g gs.Future[U], ret *F[gs.Tuple2[gs.Try[T], gs.Try[U]]]) {
		defer ret.cancel()

		cases := []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ret.Done()),
			},
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(f.Done()),
			},
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(g.Done()),
			},
		}

		var fresult gs.Try[T]
		var gresult gs.Try[U]
		var fcompleted, gcompleted bool

		for !f.Completed() || !g.Completed() {
			chosen, _, _ := reflect.Select(cases)
			if chosen == 0 {
				break
			}

			switch chosen {
			case 1:
				cases = append(cases[0:1], cases[2:]...)
				if fresult == nil {
					fresult, fcompleted = f.Get()
				} else {
					gresult, gcompleted = g.Get()
				}

			case 2:
				cases = cases[0:2]
				gresult, gcompleted = g.Get()
			}
		}

		ret.completed = fcompleted && gcompleted
		if ret.completed {
			ret.result = gs.Success(gs.T2(fresult, gresult))
		}
	}(f, g, ret)

	return ret
}

// ZipWith zip the values of futures f and g using given function op, and creates a new future holding the result.
func ZipWith[T, U, R any](ctx context.Context, f gs.Future[T], g gs.Future[U], op func(gs.Try[T], gs.Try[U]) gs.Try[R]) gs.Future[R] {
	return Transform(ctx, Zip(ctx, f, g), func(x gs.Try[gs.Tuple2[gs.Try[T], gs.Try[U]]]) gs.Try[R] {
		return try.FlatMap(x, func(tup gs.Tuple2[gs.Try[T], gs.Try[U]]) gs.Try[R] {
			return op(tup.V1, tup.V2)
		})
	})
}
