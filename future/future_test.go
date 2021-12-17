// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package future_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
	"github.com/dairaga/gs/future"
	"github.com/stretchr/testify/assert"
)

func assertResult[T any](t *testing.T, a, b gs.Try[T]) {
	t.Helper()

	assert.Equal(t, a.IsSuccess(), b.IsSuccess())

	if a.IsSuccess() {
		assert.Equal(t, a.Success(), b.Success())
	}

	if a.IsFailure() {
		assert.True(t, errors.Is(a.Failed(), b.Failed()))
	}
}

func TestRun(t *testing.T) {
	run := func() int {
		return 0
	}

	f := future.Run(context.Background(), run)
	f.Wait()
	t.Log(f)

	result, completed := f.Get()
	assert.True(t, f.Completed())
	assert.True(t, completed)

	assert.True(t, result.IsSuccess())
	assert.Equal(t, 0, result.Get())

	run = func() int {
		time.Sleep(3 * time.Second)
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	f = future.Run(ctx, run)
	<-ctx.Done()
	t.Log(f)
	result, completed = f.Get()
	assert.False(t, completed)
	assertResult(t, future.Failure[int](), result)

	run = func() int {
		panic(gs.ErrEmpty)
	}

	f = future.Run(context.Background(), run)
	f.Wait()
	result, completed = f.Get()
	assert.True(t, f.Completed())
	assert.True(t, completed)

	assert.True(t, result.IsFailure())
	assert.True(t, errors.Is(gs.ErrEmpty, result.Failed()))

}

func TestTry(t *testing.T) {
	try := func() (int, error) {
		return 0, nil
	}

	f := future.Try(context.Background(), try)
	f.Wait()

	result, completed := f.Get()
	assert.True(t, completed)
	assert.True(t, f.Completed())
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 0, result.Success())

	try = func() (int, error) {
		return 0, gs.ErrEmpty
	}

	f = future.Try(context.Background(), try)
	f.Wait()

	result, completed = f.Get()
	assert.True(t, completed)
	assert.True(t, f.Completed())
	assert.True(t, result.IsFailure())
	assert.True(t, errors.Is(gs.ErrEmpty, result.Failed()))

}

func TestCallback(t *testing.T) {
	ch := make(chan struct{}, 3)
	defer close(ch)

	try := func() (int, error) {
		return 1, nil
	}

	check := 0
	f := future.Try(context.Background(), try)
	f.OnCompleted(func(x gs.Try[int]) {
		check++
		assert.True(t, x.IsSuccess())
		assert.Equal(t, 1, x.Success())
		ch <- struct{}{}
	})
	f.OnSuccess(func(x int) {
		check++
		assert.Equal(t, 1, x)
		ch <- struct{}{}
	})
	f.OnError(func(err error) {
		check++
		ch <- struct{}{}
	})
	f.Wait()
	<-ch
	<-ch
	assert.Equal(t, 2, check)

	try = func() (int, error) {
		return 0, gs.ErrEmpty
	}

	check = 0
	f = future.Try(context.Background(), try)
	f.OnCompleted(func(x gs.Try[int]) {
		check++
		assert.True(t, x.IsFailure())
		assert.True(t, errors.Is(gs.ErrEmpty, x.Failed()))
		ch <- struct{}{}
	})
	f.OnSuccess(func(x int) {
		check++
		ch <- struct{}{}
	})

	f.OnError(func(err error) {
		check++
		assert.True(t, errors.Is(gs.ErrEmpty, err))
		ch <- struct{}{}
	})
	f.Wait()
	<-ch
	<-ch
	assert.Equal(t, 2, check)

}

