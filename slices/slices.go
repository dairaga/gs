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

func Empty[T any]() S[T] {
	return []T{}
}

func One[T any](v T) S[T] {
	return []T{v}
}

func From[T any](a ...T) S[T] {
	return a
}

func Fill[T any](size int, v T) S[T] {
	ret := make([]T, size)
	for i := range ret {
		ret[i] = v
	}
	return ret
}

func Range[T gs.Numeric](start, end, step T) S[T] {
	ret := Empty[T]()

	for i := start; i < end; i += step {
		ret = append(ret, i)
	}
	return ret
}

func Tabulate[T any](size int, op funcs.Func[int, T]) S[T] {
	if size <= 0 {
		return Empty[T]()
	}

	ret := make([]T, size)
	for i := 0; i < size; i++ {
		ret[i] = op(i)
	}

	return ret
}

// -----------------------------------------------------------------------------

// TODO: refactor the method when go 1.19 releases.

func IndexFromFunc[T, U any](s S[T], x U, from int, eq funcs.Equal[T, U]) int {
	p := func(v T) bool {
		return eq(v, x)
	}

	return s.IndexWhereFrom(p, from)
}

func IndexFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) int {
	return IndexFromFunc(s, x, 0, eq)
}

func IndexFrom[T comparable](s S[T], x T, from int) int {
	return IndexFromFunc(s, x, from, func(a, b T) bool { return a == b })
}

func Index[T comparable](s S[T], x T) int {
	return IndexFrom(s, x, 0)
}

func LastIndexFromFunc[T, U any](s S[T], x U, from int, eq funcs.Equal[T, U]) int {
	p := func(v T) bool {
		return eq(v, x)
	}

	return s.LastIndexWhereFrom(p, from)
}

func LastIndexFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) int {
	return LastIndexFromFunc(s, x, -1, eq)
}

func LastIndexFrom[T comparable](s S[T], x T, from int) int {
	return LastIndexFromFunc(s, x, from, func(a, b T) bool { return a == b })
}

func LastIndex[T comparable](s S[T], x T) int {
	return LastIndexFrom(s, x, -1)
}

func Contain[T comparable](s S[T], x T) bool {
	return Index(s, x) >= 0
}

func ContainFunc[T, U any](s S[T], x U, eq funcs.Equal[T, U]) bool {
	return IndexFunc(s, x, eq) >= 0
}

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

func Equal[T comparable](s1 S[T], s2 S[T]) bool {
	return EqualFunc(s1, s2, func(a, b T) bool { return a == b })
}

func Collect[T, U any](s S[T], p funcs.Can[T, U]) S[U] {
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

func CollectFirst[T, U any](s S[T], p funcs.Can[T, U]) gs.Option[U] {
	for i := range s {
		if u, ok := p(s[i]); ok {
			return gs.Some(u)
		}
	}
	return gs.None[U]()
}

func FoldLeft[T, U any](s S[T], z U, op func(U, T) U) (ret U) {
	ret = z
	for i := range s {
		ret = op(ret, s[i])
	}

	return
}

func FoldRight[T, U any](s S[T], z U, op func(T, U) U) (ret U) {
	ret = z
	size := len(s)
	for i := size - 1; i >= 0; i-- {
		ret = op(s[i], ret)
	}

	return ret
}

func Fold[T, U any](s S[T], z U, op func(U, T) U) U {
	return FoldLeft(s, z, op)
}

func ScanLeft[T, U any](s S[T], z U, op func(U, T) U) S[U] {
	return FoldLeft(s, One(z), func(a S[U], b T) S[U] {
		return append(a, op(a[len(a)-1], b))
	})
}

func ScanRight[T, U any](s S[T], z U, op func(T, U) U) (ret S[U]) {
	ret = FoldRight(s, One(z), func(a T, b S[U]) S[U] {
		return append(b, op(a, b[len(b)-1]))
	})

	ret.ReverseSelf()
	return
}

func Scan[T, U any](s S[T], z U, op func(U, T) U) S[U] {
	return ScanLeft(s, z, op)
}

func FlatMap[T, U any](s S[T], op funcs.Func[T, S[U]]) S[U] {
	return Fold(s, Empty[U](), func(a S[U], v T) S[U] {
		return append(a, op(v)...)
	})
}

func Map[T, U any](s S[T], op funcs.Func[T, U]) S[U] {
	return Fold(s, Empty[U](), func(a S[U], v T) S[U] {
		return append(a, op(v))
	})
}

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

func GroupBy[T any, K comparable](s S[T], key funcs.Func[T, K]) map[K]S[T] {
	return GroupMap(s, key, funcs.Self[T])
}

func GroupMapReduce[T any, K comparable, V any](s S[T], key funcs.Func[T, K], val funcs.Func[T, V], op func(V, V) V) map[K]V {
	// TODO: refactor when go 1.19 release
	// use fold of map
	m := GroupMap(s, key, val)
	ret := make(map[K]V)

	for k := range m {
		ret[k] = m[k].Reduce(op).Get()
	}
	return ret
}

func MaxBy[T any, R constraints.Ordered](s S[T], op funcs.Order[T, R]) gs.Option[T] {
	return s.Max(func(a, b T) int { return funcs.Cmp(op(a), op(b)) })
}

func MinBy[T any, R constraints.Ordered](s S[T], op funcs.Order[T, R]) gs.Option[T] {
	return s.Min(func(a, b T) int { return funcs.Cmp(op(a), op(b)) })
}

func SortBy[T any, R constraints.Ordered](s S[T], op funcs.Order[T, R]) S[T] {
	return s.Sort(func(a, b T) int {
		return funcs.Cmp(op(a), op(b))
	})
}

func IsEmpty[T any](s S[T]) bool {
	return len(s) <= 0
}
