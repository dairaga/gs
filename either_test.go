// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs_test

import (
	"errors"
	"testing"

	"github.com/dairaga/gs"
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

func TestRight(t *testing.T) {
	e := gs.Right[string](1)
	t.Log(e)

	assert.True(t, e.IsRight())
	assert.Equal(t, 1, e.Right())
	assert.Equal(t, 1, e.Get())
	assert.Panics(t, func() { e.Left() })

	v, err := e.Fetch()
	assert.Equal(t, 1, v)
	assert.Nil(t, err)

	try := e.Try()
	assert.True(t, try.IsSuccess())
	assert.Equal(t, 1, try.Get())

	opt := e.Option()
	assert.True(t, opt.IsDefined())
	assert.Equal(t, 1, opt.Get())
}

func TestLeft(t *testing.T) {
	e := gs.Left[string, int]("1")
	t.Log(e)

	assert.True(t, e.IsLeft())
	assert.Equal(t, "1", e.Left())
	assert.Panics(t, func() { e.Get() })
	assert.Panics(t, func() { e.Right() })

	v, err := e.Fetch()
	assert.Equal(t, 0, v)
	assert.True(t, errors.Is(gs.ErrLeft, err))

	try := e.Try()
	assert.True(t, try.IsFailure())
	assert.True(t, errors.Is(gs.ErrLeft, try.Failed()))

	opt := e.Option()
	assert.True(t, opt.IsEmpty())

	e1 := gs.Left[error, int](gs.ErrEmpty)

	v1, err1 := e1.Fetch()
	assert.Equal(t, 0, v1)
	assert.True(t, errors.Is(gs.ErrEmpty, err1))

	try = e1.Try()
	assert.True(t, try.IsFailure())
	assert.True(t, errors.Is(gs.ErrEmpty, try.Failed()))
}

func TestEitherExists(t *testing.T) {
	/*
	   Right(12).exists(_ > 10)   // true
	   Right(7).exists(_ > 10)    // false
	   Left(12).exists(_ => true) // false
	*/

	p1 := func(v int) bool {
		return v > 10
	}

	p2 := func(_ int) bool {
		return true
	}

	assert.True(t, gs.Right[gs.Nothing](12).Exists(p1))
	assert.False(t, gs.Right[gs.Nothing](7).Exists(p1))
	assert.False(t, gs.Left[gs.Nothing, int](gs.N()).Exists(p2))
}

func TestEitherForall(t *testing.T) {
	p1 := func(v int) bool {
		return v > 10
	}

	p2 := func(_ int) bool {
		return false
	}

	assert.True(t, gs.Right[gs.Nothing](12).Forall(p1))
	assert.False(t, gs.Right[gs.Nothing](7).Forall(p1))
	assert.True(t, gs.Left[gs.Nothing, int](gs.N()).Forall(p2))

}

func TestEitherForeach(t *testing.T) {
	sum := 0
	op := func(v int) {
		sum += v
	}

	gs.Right[gs.Nothing](1).Foreach(op)
	assert.Equal(t, 1, sum)

	sum = 0
	gs.Left[gs.Nothing, int](gs.N()).Foreach(op)
	assert.Equal(t, 0, sum)
}

func TestEitherFilterOrElse(t *testing.T) {
	/*
	   Right(12).filterOrElse(_ > 10, -1)   // Right(12)
	   Right(7).filterOrElse(_ > 10, -1)    // Left(-1)
	   Left(7).filterOrElse(_ => false, -1) // Left(7)
	*/

	p1 := func(v int) bool {
		return v > 10
	}

	p2 := func(_ int) bool {
		return false
	}

	assertEither(t, gs.Right[int](12), gs.Right[int](12).FilterOrElse(p1, -1))
	assertEither(t, gs.Left[int, int](-1), gs.Right[int](7).FilterOrElse(p1, -1))
	assertEither(t, gs.Left[int, int](7), gs.Left[int, int](7).FilterOrElse(p2, -1))
}

func TestEitherGetOrElse(t *testing.T) {
	/*
	   Right(12).getOrElse(17) // 12
	   Left(12).getOrElse(17)  // 17
	*/

	assert.Equal(t, 12, gs.Right[int](12).GetOrElse(17))
	assert.Equal(t, 17, gs.Left[int, int](12).GetOrElse(17))
}

func TestEitherOrElse(t *testing.T) {
	z := gs.Right[int](1)

	assertEither(t, gs.Right[int](10), gs.Right[int](10).OrElse(z))
	assertEither(t, gs.Right[int](1), gs.Left[int, int](10).OrElse(z))
}

func TestEitherSwap(t *testing.T) {
	assertEither(t, gs.Right[int]("left"), gs.Left[string, int]("left").Swap())
}
