// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"strconv"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/slices"
	"github.com/stretchr/testify/assert"
)

var (
	even = func(v int) bool { return (v & 0x01) == 0 }
	odd  = func(v int) bool { return (v & 0x01) == 1 }
)

type person struct {
	age int
}

func orderize(p person) int {
	return p.age
}

func personEq(v int, p person) bool { return v == p.age }

func TestSliceIsEmpty(t *testing.T) {
	assert.False(t, slices.From(1, 2, 3).IsEmpty())
	assert.True(t, slices.Empty[int]().IsEmpty())
}

func TestSliceClone(t *testing.T) {
	src := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	dst := src.Clone()

	assert.False(t, &src[0] == &dst[0])
	assert.Equal(t, src, dst)
}

func TestSliceReverseSelf(t *testing.T) {
	s := slices.From(1, 2, 3)
	s1 := s.ReverseSelf()

	assert.Equal(t, slices.From(3, 2, 1), s1)
	assert.Equal(t, slices.From(3, 2, 1), s)
	assert.True(t, &s1[0] == &s[0])
}

func TestSliceReverse(t *testing.T) {
	s := slices.From(1, 2, 3)
	s1 := s.Reverse()

	assert.Equal(t, slices.From(3, 2, 1), s1)
	assert.Equal(t, slices.From(1, 2, 3), s)
}

func TestSliceIndexWhereFrom(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	assert.Equal(t, 5, s.IndexWhereFrom(even, 0))
	assert.Equal(t, 5, s.IndexWhere(even))

	assert.Equal(t, 2, s.IndexWhereFrom(odd, 2))
	assert.Equal(t, -1, s.IndexWhereFrom(odd, 5))
	assert.Equal(t, 7, s.IndexWhereFrom(even, -2))
	assert.Equal(t, -1, s.IndexWhereFrom(odd, -2))
}

func TestSliceLastIndexWhereFrom(t *testing.T) {
	assert.Equal(t, -1, slices.Empty[int]().LastIndexWhereFrom(odd, 0))
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	assert.Equal(t, 4, s.LastIndexWhereFrom(odd, -1))
	assert.Equal(t, 4, s.LastIndexWhere(odd))

	assert.Equal(t, -1, s.LastIndexWhereFrom(even, 4))
	assert.Equal(t, -1, s.LastIndexWhereFrom(even, -5))
}

func TestSliceForall(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	p1 := func(_ int, v int) bool {
		return v >= 0
	}

	p2 := func(_ int, v int) bool {
		return v > 5
	}
	assert.True(t, slices.Empty[int]().Forall(p1))
	assert.True(t, slices.Empty[int]().Forall(p2))

	assert.True(t, s.Forall(p1))
	assert.False(t, s.Forall(p2))
}

func TestSliceExists(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	p1 := func(_ int, v int) bool {
		return v > 0
	}
	p2 := func(_ int, v int) bool {
		return v < 0
	}

	assert.False(t, slices.Empty[int]().Exists(p1))
	assert.False(t, slices.Empty[int]().Exists(p2))

	assert.True(t, s.Exists(p1))
	assert.False(t, s.Exists(p2))

}

func TestSliceForeach(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	sum := 0
	s.Foreach(func(_int, v int) {
		sum += v
	})
	assert.Equal(t, 45, sum)
}

func TestSliceHead(t *testing.T) {
	assert.False(t, slices.Empty[int]().Head().IsDefined())
	assert.Equal(t, 1, slices.From(1, 2, 3).Head().Get())
}

func TestSliceHeads(t *testing.T) {
	assert.Equal(t, slices.Empty[int](), slices.Empty[int]().Heads())
	assert.Equal(t, slices.From(1, 2), slices.From(1, 2, 3).Heads())
}

func TestSliceLast(t *testing.T) {
	assert.False(t, slices.Empty[int]().Head().IsDefined())
	assert.Equal(t, 3, slices.From(1, 2, 3).Last().Get())
}

func TestSliceTail(t *testing.T) {
	assert.Equal(t, slices.Empty[int](), slices.Empty[int]().Tail())
	assert.Equal(t, slices.From(2, 3), slices.From(1, 2, 3).Tail())
}

func TestSliceFilter(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	assert.Equal(t, slices.From(2, 4, 6, 8), s.Filter(even))
}

func TestSliceFilterNot(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	assert.Equal(t, slices.From(1, 3, 5, 7, 9), s.FilterNot(even))
}

func TestSliceFind(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	p1 := func(v int) bool {
		return v > 5
	}
	p2 := func(v int) bool {
		return v < 0
	}

	assert.Equal(t, 7, s.Find(p1).Get())
	assert.False(t, s.Find(p2).IsDefined())
	assert.Equal(t, 6, s.FindFrom(p1, 5).Get())
}

