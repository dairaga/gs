// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package either

import (
	"github.com/dairaga/gs"
	"github.com/dairaga/gs/funcs"
)

// TODO: refactor following functions to methods when go 1.19 releases.

// Fold applies given function right to value from Right,
// or function left to value From Left.
func Fold[L, R, T any](e gs.Either[L, R],
	left funcs.Func[L, T], right funcs.Func[R, T]) T {

	return funcs.BuildUnit(e.Fetch, funcs.UnitAndThen(e.Left, left), right)
}

// FlapMap returns result applying given function op on value from Right, or returns Left.
func FlatMap[L, R, T any](e gs.Either[L, R],
	op funcs.Func[R, gs.Either[L, T]]) gs.Either[L, T] {

	return Fold(e, gs.Left[L, T], op)
}

// Map returns new Right with result of applying function op if e is a Right, or return Left.
func Map[L, R, T any](e gs.Either[L, R],
	op funcs.Func[R, T]) gs.Either[L, T] {

	return Fold(e, gs.Left[L, T], funcs.AndThen(op, gs.Right[L, T]))
}
