package maps

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/slices"
)

func (m M[K, V]) Keys() slices.S[K] {
	return Fold(
		m,
		make(slices.S[K], 0, len(m)),
		func(z slices.S[K], k K, _ V) slices.S[K] {
			return append(z, k)
		},
	)
}

func (m M[K, V]) Values() slices.S[V] {
	return Fold(
		m,
		make(slices.S[V], 0, len(m)),
		func(z slices.S[V], _ K, v V) slices.S[V] {
			return append(z, v)
		},
	)
}

func (m M[K, V]) Contain(x K) (ok bool) {
	_, ok = m[x]
	return
}

func (m M[K, V]) Count(p func(K, V) bool) int {
	return Fold(m, 0, func(a int, k K, v V) int {
		return funcs.Cond(p(k, v), a+1, a)
	})
}

func (m M[K, V]) Find(p func(K, V) bool) gs.Option[Pair[K, V]] {
	for k, v := range m {
		if p(k, v) {
			return gs.Some(P(k, v))
		}
	}
	return gs.None[Pair[K, V]]()
}

func (m M[K, V]) Exists(p func(K, V) bool) bool {
	for k, v := range m {
		if p(k, v) {
			return true
		}
	}
	return false
}

func (m M[K, V]) Filter(p func(K, V) bool) slices.S[Pair[K, V]] {
	return Fold(
		m,
		slices.Empty[Pair[K, V]](),
		func(z slices.S[Pair[K, V]], k K, v V) slices.S[Pair[K, V]] {
			return funcs.Cond(p(k, v), append(z, P(k, v)), z)
		},
	)
}

func (m M[K, V]) FilterNot(p func(K, V) bool) slices.S[Pair[K, V]] {
	return m.Filter(func(k K, v V) bool { return !p(k, v) })
}

func (m M[K, V]) Forall(p func(K, V) bool) bool {
	for k, v := range m {
		if !p(k, v) {
			return false
		}
	}
	return true
}

func (m M[K, V]) Foreach(op func(K, V)) {
	for k, v := range m {
		op(k, v)
	}
}

func (m M[K, V]) Partition(p func(K, V) bool) (_, _ slices.S[Pair[K, V]]) {
	t2 := Fold(
		m,
		gs.T2(slices.Empty[Pair[K, V]](), slices.Empty[Pair[K, V]]()),
		func(z gs.Tuple2[slices.S[Pair[K, V]], slices.S[Pair[K, V]]], k K, v V) gs.Tuple2[slices.S[Pair[K, V]], slices.S[Pair[K, V]]] {
			if p(k, v) {
				z.V2 = append(z.V2, P(k, v))
			} else {
				z.V1 = append(z.V1, P(k, v))
			}
			return z
		},
	)

	return t2.V1, t2.V2
}

func (m M[K, V]) Slice() slices.S[Pair[K, V]] {
	return Fold(
		m,
		make(slices.S[Pair[K, V]], 0, len(m)),
		func(z slices.S[Pair[K, V]], k K, v V) slices.S[Pair[K, V]] {
			return append(z, P(k, v))
		},
	)
}
