package funcs

import "constraints"

type Compare[T, U any] func(T, U) int
type Order[T any, R constraints.Ordered] func(T) R
type Equal[T, U any] func(T, U) bool

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

func Max[T constraints.Ordered](a, b T) T {
	return Cond(a >= b, a, b)
}

func Min[T constraints.Ordered](a, b T) T {
	return Cond(a <= b, a, b)
}
