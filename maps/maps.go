// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package maps

import (
	"constraints"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/slices"
)

type M[K comparable, V any] map[K]V

// Pair is a 2 dimensional tuple contains key and value in map.
type Pair[K comparable, V any] struct {
	_     struct{}
	Key   K
	Value V
}

// P builds a pair from given key k, and value v.
func P[K comparable, V any](k K, v V) Pair[K, V] {
	return Pair[K, V]{
		Key:   k,
		Value: v,
	}
}

// Form builds a map from pairs.
func From[K comparable, V any](a ...Pair[K, V]) (ret M[K, V]) {
	ret = make(M[K, V], len(a))
	for i := range a {
		ret[a[i].Key] = a[i].Value
	}
	return ret
}

// Zip combines given two slices into a map.
func Zip[K comparable, V any](a slices.S[K], b slices.S[V]) (ret M[K, V]) {
	size := funcs.Min(len(a), len(b))
	ret = make(M[K, V], size)

	for i := 0; i < size; i++ {
		ret[a[i]] = b[i]
	}
	return ret
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// Fold applies given function op to a start value z and all elements of this map.
func Fold[K comparable, V, U any](m M[K, V], z U, op func(U, K, V) U) (ret U) {
	ret = z

	for k, v := range m {
		ret = op(ret, k, v)
	}

	return
}

// Collect returns a new slice by applying a partial function p to all elements of map m on which it is defined. It might return different results for different runs.
func Collect[K comparable, V, T any](m M[K, V], p func(K, V) (T, bool)) slices.S[T] {

	return Fold(
		m,
		slices.Empty[T](),
		func(z slices.S[T], k K, v V) slices.S[T] {
			if val, ok := p(k, v); ok {
				return append(z, val)
			}
			return z
		},
	)
}

// CollectMap returns a new map by applying a partial function p to all elements of map m on which it is defined.
func CollectMap[K1, K2 comparable, V1, V2 any](m M[K1, V1], p func(K1, V1) (K2, V2, bool)) M[K2, V2] {

	return Fold(
		m,
		make(M[K2, V2]),
		func(z M[K2, V2], k K1, v V1) M[K2, V2] {
			if k2, v2, ok := p(k, v); ok {
				z[k2] = v2
			}
			return z
		},
	)
}

// FlatMapSlice returns a new slices by applying given function op to all elements of map m and merge results.
func FlatMapSlice[K comparable, V any, T any](m M[K, V], op func(K, V) slices.S[T]) slices.S[T] {

	return Fold(
		m,
		slices.Empty[T](),
		func(z slices.S[T], k K, v V) slices.S[T] {
			return append(z, op(k, v)...)
		},
	)
}

// FlatMap returns a new map by applying given function op to all elements of map m and merge results.
func FlatMap[K1, K2 comparable, V1, V2 any](m M[K1, V1], op func(K1, V1) M[K2, V2]) M[K2, V2] {
	return Fold(
		m,
		make(M[K2, V2]),
		func(z M[K2, V2], k K1, v V1) M[K2, V2] {
			return z.Merge(op(k, v))
		},
	)
}

// MapSlice retuns a new slice by applying given function op to all elements of map m.
func MapSlice[K comparable, V, T any](m M[K, V], op func(K, V) T) slices.S[T] {
	return Fold(
		m,
		make(slices.S[T], 0, len(m)),
		func(z slices.S[T], k K, v V) slices.S[T] {
			return append(z, op(k, v))
		},
	)
}

// Map retuns a new map by applying given function op to all elements of map m.
func Map[K1, K2 comparable, V1, V2 any](m M[K1, V1], op func(K1, V1) (K2, V2)) M[K2, V2] {
	return Fold(
		m,
		make(M[K2, V2]),
		func(z M[K2, V2], k K1, v V1) M[K2, V2] {
			return z.Put(op(k, v))
		},
	)
}

// GroupMap partitions map m into a map of maps according to a discriminator function key. Each element in a group is transformed into a value of type V2 using function val.
func GroupMap[K1, K2 comparable, V1, V2 any](m M[K1, V1], key func(K1, V1) K2, val func(K1, V1) V2) M[K2, slices.S[V2]] {
	return Fold(
		m,
		make(M[K2, slices.S[V2]]),
		func(z M[K2, slices.S[V2]], k K1, v V1) M[K2, slices.S[V2]] {
			k2 := key(k, v)
			v2 := val(k, v)
			z[k2] = append(z[k2], v2)
			return z
		},
	)
}

// GroupBy partitions map m into a map of maps according to some discriminator function.
func GroupBy[K, K1 comparable, V any](m M[K, V], key func(K, V) K1) M[K1, M[K, V]] {
	return Fold(
		m,
		make(M[K1, M[K, V]]),
		func(z M[K1, M[K, V]], k K, v V) M[K1, M[K, V]] {
			k2 := key(k, v)
			m2, ok := z[k2]
			if !ok {
				m2 = make(M[K, V])
			}
			m2[k] = v
			z[k2] = m2
			return z
		},
	)
}

// GroupMapReduce partitions map m into a map according to a discriminator function key. All the values that have the same discriminator are then transformed by the function val and then reduced into a single value with the reduce function op.
func GroupMapReduce[K1, K2 comparable, V1, V2 any](m M[K1, V1], key func(K1, V1) K2, val func(K1, V1) V2, op func(V2, V2) V2) M[K2, V2] {
	return Fold(
		GroupMap(m, key, val),
		make(M[K2, V2]),
		func(z M[K2, V2], k K2, v slices.S[V2]) M[K2, V2] {
			z[k] = v.Reduce(op).Get()
			return z
		},
	)
}

// PartitionMap applies given function op to each element of the map and returns a pair of maps: the first one made of those values returned by f that were wrapped in scala.util.Left, and the second one made of those wrapped in scala.util.Right.
func PartitionMap[K comparable, V, A, B any](m M[K, V], op func(K, V) gs.Either[A, B]) (M[K, A], M[K, B]) {

	t2 := Fold(
		m,
		gs.T2(make(M[K, A]), make(M[K, B])),
		func(z gs.Tuple2[M[K, A], M[K, B]], k K, v V) gs.Tuple2[M[K, A], M[K, B]] {
			e := op(k, v)
			if e.IsRight() {
				z.V2[k] = e.Right()
			} else {
				z.V1[k] = e.Left()
			}
			return z
		},
	)

	return t2.V1, t2.V2
}

// MaxBy returns a maximum pair of key and value according to result of ordering function op.
func MaxBy[K comparable, V any, B constraints.Ordered](m M[K, V], op func(K, V) B) gs.Option[Pair[K, V]] {
	return slices.MaxBy(
		m.Slice(),
		func(pair Pair[K, V]) B {
			return op(pair.Key, pair.Value)
		},
	)
}

// MinBy returns a minimum pair of key and value according to result of ordering function op.
func MinBy[K comparable, V any, B constraints.Ordered](m M[K, V], op func(K, V) B) gs.Option[Pair[K, V]] {
	return slices.MinBy(
		m.Slice(),
		func(pair Pair[K, V]) B {
			return op(pair.Key, pair.Value)
		},
	)
}
