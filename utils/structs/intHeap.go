package structs

type MinHeap []int

func (h MinHeap) Len() int           {
    return len(h)
}

func (h MinHeap) Less(i, j int) bool {
    return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int)      {
    h[i], h[j] = h[j], h[i]
}

func (h MinHeap) Peek() int {
    return h[0]
}

func (h *MinHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

type MaxHeap struct {
    MinHeap
}

func (h MaxHeap) Less(i, j int) bool {
    return h.MinHeap[i] > h.MinHeap[j]
}
