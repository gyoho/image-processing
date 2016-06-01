package structs

import (
    "../../models"
    "../"
	"container/heap"
	"errors"
)

type SimilarHistHeap struct {
	maxHeap *HistMaxHeap
}

func NewSimilarHistHeap() *SimilarHistHeap {
    var sh SimilarHistHeap
    sh.maxHeap = &HistMaxHeap{}
    heap.Init(sh.maxHeap)
    return &sh
}

func (sh *SimilarHistHeap) BuildHeap(src models.ImageInfo, target []models.ImageInfo, n int) {
    for i := range target {
        var diff float64 = 0.0
        for j := range src.Histogram {
            diff += float64(utils.Abs(src.Histogram[j] - target[i].Histogram[j])) / float64(utils.Max(src.Histogram[j], target[i].Histogram[j]))
        }
        target[i].Delta = diff
        if sh.maxHeap.Len() > n {   // include itself
            if target[i].Delta < sh.maxHeap.Peek().Delta {
                heap.Pop(sh.maxHeap)
                heap.Push(sh.maxHeap, target[i])
            }
        } else {
            heap.Push(sh.maxHeap, target[i])
        }
    }
}

func (sh *SimilarHistHeap) GetNSimilarElem() ([]models.ImageInfo, error) {
    if sh.maxHeap.Len() == 0 {
		return []models.ImageInfo{}, errors.New("Heap is empty")
	}

    n := sh.maxHeap.Len()
    res := make([]models.ImageInfo, n - 1)
    for i := 0; i < n - 1; i++ {
        res[i] = heap.Pop(sh.maxHeap).(models.ImageInfo)
    }

    return res, nil
}
