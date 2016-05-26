package structs

type MedianHistogram [16]IntMedianHeap

func NewMedianHistogram() *MedianHistogram {
    mh := &MedianHistogram{}
    for i := 0; i < 15; i++ {
        mh[i] = *NewIntMedianHeap()
    }
    return mh
}

func (mh *MedianHistogram) AddHistogram(hg [16]int) {
    for i := 0; i < 15; i++ {
        mh[i].AddNum(hg[i])
    }
}

func (mh *MedianHistogram) GetMedianHistogram() [16]float64 {
    var res [16]float64
    for i := 0; i < 15; i++ {
        res[i] = mh[i].GetMedian()
    }
    return res
}
