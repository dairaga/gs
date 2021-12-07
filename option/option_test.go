package option_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/option"
	"github.com/stretchr/testify/assert"
)

func assertOption[T any](t *testing.T, a, b gs.Option[T]) {
	t.Helper()
	assert.Equal(t, a.IsDefined(), b.IsDefined())
	aval, aerr := a.Fetch()
	bval, berr := b.Fetch()
	assert.Equal(t, aval, bval)
	assert.Equal(t, aerr, berr)
}

var trueF = func() bool { return true }
var falseF = func() bool { return false }

func TestFrom(t *testing.T) {
	assertOption(t, option.From(0, true), gs.Some(0))
	assertOption(t, option.From(100, false), gs.None[int]())
}

func TestFromWithErr(t *testing.T) {
	assertOption(t, option.FromWithErr(0, nil), gs.Some(0))
	assertOption(t, option.FromWithErr(100, errors.New("test")), gs.None[int]())
}

func TestWhen(t *testing.T) {
	assertOption(t, option.When(trueF, 0), gs.Some(0))
	assertOption(t, option.When(falseF, 0), gs.None[int]())
}

func TestUnless(t *testing.T) {
	assertOption(t, option.Unless(trueF, 0), gs.None[int]())
	assertOption(t, option.Unless(falseF, 0), gs.Some(0))
}

func TestFold(t *testing.T) {
	z := `zero`
	assert.Equal(t, "1", option.Fold(gs.Some(1), z, strconv.Itoa))
	assert.Equal(t, z, option.Fold(gs.None[int](), z, strconv.Itoa))
}

func TestCollect(t *testing.T) {
	p := func(v int) (s string, ok bool) {
		if ok = (v == 100); ok {
			s = strconv.Itoa(v)
		}
		return
	}

	assertOption(t, gs.Some("100"), option.Collect(gs.Some(100), p))
	assertOption(t, gs.None[string](), option.Collect(gs.None[int](), p))
	assertOption(t, gs.None[string](), option.Collect(gs.Some(1), p))
}

func TestFlatMap(t *testing.T) {
	op := funcs.AndThen(strconv.Itoa, gs.Some[string])

	assertOption(t, gs.Some("1"), option.FlatMap(gs.Some(1), op))
	assertOption(t, gs.None[string](), option.FlatMap(gs.None[int](), op))
}

func TestMap(t *testing.T) {
	assertOption(t, gs.Some("1"), option.Map(gs.Some(1), strconv.Itoa))
	assertOption(t, gs.None[string](), option.Map(gs.None[int](), strconv.Itoa))
}

func TestCankMap(t *testing.T) {
	op := func(s string) (int, bool) {
		v, err := strconv.Atoi(s)
		return v, err == nil
	}

	opt := option.CanMap(gs.Some("1"), op)
	assertOption(t, gs.Some(1), opt)

	opt = option.CanMap(gs.Some("abc"), op)
	assertOption(t, gs.None[int](), opt)

}

func TestTryMap(t *testing.T) {
	opt := option.TryMap(gs.Some("1"), strconv.Atoi)
	assertOption(t, gs.Some(1), opt)

	opt = option.TryMap(gs.Some("abc"), strconv.Atoi)
	assertOption(t, gs.None[int](), opt)
}

func TestLeft(t *testing.T) {
	e := option.Left(gs.Some(1), "1")
	assert.True(t, e.IsLeft())
	assert.Equal(t, 1, e.Left())

	e = option.Left(gs.None[int](), "1")
	assert.True(t, e.IsRight())
	assert.Equal(t, "1", e.Right())

}

func TestRight(t *testing.T) {
	e := option.Right(gs.Some(1), "1")
	assert.True(t, e.IsRight())
	assert.Equal(t, 1, e.Right())

	e = option.Right(gs.None[int](), "1")
	assert.True(t, e.IsLeft())
	assert.Equal(t, "1", e.Left())
}
