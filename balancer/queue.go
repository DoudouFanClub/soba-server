package balancer

type node[T any] struct {
	next *node[T]
	prev *node[T]
	val  T
}

type LinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	len  int
}

type Queue[T any] struct {
	queue *LinkedList[T]
}

func CreateNode[T any]() *node[T] {
	var zero T
	return &node[T]{next: nil, prev: nil, val: zero}
}

func NewLL[T any]() *LinkedList[T] {
	root := CreateNode[T]()
	return &LinkedList[T]{head: root, tail: root}
}

func (l *LinkedList[T]) Add(val T) {
	l.tail.val = val
	l.tail.next = &node[T]{}
	l.tail = l.tail.next
	l.len += 1
}

func (l *LinkedList[T]) RemoveFront() T {
	if l.len > 0 {
		temp := l.head.val
		l.head = l.head.next
		l.len -= 1
		return temp
	}
	var zero T
	return zero
}

func NewQueue[T any]() Queue[T] {
	ll := NewLL[T]()
	return Queue[T]{queue: ll}
}

func (q *Queue[T]) Add(val T) {
	q.queue.Add(val)
}

func (q *Queue[T]) Pop() T {
	return q.queue.RemoveFront()
}

func (q *Queue[T]) Empty() bool {
	return q.queue.len != 0
}
