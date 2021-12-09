// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"fmt"

	"github.com/dairaga/gs/funcs"
)

type Try[T any] interface {
	fmt.Stringer
	Fetch() (T, error)
	Get() T

	IsSuccess() bool
	Success() T

	IsFailure() bool
	Failed() error

	Exists(funcs.Predict[T]) bool
	Forall(funcs.Predict[T]) bool
	Foreach(func(T))

	Filter(funcs.Predict[T]) Try[T]
	FilterNot(funcs.Predict[T]) Try[T]
	GetOrElse(T) T
	OrElse(Try[T]) Try[T]

	Recover(funcs.Func[error, T]) Try[T]
	RecoverWith(funcs.Func[error, Try[T]]) Try[T]

	Either() Either[error, T]
	Option() Option[T]
}

type try[T any] struct {
	*either[error, T]
}

var _ Try[int] = &try[int]{}

func (t *try[T]) String() string {
	if t.ok {
		return fmt.Sprintf(`Success(%v)`, t.right)
	}
	return fmt.Sprintf(`Failure(%s)`, t.left.Error())
}

func (t *try[T]) Fetch() (T, error) {
	return t.right, t.left
}

func (t *try[T]) IsSuccess() bool {
	return t.ok
}

func (t *try[T]) Success() T {
	return t.Get()
}

func (t *try[T]) IsFailure() bool {
	return !t.ok
}

func (t *try[T]) Failed() error {
	if t.IsFailure() {
		return t.left
	}
	return ErrUnsupported
}

func (t *try[T]) Filter(p funcs.Predict[T]) Try[T] {
	if t.Forall(p) {
		return t
	}

	return failure[T](ErrUnsatisfied)
}

func (t *try[T]) FilterNot(p funcs.Predict[T]) Try[T] {
	return t.Filter(func(v T) bool { return !p(v) })
}

func (t *try[T]) OrElse(z Try[T]) Try[T] {
	return funcs.Cond(t.ok, Try[T](t), z)
}

func (t *try[T]) Recover(r funcs.Func[error, T]) Try[T] {
	if t.ok {
		return t
	}
	return success(r(t.left))
}

func (t *try[T]) RecoverWith(r funcs.Func[error, Try[T]]) Try[T] {
	if t.ok {
		return t
	}
	return r(t.left)
}

func (t *try[T]) Either() Either[error, T] {
	return t.either
}

func (t *try[T]) Option() Option[T] {
	if t.ok {
		return some(t.right)
	}
	return none[T]()
}

func success[T any](v T) *try[T] {
	return &try[T]{
		either: right[error](v),
	}
}

func failure[T any](err error) *try[T] {
	return &try[T]{
		either: left[error, T](err),
	}
}

func Success[T any](v T) Try[T] {
	return success(v)
}

func Failure[T any](err error) Try[T] {
	return failure[T](err)
}
