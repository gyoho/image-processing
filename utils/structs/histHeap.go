package structs

import "../../models"

type HistMinHeap []models.ImageInfo

func (hmh HistMinHeap) Len() int {
    return len(hmh)
}

func (hmh HistMinHeap) Less(i, j int) bool {
    return hmh[i].Delta < hmh[j].Delta
}

func (hmh HistMinHeap) Swap(i, j int) {
    hmh[i], hmh[j] = hmh[j], hmh[i]
}

func (hmh HistMinHeap) Peek() models.ImageInfo {
    return hmh[0]
}

func (hmh *HistMinHeap) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    *hmh = append(*hmh, x.(models.ImageInfo))
}

func (hmh *HistMinHeap) Pop() interface{} {
    old := *hmh
    n := len(old)
    x := old[n-1]
    *hmh = old[0 : n-1]
    return x
}


type HistMaxHeap struct {
    HistMinHeap
}

func (h HistMaxHeap) Less(i, j int) bool {
    return h.HistMinHeap[i].Delta > h.HistMinHeap[j].Delta
}
