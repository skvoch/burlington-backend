package finder

import (
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
)

type AreaCellPair struct {
	Cost int64
	Cell models.XYZ
}

type NotVisitedHeap []AreaCellPair

func (h NotVisitedHeap) Len() int {

	return len(h)
}

func (h NotVisitedHeap) Less(i, j int) bool {
	return h[i].Cost < h[j].Cost
}

func (h NotVisitedHeap) Swap(i, j int) {
	tmp := h[i]
	h[i] = h[j]
	h[j] = tmp
}

func (h *NotVisitedHeap) Push(i interface{}) {
	*h = append(*h, i.(AreaCellPair))
}

func (h *NotVisitedHeap) Pop() interface{} {

	pre := *h
	length := len(pre)
	x := pre[length-1]
	*h = pre[0 : length-1]

	return x
}

