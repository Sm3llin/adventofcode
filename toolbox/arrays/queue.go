package arrays

import "iter"

type Queue[T comparable] struct {
	items []T
}

func NewQueue[T comparable](items []T) *Queue[T] {
	return &Queue[T]{
		items: items,
	}
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() T {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(q.items) > 0 {
			yield(q.Pop())
		}
	}
}

func (q *Queue[T]) Exists(needle T) bool {
	for _, item := range q.items {
		if item == needle {
			return true
		}
	}
	return false
}
