package repository

import (
	"bytes"
	"cat-agency/internal/common"
	"cat-agency/internal/missions/models"
	"database/sql"
	"strconv"

	"github.com/pkg/errors"
)

func (r *repo) CreateTargets(tx *sql.Tx, missionID int, targets []models.Target) error {
	targetIDs, err := r.insertTargets(tx, targets)
	if err != nil {
		return err
	}

	if err := r.linkTargetsToMission(tx, missionID, targetIDs); err != nil {
		return err
	}

	return nil
}

func (r *repo) insertTargets(tx *sql.Tx, targets []models.Target) ([]int, error) {
	var buff bytes.Buffer
	buff.WriteString("INSERT INTO targets (name, country, complete) VALUES ")
	for i, target := range targets {
		buff.WriteString("('")
		buff.WriteString(target.Name)
		buff.WriteString("', '")
		buff.WriteString(target.Country)
		buff.WriteString("', false)")
		if i != len(targets)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString(";")

	res, err := tx.Exec(buff.String())
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))
	}

	targetIDs := make([]int, 0, len(targets))
	for i, j := int(lastID)-len(targets), 0; i < int(lastID); i, j = i+1, j+1 {
		targetIDs = append(targetIDs, i+1)
	}

	for i, targetID := range targetIDs {
		if len(targets[i].Notes) == 0 {
			continue
		}
		if err := r.CreateNotes(tx, targetID, targets[i].Notes); err != nil {
			return nil, err
		}
	}

	return targetIDs, nil
}

func (r *repo) linkTargetsToMission(tx *sql.Tx, missionID int, targetIDs []int) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT INTO mission_targets (mission_id, target_id) VALUES ")
	for i, targetID := range targetIDs {
		buff.WriteString("(")
		buff.WriteString(strconv.Itoa(missionID))
		buff.WriteString(", ")
		buff.WriteString(strconv.Itoa(targetID))
		buff.WriteString(")")
		if i != len(targetIDs)-1 {
			buff.WriteString(",")
		}
	}

	_, err := tx.Exec(buff.String())
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	return nil
}

func (r *repo) missionTargetsIDs(missionID int) ([]int, error) {
	rows, err := r.db.Query("SELECT target_id FROM mission_targets WHERE mission_id=$1", missionID)
	if err == sql.ErrNoRows {
		return nil, common.NewNoRowsError(errors.WithStack(err))
	}
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))
	}

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, common.NewDatabseError(errors.WithStack(err))
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *repo) UpdateTarget(targetID int) error {
	_, err := r.db.Exec("UPDATE targets SET complete=true WHERE id=$1", targetID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}
	return nil
}

func (r *repo) IsTargetCompleted(targetID int) (bool, error) {
	row := r.db.QueryRow("SELECT complete FROM targets WHERE id=$1", targetID)
	var b bool
	err := row.Scan(&b); 
	if err == sql.ErrNoRows{
		return false, common.NewNoRowsError(err)
	}
	if err != nil {
		return false, common.NewDatabseError(errors.WithStack(err))
	}
	return b, nil
}

func (r *repo) DeleteTarget(targetID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return common.NewDatabseError(err)
	}

	if err := r.deleteTarget(tx, targetID); err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(err)
	}

	if err := tx.Commit(); err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(err)
	}

	return nil
}

func (r *repo) deleteTarget(tx *sql.Tx, targetID int) error {
	_, err := tx.Exec(`DELETE FROM notes WHERE id IN (
		SELECT note_id FROM target_notes WHERE target_id=$1
	)`, targetID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM target_notes WHERE target_id=$1`, targetID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM targets WHERE id=$1`, targetID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec(`DELETE FROM mission_targets WHERE target_id=$1`, targetID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	return nil
}

func (r *repo) CreateTarget(missionID int, target models.Target) error {
	tx, err := r.db.Begin()
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}
	if err := r.CreateTargets(tx, missionID, []models.Target{target}); err != nil {
		err = executeRollback(tx, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		err = executeRollback(tx, err)
		return err
	}

	return nil
}

func (r *repo) GetTargets(missionID int) ([]models.Target, error){
	rows, err := r.db.Query(`SELECT id, name, country, complete FROM targets INNER JOIN mission_targets 
		ON targets.id=mission_targets.target_id WHERE mission_targets.mission_id=$1`, missionID)
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))	
	}
	var targets []models.Target
	for rows.Next() {
		var target models.Target 
		if err := rows.Scan(&target.ID, &target.Name, &target.Country, &target.Complete); err != nil {
			return nil, common.NewDatabseError(errors.WithStack(err))	
		}

		notes, err := r.GetNotes(target.ID)
		if err != nil {
			return nil, err
		}
		target.Notes = notes
		
		targets = append(targets, target)
	}

	return targets, nil
}