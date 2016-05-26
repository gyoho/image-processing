package structs

import "container/heap"

type IntMedian struct {
	minHeap *MinHeap
	maxHeap *MaxHeap
}

func NewIntMedian() *IntMedian {
	var im IntMedian
	im.minHeap = &MinHeap{}
	heap.Init(im.minHeap)

	im.maxHeap = &MaxHeap{[]int{}}
	heap.Init(im.maxHeap)

	return &im
}

func (im *IntMedian) AddNum(num int) {
	heap.Push(im.maxHeap, num)

	// Size requirement
	// max-heap can contain 1 more element than min-heap
	if im.maxHeap.Len() > im.minHeap.Len()+1 {
		heap.Push(im.minHeap, heap.Pop(im.maxHeap))
	}

	// Order requirement
	// every element in the max-heap to be less than or equal to all the elements in min-heap
	if im.minHeap.Len() != 0 && im.maxHeap.Peek() > im.minHeap.Peek() {
		temp := heap.Pop(im.maxHeap)
		heap.Push(im.maxHeap, heap.Pop(im.minHeap))
		heap.Push(im.minHeap, temp)
	}
}

func (im IntMedian) FindMedian() float64 {
	if im.maxHeap.Len() == im.minHeap.Len() {
		return float64((im.maxHeap.Peek() + im.minHeap.Peek()) / 2.0)
	}

	return float64(im.maxHeap.Peek())
}
