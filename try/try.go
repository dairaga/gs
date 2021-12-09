// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package try

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// TODO: refactor the method when go 1.19 releases.

func From[T any](v T, err error) gs.Try[T] {
	return funcs.BuildWithErr(v, err, gs.Failure[T], gs.Success[T])
}

func FromWithBool[T any](v T, ok bool) gs.Try[T] {
	return From(v, funcs.Cond(ok, nil, gs.ErrUnsatisfied))
}

func Fold[T, R any](t gs.Try[T],
	fail funcs.Func[error, R], succ funcs.Func[T, R]) R {

	return funcs.Build(t.Fetch, fail, succ)
}

func Collect[T, R any](t gs.Try[T], p funcs.Try[T, R]) gs.Try[R] {
	return funcs.Build(t.Fetch, gs.Failure[R], funcs.TryRecover(p, From[R]))
}

func FlatMap[T, R any](t gs.Try[T], op funcs.Func[T, gs.Try[R]]) gs.Try[R] {
	return funcs.Build(t.Fetch, gs.Failure[R], op)
}

func Map[T, R any](t gs.Try[T], op funcs.Func[T, R]) gs.Try[R] {
	return funcs.Build(t.Fetch, gs.Failure[R], funcs.AndThen(op, gs.Success[R]))
}

func TryMap[T, R any](t gs.Try[T], op funcs.Try[T, R]) gs.Try[R] {
	return funcs.Build(t.Fetch, gs.Failure[R], funcs.TryRecover(op, From[R]))
}

func CanMap[T, R any](t gs.Try[T], op funcs.Can[T, R]) gs.Try[R] {
	return funcs.Build(
		t.Fetch,
		gs.Failure[R],
		funcs.CanTransform(op, FromWithBool[R]),
	)
}

func Transform[T, R any](t gs.Try[T],
	fail funcs.Func[error, gs.Try[R]],
	succ funcs.Func[T, gs.Try[R]]) gs.Try[R] {

	return funcs.Build(t.Fetch, fail, succ)
}