func TestSliceFindLast(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	p1 := func(v int) bool {
		return v > 5
	}
	p2 := func(v int) bool {
		return v < 0
	}

	assert.Equal(t, 8, s.FindLast(p1).Get())
	assert.False(t, s.FindLast(p2).IsDefined())
	assert.Equal(t, 9, s.FindLastFrom(p1, 5).Get())
}

func TestSlicePartition(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	a, b := s.Partition(even)

	assert.Equal(t, slices.From(1, 3, 5, 7, 9), a)
	assert.Equal(t, slices.From(2, 4, 6, 8), b)
}

func TestSliceSplitAt(t *testing.T) {
	a, b := slices.Empty[int]().SplitAt(10)
	assert.Equal(t, slices.Empty[int](), a)
	assert.Equal(t, slices.Empty[int](), b)

	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	a, b = s.SplitAt(100)

	assert.Equal(t, s, a)
	assert.Equal(t, slices.Empty[int](), b)

	a, b = s.SplitAt(-100)

	assert.Equal(t, slices.Empty[int](), a)
	assert.Equal(t, s, b)

	a, b = s.SplitAt(5)

	assert.Equal(t, s[0:5], a)
	assert.Equal(t, s[5:], b)

	a, b = s.SplitAt(-3)
	idx := len(s) - 3
	assert.Equal(t, s[:idx], a)
	assert.Equal(t, s[idx:], b)

}
func TestSliceCount(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	assert.Equal(t, 5, s.Count(odd))
}

func TestSliceTake(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	a := s.Take(100)
	assert.Equal(t, s, a)
	assert.False(t, &s[0] == &a[0])

	a = s.Take(-100)
	assert.Equal(t, s, a)
	assert.False(t, &s[0] == &a[0])

	a = s.Take(5)
	ans := s[:5]
	assert.Equal(t, ans, a)
	assert.False(t, &ans[0] == &a[0]) // test clone

	a = s.Take(-3)
	idx := len(s) - 3
	ans = s[idx:]
	assert.Equal(t, ans, a)
	assert.False(t, &ans[0] == &a[0]) // test clone
}

func TestSliceTakeWhile(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	a := s.TakeWhile(odd)
	assert.Equal(t, slices.From(1, 3, 5, 7, 9), a)
}

func TestSliceDrop(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	b := s.Drop(100)
	assert.Equal(t, slices.Empty[int](), b)

	b = s.Drop(-100)
	assert.Equal(t, slices.Empty[int](), b)

	b = s.Drop(5)
	ans := s[5:]
	assert.Equal(t, ans, b)
	assert.False(t, &ans[0] == &b[0]) // test clone

	b = s.Drop(-3)
	idx := len(s) - 3
	ans = s[:idx]
	assert.Equal(t, ans, b)
	assert.False(t, &ans[0] == &b[0]) // test clone

}

func TestSliceDropWhile(t *testing.T) {
	assert.Empty(t, slices.Empty[int]().DropWhile(odd))

	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	a := s.DropWhile(odd)
	assert.Equal(t, slices.From(2, 4, 6, 8), a)
}

func TestSliceReduceLeft(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	op1 := func(v1, v2 int) int {
		return v1 - v2
	}
	op2 := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.False(t, slices.Empty[int]().ReduceLeft(op1).IsDefined())
	assert.False(t, slices.Empty[int]().ReduceLeft(op2).IsDefined())

	assert.Equal(t, -43, s.ReduceLeft(op1).Get())
	assert.Equal(t, -43, s.Reduce(op1).Get())

	assert.Equal(t, 45, s.ReduceLeft(op2).Get())
	assert.Equal(t, 45, s.Reduce(op2).Get())

}

func TestSliceReduceRight(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	op1 := func(v1, v2 int) int {
		return v1 - v2
	}
	op2 := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.False(t, slices.Empty[int]().ReduceRight(op1).IsDefined())
	assert.False(t, slices.Empty[int]().ReduceRight(op2).IsDefined())

	assert.Equal(t, 9, s.ReduceRight(op1).Get())
	assert.Equal(t, 45, s.ReduceRight(op2).Get())
}

func TestSliceMax(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	assert.Equal(t, 9, s.Max(funcs.Order[int]).Get())
	assert.False(t, slices.Empty[int]().Max(funcs.Order[int]).IsDefined())
}

func TestSliceMin(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	assert.Equal(t, 1, s.Min(funcs.Order[int]).Get())
	assert.False(t, slices.Empty[int]().Max(funcs.Order[int]).IsDefined())
}

func TestSliceSort(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	assert.Equal(t,
		slices.From(1, 2, 3, 4, 5, 6, 7, 8, 9),
		s.Sort(funcs.Order[int]),
	)
}

