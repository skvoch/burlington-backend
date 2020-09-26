package service

import (
	"fmt"
	"github.com/rs/zerolog"
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
	Repo *couchbase.Repositories
}

func New(opts *Opts) *Service {
	return &Service{
		repo: opts.Repo,
		logger: opts.Logger,
	}
}

type Service struct {
	logger zerolog.Logger

	repo *couchbase.Repositories
}

func (s *Service) GetArea(id string)(models.Area, error){
	area, err := s.repo.Areas().Get(id)
	if err != nil {
		return models.Area{}, fmt.Errorf("failed to get area: %w", err)
	}

	return area,nil
}

func (s *Service) SetArea(area models.Area) error{
	if err := s.repo.Areas().Set(area.ID, area); err != nil {
		return fmt.Errorf("failed to set area: %w", err)
	}
	return nil
}

func (s *Service) CreateArea(area models.Area) error {
	if _, err := s.repo.Areas().Create(area); err != nil {
		return fmt.Errorf("failed to create area: %w", err)
	}

	return nil
}

func (s *Service) GetEntity(id string)(models.Entity, error){
	entity, err := s.repo.Entities().Get(id)
	if err != nil{
		return models.Entity{}, fmt.Errorf("Failed to get entity, err %w", err)
	}

	return entity, nil
}
func (s *Service) SetEntity(entity models.Entity) error{
	if err := s.repo.Entities().Set(entity.ID, entity); err != nil{
		return fmt.Errorf("failed to set area : %w", err)
	}
	return nil
}

func (s *Service) CreateEntity(entity models.Entity) error{
	if _, err := s.repo.Entities().Create(entity); err != nil{
		return fmt.Errorf("failed to sreate entity")
	}
	return nil
}

func (s *Service) GenerateQR(id string) (image.Image, error){
	if img, err := qr.Generate(id); err != nil{
		return nil, fmt.Errorf("Failed to generate QR: %w", err)
	}else{
		return img, nil
	}
}