package repository

import "github.com/skvoch/burlington-backend/tree/master/internal/models"

type Repositories interface {
	Entities() EntityRepository
	Areas() AreaRepository
}

type AreaRepository interface {
	Get(id string) (models.Area, error)
	Set(id string, area models.Area) error
}

type EntityRepository interface {
	Get(id string) (models.Entity, error)
	Set(entity models.Entity) error
	GetByAreaName(areaName string) []models.Entity
}