func TestEmpty(t *testing.T) {
	assert.Equal(t, 0, len(slices.Empty[int]()))
}

func TestOne(t *testing.T) {
	s := slices.One(0)
	assert.Equal(t, 1, len(s))
	assert.Equal(t, 0, s[0])
}

func TestFrom(t *testing.T) {
	s := slices.From(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, []int(s))
}

func TestFill(t *testing.T) {
	s := slices.Fill(3, 1)
	assert.Equal(t, slices.From(1, 1, 1), s)
}

func TestRange(t *testing.T) {
	assert.Equal(t,
		slices.From(0, 2, 4, 6, 8),
		slices.Range(0, 10, 2))
}

func TestTabulate(t *testing.T) {

	op := func(v int) string {
		return strconv.Itoa(v + 1)
	}
	assert.Equal(t, slices.From("1", "2", "3"), slices.Tabulate(3, op))
	assert.Empty(t, slices.Tabulate(-1, op))
}

func TestIndex(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 1, 3, 5, 7, 9)

	assert.Equal(t, 4, slices.IndexFromFunc(s, person{9}, 0, personEq))
	assert.Equal(t, 4, slices.IndexFunc(s, person{9}, personEq))
	assert.Equal(t, 4, slices.IndexFrom(s, 9, 0))
	assert.Equal(t, 4, slices.Index(s, 9))

	assert.Equal(t, -1, slices.IndexFromFunc(s, person{3}, -3, personEq))
	assert.Equal(t, -1, slices.IndexFrom(s, 3, -3))
	assert.Equal(t, 8, slices.IndexFrom(s, 7, -3))
	assert.Equal(t, -1, slices.IndexFrom(s, 1, 6))

}

func TestLastIndexFromFunc(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 1, 3, 5, 7, 9)

	assert.Equal(t, 9, slices.LastIndexFromFunc(s, person{9}, -1, personEq))
	assert.Equal(t, 9, slices.LastIndexFunc(s, person{9}, personEq))
	assert.Equal(t, 9, slices.LastIndexFrom(s, 9, -1))
	assert.Equal(t, 9, slices.LastIndex(s, 9))

	assert.Equal(t, -1, slices.LastIndexFromFunc(s, person{9}, -7, personEq))
	assert.Equal(t, -1, slices.LastIndexFrom(s, 9, -7))
	assert.Equal(t, 3, slices.LastIndexFrom(s, 7, -3))
	assert.Equal(t, -1, slices.LastIndexFrom(s, 9, 3))
}

func TestContain(t *testing.T) {
	s := slices.From(1, 2, 3)
	assert.True(t, slices.Contain(s, 1))
	assert.False(t, slices.Contain(s, -1))
}

func TestContainFunc(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	assert.True(t, slices.ContainFunc(s, person{9}, personEq))
	assert.False(t, slices.ContainFunc(s, person{10}, personEq))

	assert.False(t,
		slices.ContainFunc(slices.Empty[int](), person{10}, personEq))

}

func TestEqualFunc(t *testing.T) {
	e1 := slices.Empty[int]()
	e2 := slices.Empty[person]()
	assert.True(t, slices.EqualFunc(e1, e2, personEq))

	s1 := slices.From(1, 2, 3)
	s2 := slices.From(person{1}, person{2}, person{3})
	assert.True(t, slices.EqualFunc(s1, s2, personEq))

	s1 = slices.From(1, 2, 3, 4)
	s2 = slices.From(person{1}, person{2}, person{3})
	assert.False(t, slices.EqualFunc(s1, s2, personEq))

	s1 = slices.From(1, 3, 2)
	s2 = slices.From(person{1}, person{2}, person{3})
	assert.False(t, slices.EqualFunc(s1, s2, personEq))

}

func TestEqual(t *testing.T) {
	e1 := slices.Empty[int]()
	e2 := slices.Empty[int]()
	assert.True(t, slices.Equal(e1, e2))

	s1 := slices.From(1, 2, 3)
	s2 := slices.From(1, 2, 3)

	assert.False(t, &s1[0] == &s2[0])
	assert.True(t, slices.Equal(s1, s2))

	s1 = slices.From(1, 2, 3, 4)
	s2 = slices.From(1, 2, 3)

	assert.False(t, &s1[0] == &s2[0])
	assert.False(t, slices.Equal(s1, s2))

	s1 = slices.From(1, 3, 2)
	s2 = slices.From(1, 2, 3)

	assert.False(t, &s1[0] == &s2[0])
	assert.False(t, slices.Equal(s1, s2))
}

func TestCollect(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	fn := func(v int) (s string, ok bool) {
		if ok = even(v); ok {
			s = strconv.Itoa(v)
		}
		return
	}

	dst := slices.Collect(s, fn)
	assert.Equal(t, slices.From("2", "4", "6", "8"), dst)
}

