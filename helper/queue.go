// MPMS（multiple producer, multiple subscriber,多生产者、多消费者）队列
package helper

type Queue[T any] struct {
	data chan T
}

func NewQueue[T any](capacity int) Queue[T] {
	return Queue[T]{
		data: make(chan T, capacity),
	}
}

func (q Queue[T]) Push(val T) {
	q.data <- val
}

func (q Queue[T]) TryPush(val T) bool {
	select {
	case q.data <- val:
		return true
	default:
		return false
	}
}

func (q Queue[T]) Pop() T {
	return <-q.data
}

func (q Queue[T]) TryPop() (T, bool) {
	select {
	case val := <-q.data:
		return val, true
	default:
		var zero T
		return zero, false
	}
}
