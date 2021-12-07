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

type Nothing struct{}

func (n Nothing) String() string {
	return "Nothing"
}

var nothing = struct{}{}

func N() Nothing { return nothing }

type Numeric interface {
	constraints.Integer | constraints.Float
}

type Tuple2[V1, V2 any] struct {
	_  struct{}
	V1 V1
	V2 V2
}

func T2[V1, V2 any](v1 V1, v2 V2) Tuple2[V1, V2] {
	return Tuple2[V1, V2]{
		V1: v1,
		V2: v2,
	}
}
