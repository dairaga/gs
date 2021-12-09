package maps_test

import (
	"strconv"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/maps"
	"github.com/dairaga/gs/slices"
	"github.com/stretchr/testify/assert"
)

var (
	testM = maps.From(
		maps.P(1, "1"),
		maps.P(3, "3"),
		maps.P(5, "5"),
		maps.P(7, "7"),
		maps.P(9, "9"),
		maps.P(2, "2"),
		maps.P(4, "4"),
		maps.P(6, "6"),
		maps.P(8, "8"),
	)

	oddM = maps.From(
		maps.P(1, "1"),
		maps.P(3, "3"),
		maps.P(5, "5"),
		maps.P(7, "7"),
		maps.P(9, "9"),
	)

	evenM = maps.From(
		maps.P(2, "2"),
		maps.P(4, "4"),
		maps.P(6, "6"),
		maps.P(8, "8"),
	)

	odd  = func(v int, _ string) bool { return (v & 0x01) == 1 }
	even = func(v int, _ string) bool { return (v & 0x01) == 0 }
	less = func(v int) func(int, string) bool {
		return func(x int, _ string) bool {
			return x < v
		}
	}
	large = func(v int) func(int, string) bool {
		return func(x int, _ string) bool {
			return x > v
		}
	}
)

func assertMap[K comparable, V any](t *testing.T, a, b maps.M[K, V]) {
	t.Helper()
	assert.Equal(t, len(a), len(b))

	for akey, aval := range a {
		assert.True(t, b.Contain(akey))
		assert.Equal(t, aval, b[akey])
	}
}

func TestMapKeys(t *testing.T) {
	m := maps.From(maps.P(1, "1"), maps.P(2, "2"), maps.P(3, "3"))

	keys := m.Keys()
	assert.Equal(t, slices.From(1, 2, 3), keys.Sort(funcs.Cmp[int]))
}

func TestMapValues(t *testing.T) {
	m := maps.From(maps.P(1, "1"), maps.P(2, "2"), maps.P(3, "3"))

	values := m.Values()
	assert.Equal(t, slices.From("1", "2", "3"), values.Sort(funcs.Cmp[string]))
}

func TestMapAdd(t *testing.T) {
	m := maps.From(maps.P(1, "1"), maps.P(2, "2"), maps.P(3, "3"))
	m.Add(maps.P(4, "4"), maps.P(1, "11"))
	assert.True(t, m.Contain(4))
	assert.Equal(t, "4", m[4])
	assert.Equal(t, "11", m[1])

}

func TestMapContain(t *testing.T) {
	m := maps.From(maps.P(1, "1"), maps.P(2, "2"), maps.P(3, "3"))

	assert.True(t, m.Contain(1))
	assert.False(t, m.Contain(4))
}

func TestMapCount(t *testing.T) {
	assert.Equal(t, 4, testM.Count(even))
	assert.Equal(t, 5, testM.Count(odd))
}

func TestMapFind(t *testing.T) {
	assert.True(t, testM.Find(even).IsDefined())
	assert.False(t, testM.Find(less(0)).IsDefined())
}

func TestMapExists(t *testing.T) {

	assert.True(t, testM.Exists(even))
	assert.False(t, testM.Exists(less(0)))
}

func TestMapFilter(t *testing.T) {
	assertMap(
		t,
		maps.From[int, string](),
		maps.From[int, string]().Filter(odd),
	)

	assertMap(t, oddM, testM.Filter(odd))
}

func TestMapFilterNot(t *testing.T) {
	assertMap(
		t,
		maps.From[int, string](),
		maps.From[int, string]().Filter(odd),
	)

	assertMap(t, evenM, testM.FilterNot(odd))
}
func TestMapForall(t *testing.T) {
	assert.True(t, maps.From[int, string]().Forall(large(0)))

	m := maps.From(maps.P(1, "1"), maps.P(2, "2"))
	assert.True(t, m.Forall(large(0)))

	m[-1] = "1"
	assert.False(t, m.Forall(large(0)))

}

func TestMapForeach(t *testing.T) {
	sum := 0

	op := func(k int, _ string) {
		sum += k
	}
	testM.Foreach(op)

	assert.Equal(t, 45, sum)
}

func TestMapPartition(t *testing.T) {
	a, b := testM.Partition(even)

	assertMap(t, oddM, a)
	assertMap(t, evenM, b)
}
func TestMapToSlice(t *testing.T) {
	s := slices.From(
		maps.P(1, "1"),
		maps.P(3, "3"),
		maps.P(5, "5"),
		maps.P(7, "7"),
		maps.P(9, "9"),
		maps.P(2, "2"),
		maps.P(4, "4"),
		maps.P(6, "6"),
		maps.P(8, "8"),
	)

	s = slices.SortBy(s, func(p maps.Pair[int, string]) int { return p.Key })

	assert.Equal(
		t,
		s,
		slices.SortBy(
			testM.Slice(),
			func(p maps.Pair[int, string]) int { return p.Key },
		),
	)
}

