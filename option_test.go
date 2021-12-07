package gs_test

import (
	"errors"
	"testing"

	"github.com/dairaga/gs"
	"github.com/stretchr/testify/assert"
)

func assertOption[T any](t *testing.T, a, b gs.Option[T]) {
	t.Helper()
	assert.Equal(t, a.IsDefined(), b.IsDefined())
	aval, aerr := a.Fetch()
	bval, berr := b.Fetch()
	assert.Equal(t, aval, bval)
	assert.True(t, errors.Is(aerr, berr))
}

func TestSome(t *testing.T) {
	v := 0
	opt := gs.Some(v)
	assert.True(t, opt.IsDefined())
	assert.Equal(t, v, opt.Get())

	x, err := opt.Fetch()
	assert.Equal(t, v, x)
	assert.Nil(t, err)

	y, ok := opt.Check()
	assert.Equal(t, v, y)
	assert.True(t, ok)

	try := opt.Try()
	assert.True(t, try.IsSuccess())
	assert.Equal(t, v, try.Get())
	assert.Equal(t, v, try.Success())

	e := opt.Either()
	assert.True(t, e.IsRight())
	assert.Equal(t, v, e.Right())
}

func TestNone(t *testing.T) {
	opt := gs.None[int]()
	assert.True(t, opt.IsEmpty())
	assert.Panics(t, func() { opt.Get() })

	x, err := opt.Fetch()
	assert.Equal(t, 0, x) // zero value
	assert.True(t, errors.Is(gs.ErrEmpty, err))

	y, ok := opt.Check()
	assert.Equal(t, 0, y) // zero value
	assert.False(t, ok)

	try := opt.Try()
	assert.True(t, try.IsFailure())
	assert.True(t, errors.Is(gs.ErrEmpty, try.Failed()))

	e := opt.Either()
	assert.True(t, e.IsLeft())
	assert.True(t, errors.Is(gs.ErrEmpty, e.Left()))
}

func TestOptionExists(t *testing.T) {
	p := func(v int) bool {
		return v > 0
	}

	assert.True(t, gs.Some(1).Exists(p))
	assert.False(t, gs.Some(-1).Exists(p))
	assert.False(t, gs.None[int]().Exists(p))
}

func TestOptionForall(t *testing.T) {
	p1 := func(v int) bool {
		return v == 100
	}

	p2 := func(v int) bool {
		return v < 0
	}

	assert.True(t, gs.Some(100).Forall(p1))
	assert.False(t, gs.Some(100).Forall(p2))

	assert.True(t, gs.None[int]().Forall(p1))
	assert.True(t, gs.None[int]().Forall(p2))

}

func TestOptionForeach(t *testing.T) {
	sum := 123
	op := func(v int) {
		sum += v
	}

	gs.Some(100).Foreach(op)
	assert.Equal(t, 123+100, sum)

	sum = 123
	gs.None[int]().Foreach(op)
	assert.Equal(t, 123, sum)

}

func TestOptionFilter(t *testing.T) {
	p1 := func(v int) bool {
		return v == 100
	}

	p2 := func(v int) bool {
		return v < 0
	}

	assertOption(t, gs.Some(100).Filter(p1), gs.Some(100))
	assertOption(t, gs.Some(100).Filter(p2), gs.None[int]())

	assertOption(t, gs.None[int]().Filter(p1), gs.None[int]())
	assertOption(t, gs.None[int]().Filter(p2), gs.None[int]())
}

func TestOptionFilterNot(t *testing.T) {
	p1 := func(v int) bool {
		return v == 100
	}

	p2 := func(v int) bool {
		return v < 0
	}

	assertOption(t, gs.Some(100).FilterNot(p1), gs.None[int]())
	assertOption(t, gs.Some(100).FilterNot(p2), gs.Some(100))

	assertOption(t, gs.None[int]().FilterNot(p1), gs.None[int]())
	assertOption(t, gs.None[int]().FilterNot(p2), gs.None[int]())
}

func TestOptionGetOrElse(t *testing.T) {
	assert.Equal(t, 1, gs.Some(1).GetOrElse(-1))
	assert.Equal(t, -1, gs.None[int]().GetOrElse(-1))
}

func TestOptionOrElse(t *testing.T) {
	z := gs.Some(1)

	assertOption(t, gs.Some(100).OrElse(z), gs.Some(100))
	assertOption(t, gs.None[int]().OrElse(z), z)
}
