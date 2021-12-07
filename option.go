package gs

import (
	"fmt"
	"reflect"

	"github.com/dairaga/gs/funcs"
)

type Option[T any] interface {
	fmt.Stringer
	Fetch() (T, error)
	Check() (T, bool)
	Get() T

	IsDefined() bool
	IsEmpty() bool

	Exists(funcs.Predict[T]) bool
	Forall(funcs.Predict[T]) bool
	Foreach(func(T))

	Filter(funcs.Predict[T]) Option[T]
	FilterNot(funcs.Predict[T]) Option[T]
	GetOrElse(T) T
	OrElse(Option[T]) Option[T]

	Try() Try[T]
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
