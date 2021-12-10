// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"fmt"

	"github.com/dairaga/gs/funcs"
)

// Try is simplified Scala Try. Try like Either is either Success or Failure,
// and Failure contains error value.
type Try[T any] interface {
	fmt.Stringer

	// Fetch returns successful value and nil error if this is a Success,
	// or v is zero value and err is from Failure.
	Fetch() (v T, err error)

	// Get returns successful value if Try is a Success, or panic.
	Get() T

	// IsSuccess returns true if this is a Success.
	IsSuccess() bool

	// Success returns successful value if this is a Sucess, or panic.
	Success() T

	// IsFailure returns true if this is a Failure.
	IsFailure() bool

	// Failed returns error from Failure, or returns ErrUnsupported.
	Failed() error

	// Exists returns true if this is a Success and satisfies given function p.
	Exists(p funcs.Predict[T]) bool

	// Forall returns true if this is a Failure,
	// or value from Success satisfies given function p.
	Forall(p funcs.Predict[T]) bool

	// Foreach only applies given function op to value from Success.
	Foreach(op func(T))

	// Filter returns this if this is a Failure,
	// or value from Success satisfies given function p,
	// otherwise returns Failure with ErrUnsatisfied.
	Filter(p funcs.Predict[T]) Try[T]

	// FilterNot returns this if this is a Failure,
	// or value from Succes does not satisfy given function p,
	// otherwise returns Failure with ErrUnsatisfied.
	FilterNot(funcs.Predict[T]) Try[T]

	// GetOrElse returns value from Success, or returns given z.
	GetOrElse(z T) T

	// OrElse returns this if this is a Success, or return given z.
	OrElse(z Try[T]) Try[T]

	// Recover applies given function r if this is a Failure,
	// or returns this if this is a Success.
	Recover(r funcs.Func[error, T]) Try[T]

	// RecoverWith applies given function r if this is a Failure,
	// or returns this if this is a Success.
	RecoverWith(r funcs.Func[error, Try[T]]) Try[T]

	// Either returns a Either with error type in Left side.
	Either() Either[error, T]

	// Option returns Some if this is a Success, or returns None.
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