func TestResult(t *testing.T) {
	run := func() int {
		return 1
	}

	f := future.Run(context.Background(), run)
	result := f.Result(context.Background(), 5*time.Second)
	assert.True(t, result.IsSuccess())
	assert.Equal(t, 1, result.Success())

	run = func() int {
		time.Sleep(5 * time.Second)
		return 1
	}

	f = future.Run(context.Background(), run)
	result = f.Result(context.Background(), time.Second)
	assert.True(t, result.IsFailure())
	assert.True(t, errors.Is(context.DeadlineExceeded, result.Failed()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	f = future.Run(ctx, run)
	result = f.Result(context.Background(), 10*time.Second)
	cancel()
	assert.True(t, result.IsFailure())
	assert.True(t, errors.Is(context.DeadlineExceeded, result.Failed()))

}

func TestFilter(t *testing.T) {
	/*
		val f = Future { 5 }
		val g = f filter { _ % 2 == 1 }
		val h = f filter { _ % 2 == 0 }
		g foreach println // Eventually prints 5
		Await.result(h, Duration.Zero) // throw a NoSuchElementException
	*/

	p1 := func(a int) bool {
		return (a % 2) == 1
	}

	p2 := func(a int) bool {
		return (a % 2) == 0
	}

	f := future.Run(context.Background(), funcs.Id(5))
	g := f.Filter(context.Background(), p1)
	h := f.Filter(context.Background(), p2)

	result := g.Result(context.Background(), 5*time.Second)
	assert.Equal(t, 5, result.Get())

	result = h.Result(context.Background(), 5*time.Second)
	assert.True(t, errors.Is(gs.ErrUnsatisfied, result.Failed()))

}

func TestFilterNot(t *testing.T) {
	p1 := func(a int) bool {
		return (a % 2) == 1
	}

	p2 := func(a int) bool {
		return (a % 2) == 0
	}

	f := future.Run(context.Background(), funcs.Id(5))
	g := f.FilterNot(context.Background(), p1)
	h := f.FilterNot(context.Background(), p2)

	result := g.Result(context.Background(), 5*time.Second)
	assert.True(t, errors.Is(gs.ErrUnsatisfied, result.Failed()))

	result = h.Result(context.Background(), 5*time.Second)
	assert.Equal(t, 5, result.Get())
}

func TestFlatAndMap(t *testing.T) {
	f := future.Run(context.Background(), funcs.Id(5))
	g := future.Run(context.Background(), funcs.Id(3))

	h := future.FlatMap(
		context.Background(),
		f,
		func(a int) gs.Future[int] {
			return future.Map(context.Background(), g, func(b int) int {
				return a * b
			})
		},
	)

	assert.Equal(t, 5*3, h.Wait().Get())

	h = future.FlatMap(
		context.Background(),
		f,
		func(a int) gs.Future[int] {
			return future.TryMap(context.Background(), g, func(b int) (int, error) {
				return 0, gs.ErrEmpty
			})
		},
	)

	assert.True(t, errors.Is(gs.ErrEmpty, h.Wait().Failed()))

	h = future.FlatMap(
		context.Background(),
		f,
		func(a int) gs.Future[int] {
			return future.PartialMap(context.Background(), g, func(b int) (int, bool) {
				return 0, false
			})
		},
	)

	assert.True(t, errors.Is(gs.ErrUnsatisfied, h.Wait().Failed()))
}

func TestZip(t *testing.T) {
	f := future.Run(context.Background(), func() int {
		time.Sleep(2 * time.Second)
		return 1
	})

	g := future.Run(context.Background(), func() string {
		time.Sleep(5 * time.Second)
		return "hello"
	})

	h := future.Zip(context.Background(), f, g)
	result := h.Wait()

	assert.True(t, result.IsSuccess())
	assert.Equal(t, result.Success().V1.Success(), 1)
	assert.Equal(t, result.Success().V2.Success(), "hello")

	f = future.Run(context.Background(), func() int {
		time.Sleep(5 * time.Second)
		return 1
	})

	g = future.Run(context.Background(), func() string {
		time.Sleep(2 * time.Second)
		return "hello"
	})

	assert.True(t, result.IsSuccess())
	assert.Equal(t, result.Success().V1.Success(), 1)
	assert.Equal(t, result.Success().V2.Success(), "hello")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	h = future.Zip(ctx, f, g)
	result = h.Wait()
	cancel()
	assert.True(t, result.IsFailure())
	assert.True(t, errors.Is(gs.ErrEmpty, result.Failed()))
}

func TestZipWith(t *testing.T) {
	f := future.Run(context.Background(), func() int {
		time.Sleep(2 * time.Second)
		return 1
	})

	g := future.Run(context.Background(), func() string {
		time.Sleep(5 * time.Second)
		return "hello"
	})

	op := func(a gs.Try[int], b gs.Try[string]) gs.Try[string] {
		if a.IsFailure() {
			return gs.Failure[string](a.Failed())
		}

		if b.IsFailure() {
			return gs.Failure[string](b.Failed())
		}

		return gs.Success("ok")
	}

	h := future.ZipWith(context.Background(), f, g, op)
	result := h.Wait()

	assertResult(t, gs.Success("ok"), result)
}
