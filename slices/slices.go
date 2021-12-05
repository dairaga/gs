package slices

import (
	"constraints"
	"sort"

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

func Clone[T any](s S[T]) S[T] {
	if len(s) <= 0 {
		return Empty[T]()
	}

	return append(Empty[T](), s...)
}

func ReverseSelf[T any](s S[T]) S[T] {
	size := len(s)
	half := size / 2

	for i := 0; i < half; i++ {
		tmp := s[i]
		s[i] = s[size-1-i]
		s[size-1-i] = tmp
	}

	return s
}

func Reverse[T any](s S[T]) (ret S[T]) {
	ret = Clone(s)
	return ReverseSelf(ret)
}

func IndexWhereFrom[T any](s S[T], p funcs.Predict[T], from int) int {
	size := len(s)
	if size <= 0 {
		return -1
	}

	if from < 0 {
		from = size + from
	}

	for i := from; i < size; i++ {
		if p(s[i]) {
			return i
		}
	}
	return -1
}

func IndexWhere[T any](s S[T], p funcs.Predict[T]) int {
	return IndexWhereFrom(s, p, 0)
}

func IndexFromFunc[T, U any](s S[T], x U, from int, eq funcs.Equal[U, T]) int {
	p := func(v T) bool {
		return eq(x, v)
	}

	return IndexWhereFrom(s, p, from)
}

func IndexFunc[T, U any](s S[T], x U, eq funcs.Equal[U, T]) int {
	return IndexFromFunc(s, x, 0, eq)
}

func IndexFrom[T comparable](s S[T], x T, from int) int {
	return IndexFromFunc(s, x, from, func(a, b T) bool { return a == b })
}

func Index[T comparable](s S[T], x T) int {
	return IndexFrom(s, x, 0)
}

func LastIndexWhereFrom[T any](s S[T], p funcs.Predict[T], from int) int {
	size := len(s)
	if size <= 0 {
		return -1
	}

	if from < 0 {
		from = size + from
	}

	from = funcs.Min(from, size-1)

	for i := from; i >= 0; i-- {
		if p(s[i]) {
			return i
		}
	}
	return -1
}

func LastIndexWhere[T any](s S[T], p funcs.Predict[T]) int {
	return LastIndexWhereFrom(s, p, -1)
}

func LastIndexFromFunc[T, U any](s S[T], x U, from int, eq funcs.Equal[U, T]) int {
	p := func(v T) bool {
		return eq(x, v)
	}

	return LastIndexWhereFrom(s, p, from)
}

func LastIndexFunc[T, U any](s S[T], x U, eq funcs.Equal[U, T]) int {
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

func ContainFunc[T, U any](s S[T], x U, eq funcs.Equal[U, T]) bool {
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

func Collect[T, U any](s S[T], p funcs.Check[T, U]) S[U] {
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

func CollectFirst[T, U any](s S[T], p funcs.Check[T, U]) gs.Option[U] {
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

	ReverseSelf(ret)
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
		ret[k] = Reduce(m[k], op).Get()
	}
	return ret
}

func MaxBy[T any, R constraints.Ordered](s S[T], op funcs.Order[T, R]) gs.Option[T] {
	return Max(s, func(a, b T) int { return funcs.Cmp(op(a), op(b)) })
}

func MinBy[T any, R constraints.Ordered](s S[T], op funcs.Order[T, R]) gs.Option[T] {
	return Min(s, func(a, b T) int { return funcs.Cmp(op(a), op(b)) })
}

// -----------------------------------------------------------------------------

func IsEmpty[T any](s S[T]) bool {
	return len(s) <= 0
}

func Forall[T any](s S[T], p funcs.Predict[T]) bool {
	for i := range s {
		if !p(s[i]) {
			return false
		}
	}
	return true
}

func Exists[T any](s S[T], p funcs.Predict[T]) bool {
	for i := range s {
		if p(s[i]) {
			return true
		}
	}
	return false
}

func Foreach[T any](s S[T], op func(T)) {
	for i := range s {
		op(s[i])
	}
}

func Head[T any](s S[T]) gs.Option[T] {
	if IsEmpty(s) {
		return gs.None[T]()
	}
	return gs.Some(s[0])
}

func Heads[T any](s S[T]) S[T] {
	if IsEmpty(s) {
		return s
	}
	return s[:len(s)-1]
}

func Last[T any](s S[T]) gs.Option[T] {
	if IsEmpty(s) {
		return gs.None[T]()
	}
	return gs.Some(s[len(s)-1])
}

func Tail[T any](s S[T]) S[T] {
	if IsEmpty(s) {
		return s
	}
	return s[1:]
}

func Filter[T any](s S[T], p funcs.Predict[T]) S[T] {
	return Fold(s, Empty[T](), func(z S[T], v T) S[T] {
		if p(v) {
			z = append(z, v)
		}
		return z
	})
}

func FilterNot[T any](s S[T], p funcs.Predict[T]) S[T] {
	return Filter(s, func(v T) bool { return !p(v) })
}

func Find[T any](s S[T], p funcs.Predict[T]) gs.Option[T] {
	pos := IndexWhere(s, p)
	if pos >= 0 {
		return gs.Some(s[pos])
	}
	return gs.None[T]()
}

func Partition[T any](s S[T], p funcs.Predict[T]) (_, _ S[T]) {
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

func SplitAt[T any](s S[T], n int) (a, b S[T]) {
	size := len(s)
	if size <= 0 {
		a, b = Empty[T](), Empty[T]()
		return
	}

	if n < 0 {
		n = size + n
	}

	n = funcs.Max(n, size)

	return s[0:n], s[n:]
}

func Count[T any](s S[T], p funcs.Predict[T]) (ret int) {
	for i := range s {
		if p(s[i]) {
			ret++
		}
	}
	return
}

func Take[T any](s S[T], n int) S[T] {
	a, b := SplitAt(s, n)
	if n >= 0 {
		return Clone(a)
	}
	return Clone(b)
}

func TakeWhile[T any](s S[T], p funcs.Predict[T]) (ret S[T]) {
	ret = Empty[T]()

	for i := range s {
		if !p(s[i]) {
			break
		}
		ret = append(ret, s[i])
	}

	return
}

func Drop[T any](s S[T], n int) S[T] {
	a, b := SplitAt(s, n)
	if n >= 0 {
		return Clone(b)
	}
	return Clone(a)
}

func DropWhile[T any](s S[T], p funcs.Predict[T]) S[T] {
	for i := range s {
		if !p(s[i]) {
			return Clone(s[i:])
		}
	}
	return Empty[T]()
}

func ReduceLeft[T any](s S[T], op func(T, T) T) gs.Option[T] {
	head := Head(s)
	if head.IsEmpty() {
		return head
	}

	tail := Tail(s)

	return funcs.Cond(IsEmpty(tail), head, gs.Some(FoldLeft(tail, head.Get(), op)))

}

func ReduceRight[T any](s S[T], op func(T, T) T) gs.Option[T] {
	last := Last(s)
	if last.IsEmpty() {
		return last
	}

	heads := Heads(s)

	return funcs.Cond(
		IsEmpty(heads),
		last,
		gs.Some(FoldRight(heads, last.Get(), op)),
	)

}

func Reduce[T any](s S[T], op func(T, T) T) gs.Option[T] {
	return ReduceLeft(s, op)
}

func Max[T any](s S[T], cmp funcs.Compare[T, T]) gs.Option[T] {
	return Reduce(
		s,
		func(a, b T) T {
			return funcs.Cond(cmp(a, b) >= 0, a, b)
		},
	)
}

func Min[T any](s S[T], cmp funcs.Compare[T, T]) gs.Option[T] {
	return Reduce(
		s,
		func(a, b T) T {
			return funcs.Cond(cmp(a, b) <= 0, a, b)
		},
	)
}

func Sort[T any](s S[T], cmp funcs.Compare[T, T]) S[T] {
	sort.SliceStable(s, func(i, j int) bool { return cmp(s[i], s[j]) < 0 })
	return s
}
