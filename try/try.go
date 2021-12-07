package try

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// TODO: refactor the method when go 1.19 releases.

func From[T any](v T, err error) gs.Try[T] {
	return funcs.BuildFrom(v, err, gs.Success[T], gs.Failure[T])
}

func FromWithBool[T any](v T, ok bool) gs.Try[T] {
	return From(v, funcs.Cond(ok, nil, gs.ErrUnsatisfied))
}

func Fold[T, R any](t gs.Try[T], succ funcs.Func[T, R], fail funcs.Func[error, R]) R {
	return funcs.Build(t.Fetch, succ, fail)
}

func Collect[T, R any](t gs.Try[T], p funcs.Try[T, R]) gs.Try[R] {
	return Fold(t, funcs.TryRecover(p, From[R]), gs.Failure[R])
}

func FlatMap[T, R any](t gs.Try[T], op funcs.Func[T, gs.Try[R]]) gs.Try[R] {
	return Fold(t, op, gs.Failure[R])
}

func Map[T, R any](t gs.Try[T], op funcs.Func[T, R]) gs.Try[R] {
	return Fold(t, funcs.AndThen(op, gs.Success[R]), gs.Failure[R])
}

func TryMap[T, R any](t gs.Try[T], op funcs.Try[T, R]) gs.Try[R] {
	return Fold(t, funcs.TryRecover(op, From[R]), gs.Failure[R])
}

func CanMap[T, R any](t gs.Try[T], op funcs.Can[T, R]) gs.Try[R] {
	return Fold(t, funcs.CanTransform(op, FromWithBool[R]), gs.Failure[R])
}

func Transform[T, R any](t gs.Try[T], succ funcs.Func[T, gs.Try[R]], fail funcs.Func[error, gs.Try[R]]) gs.Try[R] {
	return Fold(t, succ, fail)
}
