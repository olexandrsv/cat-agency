package repository

import (
	"cat-agency/internal/common"
	"cat-agency/internal/missions/models"
	"database/sql"

	"github.com/pkg/errors"
)

func (r *repo) CreateMission(targets []models.Target) error {
	tx, err := r.db.Begin()
	if err != nil {
		return common.NewDatabseError(err)
	}

	id, err := r.insertMission(tx)
	if err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(err)
	}

	if err := r.CreateTargets(tx, int(id), targets); err != nil {
		return common.NewDatabseError(err)
	}

	if err := tx.Commit(); err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(err)
	}
	return nil
}

func (r *repo) insertMission(tx *sql.Tx) (int, error) {
	res, err := tx.Exec("INSERT INTO missions (cat_id) VALUES (null)")
	if err != nil {
		return 0, common.NewDatabseError(errors.WithStack(err))
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, common.NewDatabseError(errors.WithStack(err))
	}

	return int(id), nil
}