// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package funcs

type Fetcher[T any] func() (T, error)

func (f Fetcher[T]) Exists(p Predict[T]) bool {
	v, err := f()
	return err == nil && p(v)
}

func (f Fetcher[T]) Forall(p Predict[T]) bool {
	v, err := f()
	return err != nil || p(v)
}

func (f Fetcher[T]) Foreach(op func(T)) {
	if v, err := f(); err == nil {
		op(v)
	}
}

func (f Fetcher[T]) GetOrElse(z T) T {
	v, err := f()
	return Cond(err == nil, v, z)
}

func BuildWithErr[T, R any](v T, err error,
	fail Func[error, R], succ Func[T, R]) R {

	if err == nil {
		return succ(v)
	}
	return fail(err)
}

// -----------------------------------------------------------------------------

// TODO: refactor the method when go 1.19 releases.

func Build[T, R any](f Fetcher[T], fail Func[error, R], succ Func[T, R]) R {
	v, err := f()
	return BuildWithErr(v, err, fail, succ)
}

func BuildUnit[T, R any](f Fetcher[T], fail Unit[R], succ Func[T, R]) R {
	if v, err := f(); err == nil {
		return succ(v)
	}
	return fail()

}
