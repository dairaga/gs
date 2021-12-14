// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package slices

import (
	"sort"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// IsEmpty returns true if this is an empty slice.
func (s S[T]) IsEmpty() bool {
	return len(s) <= 0
}

// Clone returns a new copy from this.
func (s S[T]) Clone() S[T] {
	if len(s) <= 0 {
		return Empty[T]()
	}

	return append(Empty[T](), s...)
}

// ReverseSelf reverses this slice.
func (s S[T]) ReverseSelf() S[T] {
	size := len(s)
	half := size / 2

	for i := 0; i < half; i++ {
		tmp := s[i]
		s[i] = s[size-1-i]
		s[size-1-i] = tmp
	}

	return s
}

// Reverse returns a new reversed slice.
func (s S[T]) Reverse() S[T] {
	ret := s.Clone()
	return ret.ReverseSelf()
}

// IndexWhereFrom returns index of the first element that satisfies given function p after or at given start index.
func (s S[T]) IndexWhereFrom(p funcs.Predict[T], start int) int {
	size := len(s)
	if size <= 0 {
		return -1
	}

	if start < 0 {
		start = size + start
	}

	start = funcs.Min(start, size-1)
	start = funcs.Max(0, start)

	for i := start; i < size; i++ {
		if p(s[i]) {
			return i
		}
	}
	return -1
}

// IndexWhere returns index of the first element that satisfies given function p.
func (s S[T]) IndexWhere(p funcs.Predict[T]) int {
	return s.IndexWhereFrom(p, 0)
}

// LastIndexWhereFrom returns index of the last element that satisfies given function p before or at given end index.
func (s S[T]) LastIndexWhereFrom(p funcs.Predict[T], end int) int {
	size := len(s)
	if size <= 0 {
		return -1
	}

	if end < 0 {
		end = size + end
	}

	end = funcs.Min(end, size-1)
	end = funcs.Max(0, end)

	for i := end; i >= 0; i-- {
		if p(s[i]) {
			return i
		}
	}
	return -1
}

// LastIndexWhere returns index of last element satisfying given function p.
func (s S[T]) LastIndexWhere(p funcs.Predict[T]) int {
	return s.LastIndexWhereFrom(p, -1)
}

func (s S[T]) Forall(p func(int, T) bool) bool {
	for i := range s {
		if !p(i, s[i]) {
			return false
		}
	}
	return true
}

func (s S[T]) Exists(p func(int, T) bool) bool {
	for i := range s {
		if p(i, s[i]) {
			return true
		}
	}
	return false
}

func (s S[T]) Foreach(op func(int, T)) {
	for i := range s {
		op(i, s[i])
	}
}

func (s S[T]) Head() gs.Option[T] {
	if s.IsEmpty() {
		return gs.None[T]()
	}
	return gs.Some(s[0])
}

func (s S[T]) Heads() S[T] {
	if s.IsEmpty() {
		return s
	}
	return s[:len(s)-1]
}

func (s S[T]) Last() gs.Option[T] {
	if s.IsEmpty() {
		return gs.None[T]()
	}
	return gs.Some(s[len(s)-1])
}

func (s S[T]) Tail() S[T] {
	if s.IsEmpty() {
		return s
	}
	return s[1:]
}

func (s S[T]) Filter(p funcs.Predict[T]) S[T] {
	return Fold(s, Empty[T](), func(z S[T], v T) S[T] {
		if p(v) {
			z = append(z, v)
		}
		return z
	})
}

func (s S[T]) FilterNot(p funcs.Predict[T]) S[T] {
	return s.Filter(func(v T) bool { return !p(v) })
}

func (s S[T]) FindFrom(p funcs.Predict[T], from int) gs.Option[T] {
	pos := s.IndexWhereFrom(p, from)
	if pos >= 0 {
		return gs.Some(s[pos])
	}
	return gs.None[T]()
}

func (s S[T]) Find(p funcs.Predict[T]) gs.Option[T] {
	return s.FindFrom(p, 0)
}

func (s S[T]) FindLastFrom(p funcs.Predict[T], from int) gs.Option[T] {
	pos := s.LastIndexWhereFrom(p, from)
	if pos >= 0 {
		return gs.Some(s[pos])
	}
	return gs.None[T]()
}

func (s S[T]) FindLast(p funcs.Predict[T]) gs.Option[T] {
	return s.FindLastFrom(p, -1)
}

func (s S[T]) Partition(p funcs.Predict[T]) (_, _ S[T]) {
	t2 := Fold(
		s,
		gs.T2(Empty[T](), Empty[T]()),
		func(z gs.Tuple2[S[T], S[T]], x T) gs.Tuple2[S[T], S[T]] {
			if p(x) {
				z.V2 = append(z.V2, x)
			} else {
				z.V1 = append(z.V1, x)
			}
			return z
		},
	)
	return t2.V1, t2.V2
}

func (s S[T]) SplitAt(n int) (a, b S[T]) {
	size := len(s)
	if size <= 0 {
		a, b = Empty[T](), Empty[T]()
		return
	}

	if n < 0 {
		n = size + n
	}

	n = funcs.Min(n, size)
	n = funcs.Max(0, n)

	return s[0:n], s[n:]
}

func (s S[T]) Count(p funcs.Predict[T]) (ret int) {
	for i := range s {
		if p(s[i]) {
			ret++
		}
	}
	return
}

func (s S[T]) Take(n int) S[T] {
	a, b := s.SplitAt(n)
	if n >= 0 {
		return a.Clone()
	}
	return b.Clone()
}

func (s S[T]) TakeWhile(p funcs.Predict[T]) (ret S[T]) {
	ret = Empty[T]()

	for i := range s {
		if !p(s[i]) {
			break
		}
		ret = append(ret, s[i])
	}

	return
}

func (s S[T]) Drop(n int) S[T] {
	a, b := s.SplitAt(n)
	if n >= 0 {
		return b.Clone()
	}
	return a.Clone()
}

func (s S[T]) DropWhile(p funcs.Predict[T]) S[T] {
	for i := range s {
		if !p(s[i]) {
			return S[T](s[i:]).Clone()
		}
	}
	return Empty[T]()
}

func (s S[T]) ReduceLeft(op func(T, T) T) gs.Option[T] {
	head := s.Head()
	if head.IsEmpty() {
		return head
	}

	tail := s.Tail()

	return funcs.Cond(IsEmpty(tail), head, gs.Some(FoldLeft(tail, head.Get(), op)))

}

func (s S[T]) ReduceRight(op func(T, T) T) gs.Option[T] {
	last := s.Last()
	if last.IsEmpty() {
		return last
	}

	heads := s.Heads()

	return funcs.Cond(
		IsEmpty(heads),
		last,
		gs.Some(FoldRight(heads, last.Get(), op)),
	)

}

func (s S[T]) Reduce(op func(T, T) T) gs.Option[T] {
	return s.ReduceLeft(op)
}

func (s S[T]) Max(cmp funcs.Compare[T, T]) gs.Option[T] {
	return s.Reduce(
		func(a, b T) T {
			return funcs.Cond(cmp(a, b) >= 0, a, b)
		},
	)
}

func (s S[T]) Min(cmp funcs.Compare[T, T]) gs.Option[T] {
	return s.Reduce(
		func(a, b T) T {
			return funcs.Cond(cmp(a, b) <= 0, a, b)
		},
	)
}

func (s S[T]) Sort(cmp funcs.Compare[T, T]) S[T] {
	sort.SliceStable(s, func(i, j int) bool { return cmp(s[i], s[j]) < 0 })
	return s
}
