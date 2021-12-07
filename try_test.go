package gs_test

import (
	"errors"
	"testing"

	"github.com/dairaga/gs"
	"github.com/stretchr/testify/assert"
)

func assertTry[T any](t *testing.T, a, b gs.Try[T]) {
	t.Helper()
	assert.Equal(t, a.IsSuccess(), b.IsSuccess())

	aval, aerr := a.Fetch()
	bval, berr := b.Fetch()

	assert.Equal(t, aval, bval)
	assert.True(t, errors.Is(aerr, berr))
}

func TestSuccess(t *testing.T) {
	try := gs.Success(0)
	t.Log(try)
	assert.True(t, try.IsSuccess())
	assert.Equal(t, 0, try.Get())
	assert.Equal(t, 0, try.Success())
	assert.True(t, errors.Is(gs.ErrUnsupported, try.Failed()))

	v, err := try.Fetch()
	assert.Equal(t, 0, v)
	assert.Nil(t, err)

	e := try.Either()
	assert.True(t, e.IsRight())
	assert.Equal(t, try.Success(), e.Right())

	opt := try.Option()
	assert.True(t, opt.IsDefined())
	assert.Equal(t, opt.Get(), try.Success())
}

func TestFailure(t *testing.T) {
	try := gs.Failure[int](gs.ErrUnsatisfied)
	t.Log(try)
	assert.True(t, try.IsFailure())
	assert.True(t, errors.Is(gs.ErrUnsatisfied, try.Failed()))
	assert.Panics(t, func() { try.Get() })
	assert.Panics(t, func() { try.Success() })

	v, err := try.Fetch()
	assert.Equal(t, 0, v)
	assert.True(t, errors.Is(err, gs.ErrUnsatisfied))

	e := try.Either()
	assert.True(t, e.IsLeft())
	assert.True(t, errors.Is(gs.ErrUnsatisfied, e.Left()))

	opt := try.Option()
	assert.True(t, opt.IsEmpty())
}

func TestTryExists(t *testing.T) {
	p := func(v int) bool {
		return v > 0
	}

	assert.True(t, gs.Success(1).Exists(p))
	assert.False(t, gs.Success(-1).Exists(p))
	assert.False(t, gs.Failure[int](gs.ErrEmpty).Exists(p))

}
func TestTryForall(t *testing.T) {
	p := func(v int) bool {
		return v > 0
	}

	assert.True(t, gs.Success(1).Forall(p))
	assert.False(t, gs.Success(-1).Forall(p))
	assert.True(t, gs.Failure[int](gs.ErrEmpty).Forall(p))
}

func TestTryForeach(t *testing.T) {
	sum := 0
	op := func(v int) {
		sum += v
	}

	gs.Success(1).Foreach(op)
	assert.Equal(t, 1, sum)

	sum = 0
	gs.Failure[int](gs.ErrEmpty).Foreach(op)
	assert.Equal(t, 0, sum)
}

func TestTryFilter(t *testing.T) {
	p := func(v int) bool {
		return v > 0
	}

	assertTry(t, gs.Success(1), gs.Success(1).Filter(p))
	assertTry(t, gs.Failure[int](gs.ErrUnsatisfied), gs.Success(-1).Filter(p))

	err := errors.New("test")
	assertTry(t, gs.Failure[int](err), gs.Failure[int](err).Filter(p))
}

func TestTryFilterNot(t *testing.T) {
	p := func(v int) bool {
		return v > 0
	}

	assertTry(t, gs.Failure[int](gs.ErrUnsatisfied), gs.Success(1).FilterNot(p))
	assertTry(t, gs.Success(-1), gs.Success(-1).FilterNot(p))

	err := errors.New("test")
	assertTry(t, gs.Failure[int](err), gs.Failure[int](err).FilterNot(p))
}
func TestTryGetOrElse(t *testing.T) {
	assert.Equal(t, 1, gs.Success(1).GetOrElse(-1))
	assert.Equal(t, -1, gs.Failure[int](gs.ErrUnsatisfied).GetOrElse(-1))
}
func TestTryOrElse(t *testing.T) {
	z := gs.Success(-1)

	assertTry(t, gs.Success(1), gs.Success(1).OrElse(z))
	assertTry(t, z, gs.Failure[int](gs.ErrEmpty).OrElse(z))
}

func TestTryRecover(t *testing.T) {
	r := func(err error) string {
		return err.Error()
	}

	assertTry(t, gs.Success("1"), gs.Success("1").Recover(r))
	assertTry(t, gs.Success(gs.ErrEmpty.Error()), gs.Failure[string](gs.ErrEmpty).Recover(r))
}
func TestTryRecoverWith(t *testing.T) {
	r := func(err error) gs.Try[string] {
		if errors.Is(gs.ErrLeft, err) {
			return gs.Success(err.Error())
		}
		return gs.Failure[string](err)
	}

	assertTry(t, gs.Success("1"), gs.Success("1").RecoverWith(r))
	assertTry(t, gs.Success(gs.ErrLeft.Error()), gs.Failure[string](gs.ErrLeft).RecoverWith(r))
	assertTry(t, gs.Failure[string](gs.ErrEmpty), gs.Failure[string](gs.ErrEmpty).RecoverWith(r))
}
