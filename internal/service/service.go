package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/skvoch/burlington-backend/tree/master/internal/finder"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/qr"
	"github.com/skvoch/burlington-backend/tree/master/internal/repository/couchbase"
	"image"
)


const (
	name = "Administrator"
	pswd = "password"
	host = "192.168.0.58"
)

type Opts struct {
	Logger zerolog.Logger
	Repo   *couchbase.Repositories
}

func New(opts *Opts) *Service {
	return &Service{
		repo:   opts.Repo,
		logger: opts.Logger,
	}
}

type Service struct {
	logger zerolog.Logger

	repo *couchbase.Repositories
}

func (s *Service) FindPath(areaName string, start, target models.XYZ) (*finder.FindResult, error) {
	area, err := s.repo.Areas().Get(areaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get area: %w", err)
	}
	reader := newCellsReader(area)
	result, err := finder.FindPath(finder.FindPathParams{
		StartCell:  start,
		TargetCell: target,
		Reader:     reader,
	})

	return result, nil
}

func (s *Service) GetArea(id string) (models.Area, error) {
	area, err := s.repo.Areas().Get(id)
	if err != nil {
		return models.Area{}, fmt.Errorf("failed to get area: %w", err)
	}

	return area, nil
}

func (s *Service) SetArea(area models.Area) error {
	if area.Name == "" {
		return fmt.Errorf("area name shouldn't be empty")
	}

	if err := s.repo.Areas().Set(area.Name, area); err != nil {
		return fmt.Errorf("failed to set area: %w", err)
	}

	return nil
}

/*
func SetAreaModel(jsonObj models.Area) (string, error){
	rep, err := RepInit(rep)
	area := rep.Areas()
	id, err := area.Create(jsonObj)
	if err != nil{
		return "", fmt.Errorf("failed to upsert object, err: %w", err)
	}
	return fmt.Sprintf("object with id %v was added successfully", id), nil
}

//entity

//func GetEntity(id string)(models.Entity, error){
//
//}*/
