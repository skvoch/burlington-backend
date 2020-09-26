package models

import (
	"fmt"
	"strconv"
)

type CellType int

const (
	EmptyCellType = 0
	WallCellType  = 1
	PCCellType    = 2
	BooksCellType = 3
)

func (c CellType) IsBarrier() bool {
	return c != EmptyCellType
}

func (c CellType) GetMoveCost() int64 {
	return 1
}

// Сущность стола/полки/чего угодно
type Entity struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	CellID   string `json:"cell_id"`
	AreaName string `json:"area_id"`
}

// Клеточка на карте
type Cell struct {
	EntityID string   `json:"entity_ids"`
	Type     CellType `json:"type"`
}

// Карта для библиотеки
type Area struct {
	Name   string `json:"name"` // key
	Width  int64  `json:"width"`
	Height int64  `json:"height"`

	Cells [][][]Cell `json:"cells"`
}

// Координаты клетки
type XYZ struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z int64 `json:"z"`
}

func (p *XYZ) IsValid(width int64, height int64, floors int64) error {
	if (p.X < width && p.Y < height && p.Z < floors) && (p.X >= 0 && p.Y >= 0 && p.Z >= 0) {
		return nil
	} else {
		return fmt.Errorf("invalid cell coordinates, cell: %v", *p)
	}
}

func (p *XYZ) ToString() string {
	return strconv.FormatInt(p.X, 10) + "." + strconv.FormatInt(p.Y, 10) + "." + strconv.FormatInt(p.Z, 10)
}
