// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package try_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
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
	fail := func(err error) string {
		return err.Error()
	}
	assert.Equal(t, "1", try.Fold(gs.Success(1), fail, strconv.Itoa))
	assert.Equal(t,
		gs.ErrEmpty.Error(),
		try.Fold(gs.Failure[int](gs.ErrEmpty), fail, strconv.Itoa))
}

func TestCollect(t *testing.T) {
	convErr := errors.New("conv error")

	p := func(v string) (int, error) {
		if len(v) <= 2 {
			a, err := strconv.Atoi(v)
			if err == nil {
				return a, nil
			}
			return a, convErr
		}
		return 0, gs.ErrUnsatisfied
	}

	assertTry(t, gs.Success(1), try.Collect(gs.Success("1"), p))
	assertTry(
		t,
		gs.Failure[int](gs.ErrUnsatisfied),
		try.Collect(gs.Success("100"), p),
	)

	assertTry(
		t,
		gs.Failure[int](convErr),
		try.Collect(gs.Success("ab"), p),
	)

}

func TestFlatMap(t *testing.T) {
	op := func(v int) gs.Try[string] {
		return gs.Success(strconv.Itoa(v))
	}

	assertTry(t, gs.Success("1"), try.FlatMap(gs.Success(1), op))
	assertTry(t,
		gs.Failure[string](gs.ErrEmpty),
		try.FlatMap(gs.Failure[int](gs.ErrEmpty), op))
}

func TestMap(t *testing.T) {
	assertTry(
		t,
		gs.Success("1"),
		try.Map(gs.Success(1), strconv.Itoa),
	)

	assertTry(
		t,
		gs.Failure[string](gs.ErrEmpty),
		try.Map(gs.Failure[int](gs.ErrEmpty), strconv.Itoa),
	)
}

func TestTryMap(t *testing.T) {
	assertTry(t, gs.Success(1), try.TryMap(gs.Success("1"), strconv.Atoi))
	assertTry(t,
		gs.Failure[int](gs.ErrEmpty),
		try.TryMap(gs.Failure[string](gs.ErrEmpty), strconv.Atoi),
	)
}

func TestCanMap(t *testing.T) {
	op := func(s string) (int, bool) {
		a, err := strconv.Atoi(s)
		return a, err == nil
	}

	assertTry(t, gs.Success(1), try.CanMap(gs.Success("1"), op))
	assertTry(t,
		gs.Failure[int](gs.ErrEmpty),
		try.CanMap(gs.Failure[string](gs.ErrEmpty), op),
	)
}

func TestTransform(t *testing.T) {
	errConv := errors.New("conv error")
	succ := funcs.TryRecover(
		func(s string) (int, error) {
			a, err := strconv.Atoi(s)
			if err != nil {
				return a, errConv
			}
			return a, nil
		},
		try.From[int])

	fail := gs.Failure[int]

	assertTry(t,
		gs.Success(1),
		try.Transform(gs.Success("1"), fail, succ))

	assertTry(t,
		gs.Failure[int](errConv),
		try.Transform(gs.Success("abc"), fail, succ))

	assertTry(t,
		gs.Failure[int](gs.ErrEmpty),
		try.Transform(gs.Failure[string](gs.ErrEmpty), fail, succ),
	)

}
