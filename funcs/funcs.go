package funcs

type Func[T, R any] func(T) R
type Predict[T any] func(T) bool

type Unit[R any] func() R
type Condition func() bool

type Check[T, R any] func(T) (R, bool)
type Try[T, R any] func(T) (R, error)

type Transform[T, R any] func(T, bool) R
type Recover[T, R any] func(T, error) R

func Self[T any](v T) T {
	return v
}

func Id[T any](v T) Unit[T] {
	return func() T {
		return v
	}
}

func AndThen[T, U, R any](f Func[T, U], g Func[U, R]) Func[T, R] {
	return func(v T) R {
		return g(f(v))
	}
}

func UnitAndThen[T, R any](f Unit[T], g Func[T, R]) Unit[R] {
	return func() R {
		return g(f())
	}
}

func Compose[T, U, R any](f Func[U, R], g Func[T, U]) Func[T, R] {
	return func(v T) R {
		return f(g(v))
	}
}

func ComposeUnit[T, R any](f Func[T, R], g Unit[T]) Unit[R] {
	return func() R {
		return f(g())
	}
}

func CheckAndTransform[T, U, R any](f1 Check[T, U], f2 Transform[U, R]) Func[T, R] {
	return func(v T) R {
		return f2(f1(v))
	}
}

func TryAndRecover[T, U, R any](f1 Try[T, U], f2 Recover[U, R]) Func[T, R] {
	return func(v T) R {
		return f2(f1(v))
	}
}

func Cond[T any](ok bool, succ T, fail T) T {
	if ok {
		return succ
	}
	return fail
}

func ConfFunc[T any](p Condition, succ Unit[T], fail Unit[T]) T {
	if p() {
		return succ()
	}
	return fail()
}
