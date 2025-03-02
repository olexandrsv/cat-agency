package repository

import (
	"cat-agency/internal/missions/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type MissionRepository interface {
	CreateMission([]models.Target) error
	DeleteMission(int) error
	CatExists(int) (bool, error)
	AssignMission(int, int) error
	UpdateMission(int) error
	UpdateTarget(int) error
	DeleteTarget(int) error
	CreateTarget(int, models.Target) error
	IsTargetCompleted(int) (bool, error)
	IsMissionCompleted(int) (bool, error)
	IsMissionAssigned(int) (bool, error)
	CountMissionTargets(int) (int, error)
	GetMissions() ([]models.Mission, error)
	GetMission(int) (models.Mission, error)
	GetTargets(int) ([]models.Target, error)
	UpdateNote(int, string) error
}

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) MissionRepository {
	return &repo{
		db: db,
	}
}

func executeRollback(tx *sql.Tx, err error) error {
	if txErr := tx.Rollback(); txErr != nil {
		err = errors.Wrap(err, txErr.Error())
	}
	return err
}
