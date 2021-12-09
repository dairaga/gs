// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package future

import (
	"context"
	"fmt"
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

func (f *F[T]) Foreach(op func(T)) {
	f.OnCompleted(func(x gs.Try[T]) {
		x.Foreach(op)
	})
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
		return promise[T](ctx).
			assign(
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

// -----------------------------------------------------------------------------

func promise[T any](parent context.Context) *F[T] {
	ret := &F[T]{}
	ret.ctx, ret.cancel = context.WithCancel(parent)
	return ret
}

func Run[T any](parent context.Context, op func() T) gs.Future[T] {
	ret := promise[T](parent)

	go func(f *F[T]) {
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
	}(ret)

	return ret
}

func Try[T any](parent context.Context, op func() (T, error)) gs.Future[T] {
	ret := promise[T](parent)
	go func(f *F[T]) {
		f.assign(try.From(op()))
	}(ret)
	return ret
}

func Transform[T, U any](ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Try[U]) gs.Future[U] {
	ret := promise[U](ctx)

	go func() {
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
	}()

	return ret
}

func TransformWith[T, U any](ctx context.Context, f gs.Future[T], op func(gs.Try[T]) gs.Future[U]) gs.Future[U] {
	ret := promise[U](ctx)

	go func() {
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
	}()

	return ret
}

func FlatMap[T, U any](ctx context.Context, f gs.Future[T], op func(T) gs.Future[U]) gs.Future[U] {
	return TransformWith(ctx, f, func(x gs.Try[T]) gs.Future[U] {
		if x.IsSuccess() {
			return op(x.Success())
		}
		return promise[U](ctx).assign(gs.Failure[U](x.Failed()))
	})
}

func Map[T, U any](ctx context.Context, f gs.Future[T], op func(T) U) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.Map(x, op)
	})
}

func TryMap[T, U any](ctx context.Context, f gs.Future[T], op funcs.Try[T, U]) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.TryMap(x, op)
	})
}

func CanMap[T, U any](ctx context.Context, f gs.Future[T], op funcs.Can[T, U]) gs.Future[U] {
	return Transform(ctx, f, func(x gs.Try[T]) gs.Try[U] {
		return try.CanMap(x, op)
	})
}
