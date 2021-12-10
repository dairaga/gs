// Copyright © 2022 Kigi Chang <kigi.chang@gmail.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package funcs

// Fetcher gets successful result or error if something is wrong.
// Either, Try, and Option are built from Fetcher.
type Fetcher[T any] func() (result T, err error)

// Exists return true if successful result from this satisfies given function p.
func (f Fetcher[T]) Exists(p Predict[T]) bool {
	v, err := f()
	return err == nil && p(v)
}

// Forall return true if this has error, or successful result form this satisfies given function p.
func (f Fetcher[T]) Forall(p Predict[T]) bool {
	v, err := f()
	return err != nil || p(v)
}

// Foreach only applies successful result from this.
func (f Fetcher[T]) Foreach(op func(T)) {
	if v, err := f(); err == nil {
		op(v)
	}
}

// GetOrElse returns successful result or given default z.
func (f Fetcher[T]) GetOrElse(z T) T {
	v, err := f()
	return Cond(err == nil, v, z)
}

// BuildWithErr builds object R from given v and err.
// Apply given constructor succ to value v if err is nil,
// or apply given fail to err.
func BuildWithErr[T, R any](v T, err error,
	fail Func[error, R], succ Func[T, R]) R {

	if err == nil {
		return succ(v)
	}
	return fail(err)
}

// -----------------------------------------------------------------------------

// TODO: refactor following functions to methods when go 1.19 releases.

// Build builds object R from given Fetcher f.
// Apply given constructor succ to successful result from Fetcher f,
// or apply given fail to error from f.
func Build[T, R any](f Fetcher[T], fail Func[error, R], succ Func[T, R]) R {
	v, err := f()
	return BuildWithErr(v, err, fail, succ)
}

// BuildUnit builds R from given Fetcher f.
// Apply given constructor succ to successful result from Fetcher f,
// or invoke given fail to generate R.
func BuildUnit[T, R any](f Fetcher[T], fail Unit[R], succ Func[T, R]) R {
	if v, err := f(); err == nil {
		return succ(v)
	}
	return fail()

}
