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

func BuildFrom[T, R any](v T, err error, succ Func[T, R], fail Func[error, R]) R {
	if err == nil {
		return succ(v)
	}
	return fail(err)
}

// -----------------------------------------------------------------------------

// TODO: refactor the method when go 1.19 releases.

func Build[T, R any](f Fetcher[T], succ Func[T, R], fail Func[error, R]) R {
	v, err := f()
	return BuildFrom(v, err, succ, fail)
}

func BuildOrElse[T, R any](f Fetcher[T], z R, build Func[T, R]) R {
	v, err := f()
	if err == nil {
		return build(v)
	}
	return z
}

func BuildUnit[T, R any](f Fetcher[T], succ Func[T, R], fail Unit[R]) R {
	if v, err := f(); err == nil {
		return succ(v)
	}
	return fail()

}
