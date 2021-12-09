package maps

import (
	"constraints"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/slices"
)

type M[K comparable, V any] map[K]V

type Pair[K comparable, V any] struct {
	_     struct{}
	Key   K
	Value V
}

func P[K comparable, V any](k K, v V) Pair[K, V] {
	return Pair[K, V]{
		Key:   k,
		Value: v,
	}
}

func From[K comparable, V any](a ...Pair[K, V]) (ret M[K, V]) {
	ret = make(M[K, V], len(a))
	for i := range a {
		ret[a[i].Key] = a[i].Value
	}
	return ret
}

// -----------------------------------------------------------------------------

// TODO: refactor the method when go 1.19 releases.

func Fold[K comparable, V, U any](m M[K, V], z U, op func(U, K, V) U) (ret U) {
	ret = z

	for k, v := range m {
		ret = op(ret, k, v)
	}

	return
}

func Collect[K comparable, V, T any](
	m M[K, V],
	p func(K, V) (T, bool)) slices.S[T] {

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

func CollectMap[K1, K2 comparable, V1, V2 any](
	m M[K1, V1],
	p func(K1, V1) (K2, V2, bool)) M[K2, V2] {

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

func FlatMapSlice[K comparable, V any, T any](
	m M[K, V],
	op func(K, V) slices.S[T]) slices.S[T] {

	return Fold(
		m,
		slices.Empty[T](),
		func(z slices.S[T], k K, v V) slices.S[T] {
			return append(z, op(k, v)...)
		},
	)
}

func FlatMap[K1, K2 comparable, V1, V2 any](m M[K1, V1],
	op func(K1, V1) M[K2, V2]) M[K2, V2] {

	return Fold(
		m,
		make(M[K2, V2]),
		func(z M[K2, V2], k K1, v V1) M[K2, V2] {
			return z.Merge(op(k, v))
		},
	)
}

func MapSlice[K comparable, V, T any](m M[K, V], op func(K, V) T) slices.S[T] {
	return Fold(
		m,
		make(slices.S[T], 0, len(m)),
		func(z slices.S[T], k K, v V) slices.S[T] {
			return append(z, op(k, v))
		},
	)
}

func Map[K1, K2 comparable, V1, V2 any](
	m M[K1, V1],
	op func(K1, V1) (K2, V2)) M[K2, V2] {

	return Fold(
		m,
		make(M[K2, V2]),
		func(z M[K2, V2], k K1, v V1) M[K2, V2] {
			return z.Put(op(k, v))
		},
	)
}

func GroupMap[K1, K2 comparable, V1, V2 any](
	m M[K1, V1],
	key func(K1, V1) K2,
	val func(K1, V1) V2) M[K2, slices.S[V2]] {

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

func GroupBy[K, K1 comparable, V any](
	m M[K, V],
	key func(K, V) K1) M[K1, M[K, V]] {

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

func GroupMapReduce[K1, K2 comparable, V1, V2 any](
	m M[K1, V1],
	key func(K1, V1) K2,
	val func(K1, V1) V2,
	op func(V2, V2) V2) M[K2, V2] {

	return Fold(
		GroupMap(m, key, val),
		make(M[K2, V2]),
		func(z M[K2, V2], k K2, v slices.S[V2]) M[K2, V2] {
			z[k] = v.Reduce(op).Get()
			return z
		},
	)
}

func PartitionMap[K comparable, V, A, B any](
	m M[K, V],
	op func(K, V) gs.Either[A, B]) (slices.S[A], slices.S[B]) {

	t2 := Fold(
		m,
		gs.T2(slices.Empty[A](), slices.Empty[B]()),
		func(z gs.Tuple2[slices.S[A], slices.S[B]], k K, v V) gs.Tuple2[slices.S[A], slices.S[B]] {
			e := op(k, v)
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

func MaxBy[K comparable, V any, B constraints.Ordered](
	m M[K, V],
	op func(K, V) B) gs.Option[Pair[K, V]] {

	return slices.MaxBy(
		m.Slice(),
		func(pair Pair[K, V]) B {
			return op(pair.Key, pair.Value)
		},
	)
}

func MinBy[K comparable, V any, B constraints.Ordered](
	m M[K, V],
	op func(K, V) B) gs.Option[Pair[K, V]] {

	return slices.MinBy(
		m.Slice(),
		func(pair Pair[K, V]) B {
			return op(pair.Key, pair.Value)
		},
	)
}
