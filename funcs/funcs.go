// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package funcs

// Func is general function, f:T -> R.
type Func[T, R any] func(T) R

// Predict is a function (f:T -> bool) predicts given value is true or false.
type Predict[T any] func(T) bool

// Unit is a function (f:empty -> R) mapping to R from empty set.
type Unit[R any] func() R

// Condition is a function (f:empty -> bool) mapping to true or false from empty set.
type Condition func() bool

// Partial is a function (f:T -> (R, bool)) converts partial value in T to R.
type Partial[T, R any] func(T) (R, bool)

// Try is a function (f:T -> (R, error)) transfers to R from T, and maybe error returned.
type Try[T, R any] func(T) (r R, err error)

// Transform is a function (f:(T, bool) -> R) transforms to R even if given value v is not fine.
type Transform[T, R any] func(v T, ok bool) R

// Recover is a function (f:(T, error) -> R) recover to R even if given value v is failed.
type Recover[T, R any] func(v T, err error) R

// Self always return given value v.
func Self[T any](v T) T {
	return v
}

// Id returns a function always return given value v.
func Id[T any](v T) Unit[T] {
	return func() T {
		return v
	}
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// AndThen return a new function (f:T -> R) applying given function g to result from f.
func AndThen[T, U, R any](f Func[T, U], g Func[U, R]) Func[T, R] {
	return func(v T) R {
		return g(f(v))
	}
}

// UnitAndThen returns a new function (f:empty -> R) applying given function g to result from f.
func UnitAndThen[T, R any](f Unit[T], g Func[T, R]) Unit[R] {
	return func() R {
		return g(f())
	}
}

// Compose returns a new function (f:T -> R) applying given function f to result from g.
func Compose[T, U, R any](f Func[U, R], g Func[T, U]) Func[T, R] {
	return func(v T) R {
		return f(g(v))
	}
}

// ComposeUnit returns a new function (f:empty -> R) applying given function f to result from g.
func ComposeUnit[T, R any](f Func[T, R], g Unit[T]) Unit[R] {
	return func() R {
		return f(g())
	}
}

// PartialTransform returns a new function (f:T -> R) applying given Transform f2 to result from f1.
func PartialTransform[T, U, R any](f1 Partial[T, U], f2 Transform[U, R]) Func[T, R] {
	return func(v T) R {
		return f2(f1(v))
	}
}

// TryRecover returns a new function (f:T -> R) applying given Recover f2 to result from f1.
func TryRecover[T, U, R any](f1 Try[T, U], f2 Recover[U, R]) Func[T, R] {
	return func(v T) R {
		return f2(f1(v))
	}
}

// Cond is a ternary returning given value succ if ok is true, or returning fail.
func Cond[T any](ok bool, succ T, fail T) T {
	if ok {
		return succ
	}
	return fail
}

/*
func ConfFunc[T any](p Condition, succ Unit[T], fail Unit[T]) T {
	if p() {
		return succ()
	}
	return fail()
}
*/
