package finder

import "github.com/skvoch/burlington-backend/tree/master/internal/models"


type MovementHistory map[models.XYZ]models.XYZ


type Reader interface {
	GetCellType(cell models.XYZ) (models.CellType, error)
	GetCellNeighborhoods(cell models.XYZ) ([]models.XYZ, error)
	GetCellsCount() int64
}


type FindResult struct {
	TotalCost int64        `json:"total_cost"`
	Path      []models.XYZ `json:"path"`
	Found     bool         `json:"found"`
}

type FindPathParams struct {
	StartCell  models.XYZ
	TargetCell models.XYZ

	Reader Reader
}

func FindPath(params FindPathParams) (*FindResult, error) {
	movementHistory, err := calculateMovementHistory(params.StartCell, params.Reader)
	if err != nil {
		return nil, err
	}

	out, err := getFindPathResult(movementHistory, params)
	if err != nil {
		return nil, err
	}
	return out, nil
}


func calculateMovementHistory(startCell models.XYZ, reader Reader) (MovementHistory, error) {
	history := make(MovementHistory, reader.GetCellsCount())
	costs, err := NewCostMap(reader, startCell)

	if err != nil {
		return nil, err
	}

	for costs.HasUnvisited() {
		unvisitedCell := costs.GetMinimalNonVisited()

		if err != nil {
			return nil, err
		}

		neighbors, err := reader.GetCellNeighborhoods(unvisitedCell)

		if err != nil {
			return nil, err
		}

		var distanceToNeighbor int64
		var cellType models.CellType

		unvisitedCellCost := costs.Get(unvisitedCell)

		for _, cell := range neighbors {
			cellType, err = reader.GetCellType(cell)

			if err != nil {
				return nil, err
			}
			if cellType.IsBarrier() {
				continue
			}

			distanceToNeighbor = unvisitedCellCost + cellType.GetMoveCost()

			if distanceToNeighbor < costs.Get(cell) {
				costs.Set(cell, distanceToNeighbor)

				history[cell] = unvisitedCell
			}
		}
	}

	return history, nil
}

func getFindPathResult(history MovementHistory, params FindPathParams) (*FindResult, error) {
	path := getPath(history, params.TargetCell)

	totalCost, err := getTotalCost(params.Reader, path)
	if err != nil {
		return nil, err
	}
	isFound := isPathFound(path, params.StartCell, params.TargetCell)

	return &FindResult{
		Path:      path,
		TotalCost: totalCost,
		Found:     isFound,
	}, nil
}

func getPath(history MovementHistory, target models.XYZ) []models.XYZ {

	out := make([]models.XYZ, 0)
	out = append(out, target)
	current := target

	for {
		cell, ok := history[current]
		if !ok {
			break
		}

		current = cell
		out = append(out, cell)
	}

	return reverse(out)
}

func reverse(path []models.XYZ) []models.XYZ {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func getTotalCost(reader Reader, path []models.XYZ) (int64, error) {
	out := int64(0)

	for _, cell := range path {
		cellType, err := reader.GetCellType(cell)
		if err != nil {
			return 0, err
		}
		out = cellType.GetMoveCost()
	}

	return out, nil
}

func isPathFound(path []models.XYZ, from models.XYZ, to models.XYZ) bool {

	if len(path) > 2 {
		return path[0] == from && path[len(path)-1] == to
	}

	return false
}