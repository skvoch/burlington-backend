package models

type EntityType int

// Сущность стола/полки/чего угодно
type Entity struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Type   EntityType `json:"type"`
	CellID string     `json:"cell_id"`
	AreaID string     `json:"area_id"`
}

// Клеточка на карте
type Cell struct {
	EntityID string `json:"entity_ids"`
}

// Карта для библиотеки
type Area struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Cells [][][]Cell `json:"cells"`
}
