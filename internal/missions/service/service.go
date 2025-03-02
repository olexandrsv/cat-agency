package service

import (
	"cat-agency/internal/common"
	"cat-agency/internal/missions/models"
	"cat-agency/internal/missions/repository"
)

type MissionService interface {
	CreateMission([]models.Target) error
	DeleteMission(int) error
	AssignMission(int, int) error
	UpdateMission(int) error
	UpdateTarget(int) error
	DeleteTarget(int) error
	CreateTarget(int, models.Target) error
	GetMissions() ([]models.Mission, error)
	GetMission(int) (models.Mission, error)
	UpdateNote(int, int, int, string) error
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
	if len(targets) < 1 {
		return common.NewFewMissionTargetsError()
	}
	if len(targets) > 3 {
		return common.NewManyMissionTargetsError()
	}

	if err := checkIfTargetsUnique(targets); err != nil {
		return err
	}
	return s.repo.CreateMission(targets)
}

func (s *service) DeleteMission(missionID int) error {
	assigned, err := s.repo.IsMissionAssigned(missionID)
	if err != nil {
		return err
	}
	if assigned {
		return common.NewMissionAssignedError()
	}
	return s.repo.DeleteMission(missionID)
}

func (s *service) AssignMission(missionID, catID int) error {
	exists, err := s.repo.CatExists(catID)
	if err != nil {
		return err
	}
	if !exists {
		return common.NewWrongCatIDError(catID)
	}

	return s.repo.AssignMission(missionID, catID)
}

func (s *service) UpdateMission(missionID int) error {
	return s.repo.UpdateMission(missionID)
}

func (s *service) UpdateTarget(missionID int) error {
	return s.repo.UpdateTarget(missionID)
}

func (s *service) DeleteTarget(targetID int) error {
	completed, err := s.repo.IsTargetCompleted(targetID)
	if err != nil {
		return err
	}
	if completed {
		return common.NewTargetCompletedError(targetID)
	}

	return s.repo.DeleteTarget(targetID)
}

func (s *service) CreateTarget(missionID int, target models.Target) error {
	completed, err := s.repo.IsMissionCompleted(missionID)
	if err != nil {
		return err
	}
	if completed {
		return common.NewMissionCompletedError(missionID)
	}

	targets, err := s.repo.GetTargets(missionID)
	if err != nil {
		return err
	}
	if len(targets) >= 3 {
		return common.NewManyMissionTargetsError()
	}
	targets = append(targets, target)

	if err := checkIfTargetsUnique(targets); err != nil {
		return err
	}

	return s.repo.CreateTarget(missionID, target)
}

func checkIfTargetsUnique(targets []models.Target) error {
	for i := 0; i < len(targets); i++ {
		for j := 0; j < len(targets); j++ {
			if i == j {
				continue
			}
			if targets[i].Name == targets[j].Name && targets[i].Country == targets[j].Country {
				return common.NewTargetsDublicateError()
			}
		}
	}
	return nil
}

func (s *service) GetMissions() ([]models.Mission, error) {
	return s.repo.GetMissions()
}

func (s *service) GetMission(missionID int) (models.Mission, error) {
	return s.repo.GetMission(missionID)
}

func (s *service) UpdateNote(missionID, targetID, noteID int, msg string) error {
	completed, err := s.repo.IsMissionCompleted(missionID)
	if err != nil {
		return err
	}
	if completed {
		return common.NewMissionCompletedError(missionID)
	}

	completed, err = s.repo.IsTargetCompleted(targetID)
	if err != nil {
		return err
	}
	if completed {
		return common.NewTargetCompletedError(targetID)
	}

	if err := s.repo.UpdateNote(noteID, msg); err != nil {
		return err
	}

	return nil
}
