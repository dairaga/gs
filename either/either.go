package either

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

func Fold[L, R, T any](e gs.Either[L, R], left funcs.Func[L, T], right funcs.Func[R, T]) T {
	return funcs.BuildUnit(e.Fetch, right, funcs.UnitAndThen(e.Left, left))
}

func FlatMap[L, R, T any](e gs.Either[L, R], op funcs.Func[R, gs.Either[L, T]]) gs.Either[L, T] {
	return Fold(e, gs.Left[L, T], op)
}

func Map[L, R, T any](e gs.Either[L, R], op funcs.Func[R, T]) gs.Either[L, T] {
	return Fold(e, gs.Left[L, T], funcs.AndThen(op, gs.Right[L, T]))
}
