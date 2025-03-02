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

func (r *repo) DeleteMission(missionID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return common.NewDatabseError(err)
	}

	err = r.deleteMission(tx, missionID)
	if err != nil {
		err = executeRollback(tx, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		err = executeRollback(tx, err)
		return err
	}
	return nil
}

func (r *repo) IsMissionCompleted(missionID int) (bool, error) {
	row := r.db.QueryRow("SELECT complete FROM missions WHERE id=$1", missionID)
	var b bool
	if err := row.Scan(&b); err != nil {
		return false, common.NewDatabseError(err)
	}
	return b, nil
}

func (r *repo) IsMissionAssigned(missionID int) (bool, error) {
	row := r.db.QueryRow("SELECT cat_id FROM missions WHERE id=$1", missionID)
	var s sql.NullString
	if err := row.Scan(&s); err != nil {
		return false, common.NewDatabseError(err)
	}
	return s.Valid, nil
}

func (r *repo) deleteMission(tx *sql.Tx, missionID int) error {
	_, err := tx.Exec(`DELETE FROM notes WHERE id IN (
		SELECT note_id FROM target_notes WHERE target_id IN (
			SELECT target_id FROM mission_targets WHERE mission_id=$1
		)
	)`, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM target_notes WHERE target_id IN (
		SELECT target_id FROM mission_targets WHERE mission_id=$1
	)`, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM targets WHERE id IN (
		SELECT target_id FROM mission_targets WHERE mission_id=$1
	)`, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM mission_targets WHERE mission_id=$1`, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM missions WHERE id=$1`, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	return nil
}

func (r *repo) CatExists(catID int) (bool, error) {
	row := r.db.QueryRow("SELECT count(id) FROM cats WHERE id=$1", catID)
	var count int
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		return false, common.NewNoRowsError(err)
	}
	if err != nil {
		return false, common.NewDatabseError(errors.WithStack(err))
	}

	return count != 0, nil
}

func (r *repo) AssignMission(missionID, catID int) error {
	_, err := r.db.Exec("UPDATE missions SET cat_id=$1 WHERE id=$2", catID, missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	return nil
}

func (r *repo) UpdateMission(missionID int) error {
	_, err := r.db.Exec("UPDATE missions SET complete=true WHERE id=$1", missionID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}
	return nil
}

func (r *repo) CountMissionTargets(missionID int) (int, error) {
	row := r.db.QueryRow("SELECT count(mission_id) FROM mission_targets WHERE mission_id=$1", missionID)
	var count int
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		return 0, common.NewNoRowsError(err)
	}
	if err != nil {
		return 0, common.NewDatabseError(errors.WithStack(err))
	}

	return count, nil
}

func (r *repo) GetMissions() ([]models.Mission, error) {
	rows, err := r.db.Query("SELECT id, cat_id, complete FROM missions")
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))
	}
	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		var catID sql.NullInt64
		if err := rows.Scan(&mission.ID, &catID, &mission.Complete); err != nil {
			return nil, common.NewDatabseError(errors.WithStack(err))
		}
		if catID.Valid {
			mission.CatID = int(catID.Int64)
		}

		targets, err := r.GetTargets(mission.ID)
		if err != nil {
			return nil, err
		}
		mission.Targets = targets

		missions = append(missions, mission)
	}

	return missions, nil
}

func (r *repo) GetMission(missionID int) (models.Mission, error) {
	row := r.db.QueryRow("SELECT id, cat_id, complete FROM missions WHERE id=$1", missionID)
	var mission models.Mission
	var catID sql.NullInt64
	if err := row.Scan(&mission.ID, &catID, &mission.Complete); err != nil {
		return models.Mission{}, common.NewDatabseError(err)
	}
	if catID.Valid {
		mission.CatID = int(catID.Int64)
	}

	targets, err := r.GetTargets(mission.ID)
	if err != nil {
		return models.Mission{}, err
	}
	mission.Targets = targets

	return mission, nil
}
