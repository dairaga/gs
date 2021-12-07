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

func Fold[T, R any](o gs.Option[T], z R, op funcs.Func[T, R]) R {
	// FIXME: reference scala
	return funcs.BuildOrElse(o.Fetch, z, op)
}

func Collect[T, R any](o gs.Option[T], p funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, funcs.CanTransform(p, From[R]), gs.None[R])
}

func FlatMap[T, R any](o gs.Option[T], op funcs.Func[T, gs.Option[R]]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, op, gs.None[R])
}

func Map[T, R any](o gs.Option[T], op funcs.Func[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, funcs.AndThen(op, gs.Some[R]), gs.None[R])
}

func CanMap[T, R any](o gs.Option[T], op funcs.Can[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, funcs.CanTransform(op, From[R]), gs.None[R])
}

func TryMap[T, R any](o gs.Option[T], op funcs.Try[T, R]) gs.Option[R] {
	return funcs.BuildUnit(o.Fetch, funcs.TryRecover(op, FromWithErr[R]), gs.None[R])
}

func Left[T, R any](o gs.Option[T], z R) gs.Either[T, R] {
	return funcs.BuildUnit(o.Fetch, gs.Left[T, R], funcs.UnitAndThen(funcs.Id(z), gs.Right[T, R]))
}

func Right[L, T any](o gs.Option[T], z L) gs.Either[L, T] {
	return funcs.BuildUnit(o.Fetch, gs.Right[L, T], funcs.UnitAndThen(funcs.Id(z), gs.Left[L, T]))
}
