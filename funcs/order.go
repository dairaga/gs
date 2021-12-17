// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package funcs

import "constraints"

// Ordering is a function to find ordering relation of two elements.
type Ordering[T, U any] func(T, U) int

// Orderize is a function to convert unordered value to ordered one.
type Orderize[T any, R constraints.Ordered] func(T) R

// Equal is a function to check two elements is equal.
type Equal[T, U any] func(T, U) bool

// Cmp compares two ordered value and returns 1 if given a is larger than b,
// or returns -1 if a is less than b, otherwise returns 0.
func Cmp[T constraints.Ordered](a, b T) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Max returns the maximum from given ordered value a and b.
func Max[T constraints.Ordered](a, b T) T {
	return Cond(a >= b, a, b)
}

// Min returns the minimum from given ordered value a and b.
func Min[T constraints.Ordered](a, b T) T {
	return Cond(a <= b, a, b)
}
