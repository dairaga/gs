// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"fmt"

	"github.com/dairaga/gs/funcs"
)

type Either[L, R any] interface {
	fmt.Stringer
	Fetch() (R, error)
	Get() R

	IsRight() bool
	Right() R

	IsLeft() bool
	Left() L

	Exists(funcs.Predict[R]) bool
	Forall(funcs.Predict[R]) bool
	Foreach(func(R))

	FilterOrElse(funcs.Predict[R], L) Either[L, R]
	GetOrElse(R) R
	OrElse(Either[L, R]) Either[L, R]

	Swap() Either[R, L]

	Try() Try[R]
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
