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

func (r *repo) DeleteNotes(tx *sql.Tx, targetID int) error {
	_, err := tx.Exec("DELETE FROM notes WHERE id IN (SELECT note_id FROM target_notes WHERE target_id=$1)", targetID)
	if err != nil {
		return common.NewDatabseError(err)
	}

	_, err = tx.Exec("DELETE FROM target_notes WHERE target_id=$1", targetID)
	if err != nil {
		return common.NewDatabseError(err)
	}
	return nil
}

func (r *repo) GetNotes(targetID int) ([]models.Note, error) {
	rows, err := r.db.Query(`SELECT id, msg FROM notes INNER JOIN target_notes 
		ON target_notes.note_id=notes.id WHERE target_notes.target_id=$1`, targetID)
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))	
	}
	var notes []models.Note
	for rows.Next() {
		var note models.Note 
		if err := rows.Scan(&note.ID, &note.Message); err != nil {
			return nil, common.NewDatabseError(errors.WithStack(err))	
		}
		
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repo) UpdateNote(noteID int, msg string) error {
	_, err := r.db.Exec("UPDATE notes SET msg=$1 WHERE id=$2", msg, noteID)
	if err == sql.ErrNoRows {
		return common.NewNoRowsError(err)
	}
	if err != nil {
		return common.NewDatabseError(err)
	}

	return nil
}