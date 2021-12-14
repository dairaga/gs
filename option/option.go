// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package option

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// From returns a Some with given v if given ok is true, or returns a None.
func From[T any](v T, ok bool) gs.Option[T] {
	if ok {
		return gs.Some(v)
	}
	return gs.None[T]()
}

// FromWithErr returns a Some with given v if given err is nil, or returns a None.
func FromWithErr[T any](v T, err error) gs.Option[T] {
	return From(v, err == nil)
}

// When returns a Some with given v if result of given function p is true, or returns a None.
func When[T any](p funcs.Condition, z T) gs.Option[T] {
	return From(z, p())
}

// Unless returns a Some with given v if result of given function p is false, or returns a None.
func Unless[T any](p funcs.Condition, z T) gs.Option[T] {
	return From(z, !p())
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// Fold returns result from applying given function succ if o is defined, or returns given default value z.
func Fold[T, R any](o gs.Option[T], z R, succ funcs.Func[T, R]) R {
	return funcs.BuildUnit(o.Fetch, funcs.Id(z), succ)
}

// Collect returns a Some with result from applying given function p if o is defined and value of o satifies p, or returns a None.
func Collect[T, R any](o gs.Option[T], p funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.CanTransform(p, From[R]))
}

// FlatMap returns result from applying given function op if o is defined, or returns a None.
func FlatMap[T, R any](o gs.Option[T], op funcs.Func[T, gs.Option[R]]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], op)
}

// Map returns a Some with result from applying given function op if o is defined, or returns a None.
func Map[T, R any](o gs.Option[T], op funcs.Func[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.AndThen(op, gs.Some[R]))
}

// CanMap returns a Some with result from applying given function op if o is defined and satisfies op, or returns a None.
func CanMap[T, R any](o gs.Option[T], op funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.CanTransform(op, From[R]))
}

// TryMap returns a Some with result from applying given function op if o is defined and converts to R successfully, or returns a None.
func TryMap[T, R any](o gs.Option[T], op funcs.Try[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, gs.None[R], funcs.TryRecover(op, FromWithErr[R]))
}

// Left returns a Left with value from o if o is defined, or returns a Right with given z.
func Left[T, R any](o gs.Option[T], z R) gs.Either[T, R] {
	return funcs.BuildUnit(o.Fetch, funcs.UnitAndThen(funcs.Id(z), gs.Right[T, R]), gs.Left[T, R])
}

// Right returns a Right with value from o if o is defined, or returns Left with given z.
func Right[L, T any](o gs.Option[T], z L) gs.Either[L, T] {
	return funcs.BuildUnit(o.Fetch, funcs.UnitAndThen(funcs.Id(z), gs.Left[L, T]), gs.Right[L, T])
}
