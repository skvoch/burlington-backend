package finder

import (
	"container/heap"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"math"
)

func NewCostMap(reader Reader, startCell models.XYZ) (*CostMap, error) {
	costMap := &CostMap{
		costs:          make(map[models.XYZ]int64, reader.GetCellsCount()),
		notVisitedHeap: &NotVisitedHeap{},
	}

	heap.Init(costMap.notVisitedHeap)

	neighbors, err := reader.GetCellNeighborhoods(startCell)
	if err != nil {
		return nil, err
	}

	for _, cell := range neighbors {
		cellType, err := reader.GetCellType(cell)
		if err != nil {
			return nil, err
		}

		if !cellType.IsBarrier() {
			costMap.Set(cell, math.MaxInt64)
		}
	}

	costMap.Set(startCell, 0)

	return costMap, nil
}

type CostMap struct {
	costs          map[models.XYZ]int64
	notVisitedHeap *NotVisitedHeap
}

func (c *CostMap) Set(cell models.XYZ, cost int64) {
	c.costs[cell] = cost

	heap.Push(c.notVisitedHeap, AreaCellPair{
		Cost: cost,
		Cell: cell,
	})
}

func (c *CostMap) Get(cell models.XYZ) int64 {
	cost, ok := c.costs[cell]

	if !ok {
		return math.MaxInt64
	}

	return cost
}

func (c *CostMap) HasUnvisited() bool {
	return c.notVisitedHeap.Len() > 0
}

func (c *CostMap) GetMinimalNonVisited() models.XYZ {
	pair := heap.Pop(c.notVisitedHeap).(AreaCellPair)

	return pair.Cell
}

