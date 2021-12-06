package gs

import "fmt"

type Future[T any] interface {
	fmt.Stringer

	Done() <-chan struct{}
	Completed() bool
	Get() (Try[T], bool)

	OnCompleted(func(Try[T]))
	Foreach(func(T))
	Wait()
}
