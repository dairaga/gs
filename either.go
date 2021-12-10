// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"fmt"

	"github.com/dairaga/gs/funcs"
)

// Either is simplified Scala Either. Concept of Ether is either A or B.
// Either has two subclass Right and Left in Scala.
// Conventionally, Right means Positive and True; Left means Negative and False.
type Either[L, R any] interface {
	fmt.Stringer

	// Fetch returns right value if this is a Right.
	// err is original value if this is a Left with error type, or ErrLeft.
	Fetch() (r R, err error)

	// Get returns Right value if this is a Right, or panic.
	Get() R

	// IsRight returns true if this is a Right.
	IsRight() bool

	// Right returns Right value if this is a Right, or panic.
	Right() R

	// IsLeft returns true if this is a Left.
	IsLeft() bool

	// Left returns Left value if this is a Left, or panic.
	Left() L

	// Exists returns true if this is a Right,
	// and value satisfies given function p.
	Exists(p funcs.Predict[R]) bool

	// Forall returns true if this is a Left,
	// or value from Right satisfies given function p.
	Forall(p funcs.Predict[R]) bool

	// Foreach only applies given function op to value from Right.
	Foreach(op func(R))

	// FilterNoElse returns this if this is a Right and satisfies given function p,
	// or returns Left with z.
	FilterOrElse(p funcs.Predict[R], z L) Either[L, R]

	// GetOrElse returns value from Right, or returns z.
	GetOrElse(z R) R

	// OrElse returns this if this is a Right, or returns z.
	OrElse(z Either[L, R]) Either[L, R]

	// Swap returns new Either swapping Right and Left side.
	Swap() Either[R, L]

	// Try converts to Try. Right converts to Success,
	// and Left converts to Failure.
	// Left value is reserved if Left type is error,
	// or returns Failure with ErrLeft.
	Try() Try[R]

	// Option converts to Option. Right converts to Some,
	// and Left converts to None.
	Option() Option[R]
}

type either[L, R any] struct {
	_     struct{}
	ok    bool
	left  L
	right R
}

var _ Either[int, int] = &either[int, int]{}

func (e *either[L, R]) String() string {
	if e.ok {
		return fmt.Sprintf(`Right(%v)`, e.right)
	}
	return fmt.Sprintf(`Left(%v)`, e.left)
}

func (e *either[L, R]) Fetch() (R, error) {
	return e.right, funcs.Cond(e.ok, nil, err(e.left))
}

func (e *either[L, R]) Get() R {
	if e.ok {
		return e.right
	}
	panic(ErrEmpty)
}

func (e *either[L, R]) IsRight() bool {
	return e.ok
}

func (e *either[L, R]) Right() R {
	return e.Get()
}

func (e *either[L, R]) IsLeft() bool {
	return !e.ok
}
func (e *either[L, R]) Left() L {
	if !e.ok {
		return e.left
	}
	panic(ErrEmpty)
}

func (e *either[L, R]) Exists(p funcs.Predict[R]) bool {
	return funcs.Fetcher[R](e.Fetch).Exists(p)
}

func (e *either[L, R]) Forall(p funcs.Predict[R]) bool {
	return funcs.Fetcher[R](e.Fetch).Forall(p)
}

func (e *either[L, R]) Foreach(op func(R)) {
	funcs.Fetcher[R](e.Fetch).Foreach(op)
}

func (e *either[L, R]) FilterOrElse(p funcs.Predict[R], z L) Either[L, R] {
	if e.Forall(p) {
		return e
	}
	return left[L, R](z)
}
func (e *either[L, R]) GetOrElse(z R) R {
	return funcs.Fetcher[R](e.Fetch).GetOrElse(z)
}

func (e *either[L, R]) OrElse(z Either[L, R]) Either[L, R] {
	return funcs.Cond(e.ok, Either[L, R](e), z)
}

func (e *either[L, R]) Swap() Either[R, L] {
	return &either[R, L]{
		ok:    !e.ok,
		left:  e.right,
		right: e.left,
	}
}

func (e *either[L, R]) Try() Try[R] {
	if e.ok {
		return success(e.right)
	}
	return failure[R](err(e.left))
}

func err(x interface{}) error {
	switch v := x.(type) {
	case error:
		return v
	default:
		return ErrLeft
	}
}

func (e *either[L, R]) Option() Option[R] {
	if e.ok {
		return some(e.right)
	}
	return none[R]()
}

func left[L, R any](v L) *either[L, R] {
	return &either[L, R]{
		ok:   false,
		left: v,
	}
}

func right[L, R any](v R) *either[L, R] {
	return &either[L, R]{
		ok:    true,
		right: v,
	}
}

func Left[L, R any](v L) Either[L, R] {
	return left[L, R](v)
}

func Right[L, R any](v R) Either[L, R] {
	return right[L](v)
}
