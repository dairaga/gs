// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gs

import (
	"constraints"
	"errors"
)

var (

	// ErrEmpty represents object is empty or undefined.
	ErrEmpty = errors.New("empty")

	// ErrLeft represents Either is Left.
	ErrLeft = errors.New("Left")

	// ErrUnsupported represents behavior is not supported.
	ErrUnsupported = errors.New("unsupported")

	// ErrUnsatisfied represents prediction is failed.
	ErrUnsatisfied = errors.New("unsatisfied")
)

// Nothing represents Nothing in scala.
type Nothing struct{}

func (n Nothing) String() string {
	return "Nothing"
}

var nothing = struct{}{}

// N returns Nothing.
func N() Nothing { return nothing }

// Numeric includes Integer and Float.
type Numeric interface {
	constraints.Integer | constraints.Float
}

// Tuple2 is tuple with size 2.
type Tuple2[V1, V2 any] struct {
	_  struct{}
	V1 V1
	V2 V2
}

// T2 returns Tuple2.
func T2[V1, V2 any](v1 V1, v2 V2) Tuple2[V1, V2] {
	return Tuple2[V1, V2]{
		V1: v1,
		V2: v2,
	}
}
