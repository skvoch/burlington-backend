package repository

import "github.com/skvoch/burlington-backend/tree/master/internal/models"

type Repositories interface {
	Entities() EntityRepository
	Areas() AreaRepository
}

type AreaRepository interface {
	Create(area models.Area) (string, error)
	Get(id string) (models.Area, error)
	Set(id string, area models.Area) error
}

type EntityRepository interface {
	Create(entity models.Entity) (string, error)
	Get(id string) (models.Entity, error)
	GetByAreaName(areaName string) []models.Entity
	Set(id string, entity models.Entity) error
}
