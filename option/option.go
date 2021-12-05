package option

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

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

func Fold[T, R any](o gs.Option[T], succ funcs.Func[T, R], fail funcs.Unit[R]) R {
	return funcs.BuildUnit(o.Fetch, succ, fail)
}

func Collect[T, R any](o gs.Option[T], p funcs.Check[T, R]) gs.Option[R] {
	return Fold(o, funcs.CheckAndTransform(p, From[R]), gs.None[R])
}

func FlatMap[T, R any](o gs.Option[T], op funcs.Func[T, gs.Option[R]]) gs.Option[R] {
	return Fold(o, op, gs.None[R])
}

func Map[T, R any](o gs.Option[T], op funcs.Func[T, R]) gs.Option[R] {
	return Fold(o, funcs.AndThen(op, gs.Some[R]), gs.None[R])
}

func CheckMap[T, R any](o gs.Option[T], op funcs.Check[T, R]) gs.Option[R] {
	return Fold(o, funcs.CheckAndTransform(op, From[R]), gs.None[R])
}

func TryMap[T, R any](o gs.Option[T], op funcs.Try[T, R]) gs.Option[R] {
	return Fold(o, funcs.TryAndRecover(op, FromWithErr[R]), gs.None[R])
}

func Left[T, R any](o gs.Option[T], z R) gs.Either[T, R] {
	return Fold(o, gs.Left[T, R], funcs.UnitAndThen(funcs.Id(z), gs.Right[T, R]))
}

func Right[L, T any](o gs.Option[T], z L) gs.Either[L, T] {
	return Fold(o, gs.Right[L, T], funcs.UnitAndThen(funcs.Id(z), gs.Left[L, T]))
}
