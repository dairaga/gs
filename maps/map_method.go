// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package maps

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/slices"
)

// Keys returns a slices of all keys.
func (m M[K, V]) Keys() slices.S[K] {
	return Fold(
		m,
		make(slices.S[K], 0, len(m)),
		func(z slices.S[K], k K, _ V) slices.S[K] {
			return append(z, k)
		},
	)
}

// Values returns a slice of all values.
func (m M[K, V]) Values() slices.S[V] {
	return Fold(
		m,
		make(slices.S[V], 0, len(m)),
		func(z slices.S[V], _ K, v V) slices.S[V] {
			return append(z, v)
		},
	)
}

// Add adds pairs into map.
func (m M[K, V]) Add(pairs ...Pair[K, V]) M[K, V] {
	for _, p := range pairs {
		m[p.Key] = p.Value
	}
	return m
}

// Put add key and value into map.
func (m M[K, V]) Put(key K, val V) M[K, V] {
	m[key] = val
	return m
}

// Merge merges another map a into this. Values in this maybe overwritten by values in a if keys are in a and m.
func (m M[K, V]) Merge(a M[K, V]) M[K, V] {
	for k, v := range a {
		m[k] = v
	}
	return m
}

// Contain returns true if m has given key x.
func (m M[K, V]) Contain(x K) (ok bool) {
	_, ok = m[x]
	return
}

// Count returns numbers of elements in m satisfying given function p.
func (m M[K, V]) Count(p func(K, V) bool) int {
	return Fold(m, 0, func(a int, k K, v V) int {
		return funcs.Cond(p(k, v), a+1, a)
	})
}

// Find returns the first key-value pair of m satisfying given function p. It might return different results for different runs.
func (m M[K, V]) Find(p func(K, V) bool) gs.Option[Pair[K, V]] {
	for k, v := range m {
		if p(k, v) {
			return gs.Some(P(k, v))
		}
	}
	return gs.None[Pair[K, V]]()
}

// Exists return true if at least one element in m satisfies given function p.
func (m M[K, V]) Exists(p func(K, V) bool) bool {
	for k, v := range m {
		if p(k, v) {
			return true
		}
	}
	return false
}

// Filter returns a new map made of elements in m satisfying given function p.
func (m M[K, V]) Filter(p func(K, V) bool) M[K, V] {
	return Fold(
		m,
		make(M[K, V]),
		func(z M[K, V], k K, v V) M[K, V] {
			if p(k, v) {
				z[k] = v
			}
			return z
		},
	)
}

// Filter returns a new map made of elements in m not satisfying given function p.
func (m M[K, V]) FilterNot(p func(K, V) bool) M[K, V] {
	return m.Filter(func(k K, v V) bool { return !p(k, v) })
}

// Forall returns true if this is a empty map or all elements satisfy given function p.
func (m M[K, V]) Forall(p func(K, V) bool) bool {
	for k, v := range m {
		if !p(k, v) {
			return false
		}
	}
	return true
}

// Foreach applies given function op to each element in this.
func (m M[K, V]) Foreach(op func(K, V)) {
	for k, v := range m {
		op(k, v)
	}
}

// Partition partitions this into two maps according to given function p. The first map made of elements in m not satisfying the function p, and the second map made of elements satisfying the function p.
func (m M[K, V]) Partition(p func(K, V) bool) (_, _ M[K, V]) {
	t2 := Fold(
		m,
		gs.T2(make(M[K, V]), make(M[K, V])),
		func(
			z gs.Tuple2[M[K, V], M[K, V]],
			k K,
			v V) gs.Tuple2[M[K, V], M[K, V]] {

			if p(k, v) {
				z.V2[k] = v
			} else {
				z.V1[k] = v
			}
			return z
		},
	)

	return t2.V1, t2.V2
}

// Slice returns a slice containing key-value pairs from this.
func (m M[K, V]) Slice() slices.S[Pair[K, V]] {
	return Fold(
		m,
		make(slices.S[Pair[K, V]], 0, len(m)),
		func(z slices.S[Pair[K, V]], k K, v V) slices.S[Pair[K, V]] {
			return append(z, P(k, v))
		},
	)
}
