package structs

import "../../models"

type MedianHistogram [16]IntMedianHeap

func NewMedianHistogram() *MedianHistogram {
    mh := &MedianHistogram{}
    for i := 0; i < 15; i++ {
        mh[i] = *NewIntMedianHeap()
    }
    return mh
}

func (mh *MedianHistogram) AddHistogram(hg models.Histogram) {
    for i := 0; i < 15; i++ {
        mh[i].AddNum(hg[i])
    }
}

func (mh *MedianHistogram) GetMedianHistogram() (models.Histogram, error) {
    var res models.Histogram
    for i := 0; i < 15; i++ {
        mid, err := mh[i].GetMedian()
        if err != nil {
            return models.Histogram{}, err
        }
        res[i] = mid

    }
    return res, nil
}