func TestFrom(t *testing.T) {
	assertMap(t, make(maps.M[int, int]), maps.From[int, int]())

	assertMap(t,
		map[int]int{1: 1, 2: 2, 3: 3},
		maps.From(maps.P(1, 1), maps.P(2, 2), maps.P(3, 3)),
	)
}

func TestFold(t *testing.T) {
	s := maps.Fold(
		testM,
		slices.Empty[string](),
		func(z slices.S[string], k int, v string) slices.S[string] {
			return append(z, v+strconv.Itoa(k))
		},
	)

	assert.Equal(t,
		slices.From("11", "22", "33", "44", "55", "66", "77", "88", "99"),
		s.Sort(funcs.Cmp[string]),
	)
}

func TestCollect(t *testing.T) {
	s := maps.Collect(testM, func(k int, v string) (string, bool) {
		return v + strconv.Itoa(k), even(k, v)
	})

	assert.Equal(t,
		slices.From("22", "44", "66", "88"),
		s.Sort(funcs.Cmp[string]),
	)
}

func TestCollectMap(t *testing.T) {
	m1 := maps.CollectMap(testM, func(k int, v string) (int, string, bool) {
		return k, v, even(k, v)
	})

	assert.Equal(t,
		evenM,
		m1,
	)
}

func TestFlatMapSlice(t *testing.T) {

	s := maps.FlatMapSlice(
		testM,
		func(k int, v string) slices.S[string] {
			return slices.One(v + v)
		},
	)

	assert.Equal(t,
		slices.From("11", "22", "33", "44", "55", "66", "77", "88", "99"),
		s.Sort(funcs.Cmp[string]),
	)
}

func TestFlatMap(t *testing.T) {
	m := maps.FlatMap(testM, func(key int, val string) maps.M[string, int] {
		return maps.From(maps.P(val, key))
	})

	testM.Foreach(func(k int, v string) {
		assert.Equal(t, m[v], k)
	})
}

func TestMapSlice(t *testing.T) {
	s := maps.MapSlice(testM, func(k int, v string) string {
		return v + strconv.Itoa(k)
	})

	assert.Equal(t,
		slices.From("11", "22", "33", "44", "55", "66", "77", "88", "99"),
		s.Sort(funcs.Cmp[string]),
	)
}

func TestMap(t *testing.T) {
	m := maps.Map(testM, func(key int, val string) (string, int) {
		return val, key
	})

	testM.Foreach(func(k int, v string) {
		assert.Equal(t, m[v], k)
	})
}

func TestGroupMap(t *testing.T) {
	m := maps.GroupMap(
		testM,
		odd,
		func(k int, v string) string { return v + strconv.Itoa(k) },
	)

	m.Foreach(func(_ bool, v slices.S[string]) {
		v = v.Sort(funcs.Cmp[string])
	})

	odds := slices.Map(oddM.Values(), func(v string) string { return v + v })
	odds = odds.Sort(funcs.Cmp[string])

	evens := slices.Map(evenM.Values(), func(v string) string { return v + v })
	evens = evens.Sort(funcs.Cmp[string])

	assertMap(t, maps.From(maps.P(true, odds), maps.P(false, evens)), m)
}

func TestGroupBy(t *testing.T) {
	m := maps.GroupBy(testM, odd)
	assertMap(t, m[true], oddM)
	assertMap(t, m[false], evenM)
}

func TestGroupMapReduce(t *testing.T) {
	sum := func(a, b int) int { return a + b }

	m := maps.GroupMapReduce(
		testM,
		odd,
		func(k int, val string) int { return k },
		sum,
	)

	oddSum := oddM.Keys().Reduce(sum).Get()
	evenSum := evenM.Keys().Reduce(sum).Get()

	assertMap(t, maps.From(maps.P(true, oddSum), maps.P(false, evenSum)), m)
}

func TestPartitionMap(t *testing.T) {
	a, b := maps.PartitionMap(
		testM,
		func(k int, v string) gs.Either[int, string] {
			if odd(k, v) {
				return gs.Right[int](v)
			} else {
				return gs.Left[int, string](k)
			}
		},
	)

	assert.Equal(
		t,
		slices.From(2, 4, 6, 8),
		a.Sort(funcs.Cmp[int]),
	)

	assert.Equal(
		t,
		slices.From("1", "3", "5", "7", "9"),
		b.Sort(funcs.Cmp[string]),
	)
}

func TestMaxBy(t *testing.T) {
	assert.False(
		t,
		maps.MaxBy(
			maps.From[int, int](),
			func(key int, _ int) int { return key },
		).IsDefined())

	assert.Equal(
		t,
		gs.Some(maps.P(9, "9")),
		maps.MaxBy(
			testM,
			func(key int, _ string) int { return key },
		))
}

func TestMinBy(t *testing.T) {
	assert.False(
		t,
		maps.MinBy(
			maps.From[int, int](),
			func(key int, _ int) int { return key },
		).IsDefined())

	assert.Equal(
		t,
		gs.Some(maps.P(1, "1")),
		maps.MinBy(
			testM,
			func(key int, _ string) int { return key },
		))
}
