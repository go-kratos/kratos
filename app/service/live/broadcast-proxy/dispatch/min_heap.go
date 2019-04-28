package dispatch

import (
	"container/heap"
	"errors"
)

type HeapData []*HeapDataItem

type HeapDataItem struct {
	value  interface{}
	weight float64
}

func (d HeapData) Len() int {
	return len(d)
}

func (d HeapData) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d *HeapData) Push(x interface{}) {
	item := x.(*HeapDataItem)
	*d = append(*d, item)
}

func (d *HeapData) Pop() interface{} {
	old := *d
	n := len(old)
	item := old[n-1]
	*d = old[0 : n-1]
	return item
}

type MinHeapData struct {
	HeapData
}

func (d MinHeapData) Less(i, j int) bool {
	return d.HeapData[i].weight < d.HeapData[j].weight
}

type MinHeap struct {
	data MinHeapData
}

func NewMinHeap() *MinHeap {
	h := new(MinHeap)
	heap.Init(&h.data)
	return h
}

func (h *MinHeap) HeapPush(value interface{}, weight float64) {
	heap.Push(&h.data, &HeapDataItem{
		value:  value,
		weight: weight,
	})
}

func (h *MinHeap) HeapPop() (interface{}, float64, error) {
	if h.data.Len() == 0 {
		return nil, 0, errors.New("heap is empty")
	}
	item := heap.Pop(&h.data).(*HeapDataItem)
	return item.value, item.weight, nil
}

func (h *MinHeap) HeapTop() (interface{}, float64, error) {
	if h.data.Len() == 0 {
		return nil, 0, errors.New("heap is empty")
	}
	item := h.data.HeapData[0]
	return item.value, item.weight, nil
}

func (h *MinHeap) HeapLength() int {
	return h.data.Len()
}
