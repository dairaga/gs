package try_test

import (
	"errors"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/try"
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

func TestFrom(t *testing.T) {
	assertTry(t, gs.Success(0), try.From(0, nil))
	assertTry(t, gs.Failure[int](gs.ErrLeft), try.From(0, gs.ErrLeft))
}

func TestFromWithBool(t *testing.T) {
	assertTry(t, gs.Success(0), try.FromWithBool(0, true))
	assertTry(t, gs.Failure[int](gs.ErrUnsatisfied), try.FromWithBool(0, false))
}

func TestFold(t *testing.T) {

}

func TestCollect(t *testing.T) {}

func TestFlatMap(t *testing.T) {}

func TestMap(t *testing.T) {}

func TestTryMap(t *testing.T) {}

func TestCanMap(t *testing.T) {}

func TestTransform(t *testing.T) {}
