package either_test

import (
	"strconv"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/either"
	"github.com/stretchr/testify/assert"
)

func assertEither[L, R any](t *testing.T, a, b gs.Either[L, R]) {
	t.Helper()

	assert.Equal(t, a.IsRight(), b.IsRight())

	if a.IsRight() {
		assert.Equal(t, a.Right(), b.Right())
	}

	if a.IsLeft() {
		assert.Equal(t, a.Left(), b.Left())
	}
}

func TestFold(t *testing.T) {
	fa := func(v int) int64 {
		return int64(v + 10)
	}

	fb := func(s string) int64 {
		a, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return -1
		}
		return a * 10
	}

	r := gs.Right[int]("1000")
	e := either.Fold(r, fa, fb)
	assert.Equal(t, int64(1000*10), e)

	r = gs.Right[int]("abc")
	e = either.Fold(r, fa, fb)
	assert.Equal(t, int64(-1), e)

	l := gs.Left[int, string](100)
	e = either.Fold(l, fa, fb)
	assert.Equal(t, int64(100+10), e)
}

func TestFlatMap(t *testing.T) {
	f := func(v string) gs.Either[int, int64] {
		a, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return gs.Left[int, int64](0)
		}
		return gs.Right[int](a)
	}

	assertEither(t,
		gs.Right[int](int64(1000)),
		either.FlatMap(gs.Right[int]("1000"), f))

	assertEither(t,
		gs.Left[int, int64](0),
		either.FlatMap(gs.Right[int]("abc"), f),
	)

	assertEither(t,
		gs.Left[int, int64](100),
		either.FlatMap(gs.Left[int, string](100), f),
	)
}

func TestMap(t *testing.T) {
	/*
	   Right(12).map(x => "flower") // Result: Right("flower")
	   Left(12).map(x => "flower")  // Result: Left(12)
	*/
	f := func(_ int) string {
		return "flower"
	}

	assertEither(t, gs.Right[int]("flower"), either.Map(gs.Right[int](12), f))
	assertEither(t, gs.Left[int, string](12), either.Map(gs.Left[int, int](12), f))
}
