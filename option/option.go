// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package option

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// TODO: refactor the method when go 1.19 releases.

func From[T any](v T, ok bool) gs.Option[T] {
	if ok {
		return gs.Some(v)
	}
	return gs.None[T]()
}

func FromWithErr[T any](v T, err error) gs.Option[T] {
	return From(v, err == nil)
}

func When[T any](p funcs.Condition, z T) gs.Option[T] {
	return From(z, p())
}

func Unless[T any](p funcs.Condition, z T) gs.Option[T] {
	return From(z, !p())
}

func Fold[T, R any](o gs.Option[T], z R, succ funcs.Func[T, R]) R {
	return funcs.BuildUnit(o.Fetch, funcs.Id(z), succ)
}

func Collect[T, R any](o gs.Option[T], p funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.CanTransform(p, From[R]))
}

func FlatMap[T, R any](o gs.Option[T], op funcs.Func[T, gs.Option[R]]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], op)
}

func Map[T, R any](o gs.Option[T], op funcs.Func[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.AndThen(op, gs.Some[R]))
}

func CanMap[T, R any](o gs.Option[T], op funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.CanTransform(op, From[R]))
}

func TryMap[T, R any](o gs.Option[T], op funcs.Try[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.TryRecover(op, FromWithErr[R]))
}

func Left[T, R any](o gs.Option[T], z R) gs.Either[T, R] {
	return funcs.BuildUnit(o.Fetch, funcs.UnitAndThen(funcs.Id(z), gs.Right[T, R]), gs.Left[T, R])
}

func Right[L, T any](o gs.Option[T], z L) gs.Either[L, T] {
	return funcs.BuildUnit(o.Fetch, funcs.UnitAndThen(funcs.Id(z), gs.Left[L, T]), gs.Right[L, T])
}
