// Copyright (C) 2025 wangyusong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package heap

import (
	"container/heap"
	"sort"
)

type Heap[T any] struct {
	inner *innerHeap[T]
	limit int
}

func New[T any](data []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{
		inner: newInnerHeap(data, less),
		limit: cap(data),
	}
	heap.Init(h.inner)

	return h
}

func (h *Heap[T]) TryEvictPush(x T) {
	switch {
	case h.Len() < h.limit:
	case h.inner.less(h.Peek(), x):
		h.Pop()
	default:
		return
	}

	h.Push(x)
}

func (h *Heap[T]) Push(x T) {
	heap.Push(h.inner, x)
}

func (h *Heap[T]) Pop() T {
	return heap.Pop(h.inner).(T)
}

func (h *Heap[T]) PopLast() T {
	return heap.Remove(h.inner, h.Len()-1).(T)
}

func (h *Heap[T]) Peek() T {
	if h.Len() == 0 {
		var zero T

		return zero
	}

	return h.inner.data[0]
}

func (h *Heap[T]) Len() int {
	return h.inner.Len()
}

func (h *Heap[T]) Cap() int {
	return h.limit
}

func (h *Heap[T]) Slice() []T {
	return h.inner.data
}

func (h *Heap[T]) DESCSort() {
	sort.Slice(h.inner.data, func(i, j int) bool {
		return !h.inner.less(h.inner.data[i], h.inner.data[j])
	})
}

type innerHeap[T any] struct {
	data []T
	less func(a, b T) bool
}

func newInnerHeap[T any](data []T, less func(a, b T) bool) *innerHeap[T] {
	return &innerHeap[T]{
		data: data,
		less: less,
	}
}

func (h *innerHeap[T]) Len() int {
	return len(h.data)
}

func (h *innerHeap[T]) Less(i, j int) bool {
	return h.less(h.data[i], h.data[j])
}

func (h *innerHeap[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *innerHeap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *innerHeap[T]) Pop() any {
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[:n-1]

	return x
}
