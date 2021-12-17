// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package slices

import (
	"constraints"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

type S[T any] []T

// Empty returns an empty slice.
func Empty[T any]() S[T] {
	return []T{}
}

// One returns an one element slice.
func One[T any](v T) S[T] {
	return []T{v}
}

// From returns a slice from given elements.
func From[T any](a ...T) S[T] {
	return a
}

// Fill returns a slice with length n and filled with given element.
func Fill[T any](size int, v T) S[T] {
	ret := make([]T, size)
	for i := range ret {
		ret[i] = v
	}
	return ret
}

// FillWith returns a slice that contains the results of applying given function op with n times.
func FillWith[T any](n int, op funcs.Unit[T]) S[T] {
	ret := make(S[T], n)
	for i := range ret {
		ret[i] = op()
	}
	return ret
}

// Range returns a slice containing equally spaced values in a gevin interval.
func Range[T gs.Numeric](start, end, step T) S[T] {
	ret := Empty[T]()

	for i := start; i < end; i += step {
		ret = append(ret, i)
	}
	return ret
}

// Tabulate returns a slice containing values of a given function op over a range of integer values starting from 0 to given n.
func Tabulate[T any](n int, op funcs.Func[int, T]) S[T] {
	if n <= 0 {
		return Empty[T]()
	}

	ret := make([]T, n)
	for i := 0; i < n; i++ {
		ret[i] = op(i)
	}

	return ret
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// IndexFromFunc returns index of the first element is same as given x after or at given start index.
func IndexFromFunc[T, U any](s S[T], x U, start int, eq funcs.Equal[T, U]) int {
	p := funcs.EqualTo(x, eq)
	return s.IndexWhereFrom(p, start)
}

// IndexFunc returns index of the first element is same as given x.
func IndexFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) int {
	return IndexFromFunc(s, x, 0, eq)
}

// IndexFrom returns index of the first element is same as given x after or at given start index.
func IndexFrom[T comparable](s S[T], x T, start int) int {
	return IndexFromFunc(s, x, start, funcs.Same[T])
}

// IndexFrom returns index of the first element is same as given x.
func Index[T comparable](s S[T], x T) int {
	return IndexFrom(s, x, 0)
}

// LastIndexFromFunc returns index of the last element is same as given x before or at given end index.
func LastIndexFromFunc[T, U any](s S[T], x U, end int, eq funcs.Equal[T, U]) int {
	p := funcs.EqualTo(x, eq)

	return s.LastIndexWhereFrom(p, end)
}

// LastIndexFunc returns index of the last element is same as given x.
func LastIndexFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) int {
	return LastIndexFromFunc(s, x, -1, eq)
}

// LastIndexFrom returns index of the last element is same as given x before or at given end index.
func LastIndexFrom[T comparable](s S[T], x T, end int) int {
	return LastIndexFromFunc(s, x, end, funcs.Same[T])
}

// LastIndex returns index of the last element is same as given x.
func LastIndex[T comparable](s S[T], x T) int {
	return LastIndexFrom(s, x, -1)
}

// Contain returns true if s contains given x.
func Contain[T comparable](s S[T], x T) bool {
	return Index(s, x) >= 0
}

// ContainFunc returns true if s contains an element is equal to given x.
func ContainFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) bool {
	return IndexFunc(s, x, eq) >= 0
}

// EqualFunc returns true if s1 is equal to s2.
func EqualFunc[T, U any](s1 S[T], s2 S[U], eq funcs.Equal[T, U]) bool {
	size1 := len(s1)
	size2 := len(s2)

	if size1 != size2 {
		return false
	}

	for i := range s1 {
		if !eq(s1[i], s2[i]) {
			return false
		}
	}
	return true
}

// EqualFunc returns true if s1 is same as s2.
func Equal[T comparable](s1 S[T], s2 S[T]) bool {
	return EqualFunc(s1, s2, funcs.Same[T])
}

// Collect returns a new slice containing results applying given partial function p on which it is defined.
func Collect[T, U any](s S[T], p funcs.Partial[T, U]) S[U] {
	return Fold(
		s,
		Empty[U](),
		func(z S[U], v T) S[U] {
			if u, ok := p(v); ok {
				return append(z, u)
			}
			return z
		},
	)
}

// CollectFirst returns the first result applying given function p successfully.
func CollectFirst[T, U any](s S[T], p funcs.Partial[T, U]) gs.Option[U] {
	for i := range s {
		if u, ok := p(s[i]); ok {
			return gs.Some(u)
		}
	}
	return gs.None[U]()
}

// FoldLeft applies given function op to given start value z and all elements in slice s from left to right.
func FoldLeft[T, U any](s S[T], z U, op func(U, T) U) (ret U) {
	ret = z
	for i := range s {
		ret = op(ret, s[i])
	}

	return
}

