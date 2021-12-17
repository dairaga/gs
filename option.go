// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"fmt"
	"reflect"

	"github.com/dairaga/gs/funcs"
)

// Option is simplified Scala Option. Option like Either is either Some or None.
// Some means this is defined and has a value; None means Nothing or Nil.
// Suggest to return Option from function or method instead of nil.
type Option[T any] interface {
	fmt.Stringer

	// Fetch returns value v from Some and err is nil if this is a Some, otherwise v is a zero value and err is ErrEmpty.
	Fetch() (v T, err error)

	// Check returns value v from Some and ok is true if this is a Some, otherwise v is a zero value and ok is false.
	Check() (v T, ok bool)

	// Get returns value from Some, or panic.
	Get() T

	// IsDefined returns true if this is a Some.
	IsDefined() bool

	// IsEmpty returns true if this is a None.
	IsEmpty() bool

	// Exists returns true if this is a Some and value satisifies given function p.
	Exists(p funcs.Predict[T]) bool

	// Forall returns true if this is a None or value from Succes satisfies given function p.
	Forall(p funcs.Predict[T]) bool

	// Foreach only applies given function op to value from Success.
	Foreach(op func(T))

	// Filter returns this if this is a None or value from Some satisfies given function p,
	// otherwise returns None.
	Filter(p funcs.Predict[T]) Option[T]

	// FilterNot returns this if this is a None or value from Some satisfies given function p, otherwise returns None.
	FilterNot(funcs.Predict[T]) Option[T]

	// GetOrElse returns value from Some, or returns given z.
	GetOrElse(z T) T

	// OrElse returns this if this is a Some, or returns given z.
	OrElse(Option[T]) Option[T]

	// Try returns Success with value of Some, or Failure with ErrEmpty.
	Try() Try[T]

	// Either returns Rigt with value from Some, or Left with ErrEmpty.
	Either() Either[error, T]
}

type option[T any] struct {
	*try[T]
}

var _ Option[int] = &option[int]{}

func (o *option[T]) String() string {
	if o.ok {
		return fmt.Sprintf(`Some(%v)`, o.right)
	}
	return fmt.Sprintf(`None(%s)`, reflect.TypeOf(o.right).String())
}

func (o *option[T]) Fetch() (T, error) {
	return o.right, o.left
}

func (o *option[T]) Check() (T, bool) {
	return o.right, o.ok
}

func (o *option[T]) IsDefined() bool {
	return o.ok
}

func (o *option[T]) IsEmpty() bool {
	return !o.ok
}

func (o *option[T]) Filter(p funcs.Predict[T]) Option[T] {
	if o.Forall(p) {
		return o
	}
	return none[T]()
}

func (o *option[T]) FilterNot(p funcs.Predict[T]) Option[T] {
	return o.Filter(func(v T) bool { return !p(v) })
}

func (o *option[T]) OrElse(z Option[T]) Option[T] {
	return funcs.Cond(o.ok, Option[T](o), z)
}

func (o *option[T]) Try() Try[T] {
	return o.try
}

func some[T any](v T) *option[T] {
	return &option[T]{
		try: success(v),
	}
}

func none[T any]() *option[T] {
	return &option[T]{
		try: failure[T](ErrEmpty),
	}
}

func Some[T any](v T) Option[T] {
	return some(v)
}

func None[T any]() Option[T] {
	return none[T]()
}
