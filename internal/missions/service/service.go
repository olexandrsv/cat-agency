package service

import (
	"cat-agency/internal/missions/models"
	"cat-agency/internal/missions/repository"
)

type MissionService interface {
	CreateMission([]models.Target) error
}

type service struct {
	repo repository.MissionRepository
}

func New(repo repository.MissionRepository) MissionService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateMission(targets []models.Target) error {
	return s.repo.CreateMission(targets)
}