func TestCollectFirst(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	op := func(v int) (s string, ok bool) {
		if ok = even(v); ok {
			s = strconv.Itoa(v)
		}
		return
	}

	dst := slices.CollectFirst(s, op)
	assert.Equal(t, "2", dst.Get())

	op = func(_ int) (string, bool) {
		return "", false
	}

	assert.False(t, slices.CollectFirst(s, op).IsDefined())
}

func TestFoldLeft(t *testing.T) {
	var src = slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	fn1 := func(v1, v2 int) int {
		return v1 - v2
	}

	fn2 := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.Equal(t, -45, slices.FoldLeft(src, 0, fn1))
	assert.Equal(t, -45, slices.Fold(src, 0, fn1))

	assert.Equal(t, 45, slices.FoldLeft(src, 0, fn2))
	assert.Equal(t, 45, slices.Fold(src, 0, fn2))
}

func TestFoldRight(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	op1 := func(v1, v2 int) int {
		return v1 - v2
	}

	op2 := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.Equal(t, 9, slices.FoldRight(s, 0, op1))
	assert.Equal(t, 45, slices.FoldRight(s, 0, op2))
}

func TestScanLeft(t *testing.T) {
	op := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.Equal(t,
		slices.One(100),
		slices.ScanLeft(slices.Empty[int](), 100, op),
	)

	var src = slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	dst := slices.ScanLeft(src, 100, op)
	ans := slices.From(100, 101, 104, 109, 116, 125, 127, 131, 137, 145)
	assert.Equal(t, ans, dst)

	dst = slices.Scan(src, 100, op)
	assert.Equal(t, ans, dst)
}

func TestScanRight(t *testing.T) {

	op := func(v1, v2 int) int {
		return v1 + v2
	}

	assert.Equal(t,
		slices.One(100),
		slices.ScanRight(slices.Empty[int](), 100, op),
	)

	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	dst := slices.ScanRight(s, 100, op)
	ans := slices.From(145, 144, 141, 136, 129, 120, 118, 114, 108, 100)
	assert.Equal(t, ans, dst)
}

func TestFlatMap(t *testing.T) {
	dst := slices.FlatMap(
		slices.From(1, 2, 3),
		func(v int) slices.S[int] {
			return slices.Map(
				slices.From(4, 5, 6),
				func(x int) int {
					return v * x
				})
		})

	ans := slices.From(4, 5, 6, 8, 10, 12, 12, 15, 18)
	assert.Equal(t, ans, dst)
}

func TestPartitionMap(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	op := func(v int) gs.Either[int, int] {
		if even(v) {
			return gs.Right[int](v)
		}

		return gs.Left[int, int](10 + v)
	}

	a, b := slices.PartitionMap(s, op)
	assert.Equal(t, slices.From(11, 13, 15, 17, 19), a)
	assert.Equal(t, slices.From(2, 4, 6, 8), b)
}

func TestGroupMap(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	m := slices.GroupMap(s, even, strconv.Itoa)
	assert.Equal(t, m[true], slices.From("2", "4", "6", "8"))
	assert.Equal(t, m[false], slices.From("1", "3", "5", "7", "9"))
}

func TestGroupBy(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	m := slices.GroupBy(s, even)

	assert.Equal(t, m[true], slices.From(2, 4, 6, 8))
	assert.Equal(t, m[false], slices.From(1, 3, 5, 7, 9))
}

func TestGroupMapReduce(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)

	m := slices.GroupMapReduce(
		s,
		even,
		func(v int) int { return v + 1 },
		func(a, b int) int { return a + b },
	)

	assert.Equal(t, 24, m[true])
	assert.Equal(t, 30, m[false])
}

func TestMaxBy(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	people := slices.Map(s, func(v int) person {
		return person{age: v}
	})

	p := slices.MaxBy(people, orderize)
	assert.Equal(t, 9, p.Get().age)
}

func TestMinBy(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	people := slices.Map(s, func(v int) person {
		return person{age: v}
	})

	p := slices.MinBy(people, orderize)
	assert.Equal(t, 1, p.Get().age)
}

func TestSoryBy(t *testing.T) {
	s := slices.From(1, 3, 5, 7, 9, 2, 4, 6, 8)
	people := slices.Map(s, func(v int) person {
		return person{age: v}
	})

	p := slices.SortBy(people, orderize)

	assert.Equal(t,
		slices.From(1, 2, 3, 4, 5, 6, 7, 8, 9),
		slices.Map(p, func(x person) int { return x.age }),
	)
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, slices.IsEmpty(slices.Empty[int]()))
	assert.False(t, slices.IsEmpty(slices.One(0)))
}
