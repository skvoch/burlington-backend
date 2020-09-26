package service

import (
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
)

func newCellsReader(area models.Area) *cellsReader {
	return &cellsReader{
		area: area,
	}
}

type cellsReader struct {
	area models.Area
}

func (c *cellsReader) GetCellType(coords models.XYZ) (models.CellType, error) {
	if err := coords.IsValid(c.area.Width, c.area.Height, 0); err != nil {
		return 0, err
	}
	cell := c.area.Cells[coords.X][coords.Y][coords.Z]
	return cell.Type, nil
}

func (c *cellsReader) GetCellNeighborhoods(cell models.XYZ) ([]models.XYZ, error) {
	if err := cell.IsValid(c.area.Width, c.area.Height, 0); err == nil {
		out := make([]models.XYZ, 0, 5)

		left := cell
		left.X -= 1
		if err := left.IsValid(c.area.Width, c.area.Height, 0); err == nil {
			out = append(out, left)
		}

		right := cell
		right.X += 1
		if err := right.IsValid(c.area.Width, c.area.Height, 0); err == nil {
			out = append(out, right)
		}

		top := cell
		top.Y += 1
		if err := top.IsValid(c.area.Width, c.area.Height, 0); err == nil {
			out = append(out, top)
		}

		bottom := cell
		bottom.Y -= 1
		if err := bottom.IsValid(c.area.Width, c.area.Height, 0); err == nil {
			out = append(out, bottom)
		}

		return out, nil
	} else {
		return nil, err
	}
}

func (c *cellsReader) GetCellsCount() int64 {
	return int64(c.area.Width * c.area.Height)
}
