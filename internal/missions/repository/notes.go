package repository

import (
	"bytes"
	"cat-agency/internal/common"
	"cat-agency/internal/missions/models"
	"database/sql"
	"strconv"

	"github.com/pkg/errors"
)

func (r *repo) CreateNotes(tx *sql.Tx, targetID int, notes []models.Note) error {
	notesIDs, err := r.insertNotes(tx, notes)
	if err != nil {
		return err
	}

	err = r.linkNotesToTargets(tx, targetID, notesIDs)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) insertNotes(tx *sql.Tx, notes []models.Note) ([]int, error) {
	var buff bytes.Buffer
	buff.WriteString("INSERT INTO notes (msg) VALUES ")
	for i, note := range notes {
		buff.WriteString("('")
		buff.WriteString(note.Message)
		buff.WriteString("')")
		if i != len(notes)-1 {
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

	notesIDs := make([]int, 0, len(notes))
	for i := int(lastID) - len(notes); i < int(lastID); i++ {
		notesIDs = append(notesIDs, i+1)
	}

	return notesIDs, nil
}

func (r *repo) linkNotesToTargets(tx *sql.Tx, targetID int, notesIDs []int) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT INTO target_notes (target_id, note_id) VALUES ")
	for i, noteID := range notesIDs {
		buff.WriteString("(")
		buff.WriteString(strconv.Itoa(targetID))
		buff.WriteString(", ")
		buff.WriteString(strconv.Itoa(noteID))
		buff.WriteString(")")
		if i != len(notesIDs)-1 {
			buff.WriteString(",")
		}
	}

	_, err := tx.Exec(buff.String())
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	return nil
}