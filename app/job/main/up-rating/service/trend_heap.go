package service

import (
	"container/heap"

	"go-common/app/job/main/up-rating/model"
)

// Heap for diff topK
type Heap interface {
	Put(*model.Diff)
	Result() []*model.Diff
}

// AscHeap for asc
type AscHeap struct {
	heap  []*model.Diff
	ctype int
}

// Put put diff to heap
func (a *AscHeap) Put(diff *model.Diff) {
	heap.Push(a, diff)
	if a.Len() > 200 {
		heap.Pop(a)
	}
}

// Result get result
func (a *AscHeap) Result() []*model.Diff {
	return a.heap
}

// Len len
func (a *AscHeap) Len() int { return len(a.heap) }

// Less less
func (a *AscHeap) Less(i, j int) bool {
	return a.heap[i].GetScore(a.ctype) < a.heap[j].GetScore(a.ctype)
}

// Swap swap
func (a *AscHeap) Swap(i, j int) { a.heap[i], a.heap[j] = a.heap[j], a.heap[i] }

// Push push to heap
func (a *AscHeap) Push(x interface{}) {
	a.heap = append(a.heap, x.(*model.Diff))
}

// Pop pop from heap
func (a *AscHeap) Pop() interface{} {
	old := a.heap
	n := len(old)
	x := old[n-1]
	a.heap = old[0 : n-1]
	return x
}

// DescHeap for desc
type DescHeap struct {
	heap  []*model.Diff
	ctype int
}

// Put to descHeap
func (d *DescHeap) Put(diff *model.Diff) {
	heap.Push(d, diff)
	if d.Len() > 200 {
		heap.Pop(d)
	}
}

// Result desc heap result
func (d *DescHeap) Result() []*model.Diff {
	return d.heap
}

// Len len
func (d *DescHeap) Len() int { return len(d.heap) }

// Less less
func (d *DescHeap) Less(i, j int) bool {
	return d.heap[i].GetScore(d.ctype) > d.heap[j].GetScore(d.ctype)
}

// Swap swap
func (d *DescHeap) Swap(i, j int) { d.heap[i], d.heap[j] = d.heap[j], d.heap[i] }

// Push push to desc heap
func (d *DescHeap) Push(x interface{}) {
	d.heap = append(d.heap, x.(*model.Diff))
}

// Pop pop from desc heap
func (d *DescHeap) Pop() interface{} {
	old := d.heap
	n := len(old)
	x := old[n-1]
	d.heap = old[0 : n-1]
	return x
}