// FoldRight applies given function op to given start value z and all elements in slice s from right to left.
func FoldRight[T, U any](s S[T], z U, op func(T, U) U) (ret U) {
	ret = z
	size := len(s)
	for i := size - 1; i >= 0; i-- {
		ret = op(s[i], ret)
	}

	return ret
}

// Fold is same as FoldLeft.
func Fold[T, U any](s S[T], z U, op func(U, T) U) U {
	return FoldLeft(s, z, op)
}

// ScanLeft produces a new slice containing cumulative results of applying the given function op to all elements in slices s from left to right.
func ScanLeft[T, U any](s S[T], z U, op func(U, T) U) S[U] {
	return FoldLeft(s, One(z), func(a S[U], b T) S[U] {
		return append(a, op(a[len(a)-1], b))
	})
}

// ScanRight produces a new slice containing cumulative results of applying the given function op to all elements in slices s from right to left.
func ScanRight[T, U any](s S[T], z U, op func(T, U) U) (ret S[U]) {
	ret = FoldRight(s, One(z), func(a T, b S[U]) S[U] {
		return append(b, op(a, b[len(b)-1]))
	})

	ret.ReverseSelf()
	return
}

// Scan is same as ScanLeft.
func Scan[T, U any](s S[T], z U, op func(U, T) U) S[U] {
	return ScanLeft(s, z, op)
}

// FlatMap returns a new slice s by applying given function op to all elements of slice s.
func FlatMap[T, U any](s S[T], op funcs.Func[T, S[U]]) S[U] {
	return Fold(s, Empty[U](), func(a S[U], v T) S[U] {
		return append(a, op(v)...)
	})
}

// Map returns a new slice by applying given function op to all elements of slices s.
func Map[T, U any](s S[T], op funcs.Func[T, U]) S[U] {
	return Fold(s, Empty[U](), func(a S[U], v T) S[U] {
		return append(a, op(v))
	})
}

// PartitionMap applies given function op to each element of a slice and returns a pair of slices: the first contains Left result from op, and the second one contains Right result from op.
func PartitionMap[T, A, B any](s S[T], op func(T) gs.Either[A, B]) (S[A], S[B]) {
	t2 := Fold(
		s,
		gs.T2(Empty[A](), Empty[B]()),
		func(z gs.Tuple2[S[A], S[B]], v T) gs.Tuple2[S[A], S[B]] {
			e := op(v)
			if e.IsRight() {
				z.V2 = append(z.V2, e.Right())
			} else {
				z.V1 = append(z.V1, e.Left())
			}
			return z
		},
	)
	return t2.V1, t2.V2
}

// GroupMap partitions a slice into a map according to a discriminator function key. All the values that have the same discriminator are then transformed by the function val
func GroupMap[T any, K comparable, V any](s S[T], key funcs.Func[T, K], val funcs.Func[T, V]) map[K]S[V] {

	return Fold(
		s,
		make(map[K]S[V]),
		func(z map[K]S[V], x T) map[K]S[V] {
			k := key(x)
			v := val(x)
			z[k] = append(z[k], v)
			return z
		},
	)
}

// GroupBy partitions a slice into a map of arrays according to the discriminator function key.
func GroupBy[T any, K comparable](s S[T], key funcs.Func[T, K]) map[K]S[T] {
	return GroupMap(s, key, funcs.Self[T])
}

// GroupMapReduce partitions a slice into a map according to a discriminator function key. All the values that have the same discriminator are then transformed by the function val and then reduced into a single value with the reduce function op.
func GroupMapReduce[T any, K comparable, V any](s S[T], key funcs.Func[T, K], val funcs.Func[T, V], op func(V, V) V) map[K]V {
	m := GroupMap(s, key, val)
	ret := make(map[K]V)

	for k := range m {
		ret[k] = m[k].Reduce(op).Get()
	}
	return ret
}

// MaxBy returns Some with the maximum value in s according to the ordered results transformed from given ordering function op.
func MaxBy[T any, R constraints.Ordered](s S[T], op funcs.Orderize[T, R]) gs.Option[T] {
	return s.Max(func(a, b T) int { return funcs.Order(op(a), op(b)) })
}

// MinBy returns Some with minimum value in s according to the ordered results transformed from given ordering function op.
func MinBy[T any, R constraints.Ordered](s S[T], op funcs.Orderize[T, R]) gs.Option[T] {
	return s.Min(func(a, b T) int { return funcs.Order(op(a), op(b)) })
}

// SortBy sorts a slice according to the ordered results transformed from given ordering function op.
func SortBy[T any, R constraints.Ordered](s S[T], op funcs.Orderize[T, R]) S[T] {
	return s.Sort(func(a, b T) int {
		return funcs.Order(op(a), op(b))
	})
}

// IsEmpty returns true if the given slice is empty.
func IsEmpty[T any](s S[T]) bool {
	return len(s) <= 0
}
