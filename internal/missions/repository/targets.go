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