package balancer

import (
	"container/list"
)

type Queue struct {
	queue *list.List
	size  int
}

func NewQueue() Queue {
	return Queue{queue: list.New(), size: 0}
}

func (q *Queue) Add(t any) {
	q.queue.PushBack(t)
	q.size += 1
}

func (q *Queue) Pop() {
	q.queue.Remove(q.queue.Front())
	q.size -= 1
}

func (q *Queue) Empty() bool {
	return q.size != 0
}
