package repository

import (
	"cat-agency/internal/missions/models"
	"database/sql"

	"github.com/pkg/errors"
)

type MissionRepository interface {
	CreateMission([]models.Target) error
}

type repo struct {
	db *sql.DB
}

func New() MissionRepository {
	db, err := sql.Open("sqlite3", "./../cat-agency")
	if err != nil {
		panic(err)
	}

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
