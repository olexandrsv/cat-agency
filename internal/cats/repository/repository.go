package repository

import (
	"cat-agency/internal/cats/models"
	"cat-agency/internal/common"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type CatsRepository interface {
	CreateCat(int, string, float64) error
	UpdateCat(int, float64) error
	DeleteCat(int) error
	GetCats() ([]models.Cat, error)
	GetCat(int) (models.Cat, error)
}

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) CatsRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateCat(experience int, breed string, salary float64) error {
	_, err := r.db.Exec("INSERT INTO Cats (experience, breed, salary) VALUES ($1, $2, $3)", experience, breed, salary)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}
	return nil
}

func (r *repo) DeleteCat(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec("DELETE FROM cats WHERE id=$1", id)
	if err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(errors.WithStack(err))
	}

	_, err = tx.Exec("UPDATE missions SET cat_id=null WHERE cat_id=$1", id)
	if err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(errors.WithStack(err))
	}

	if err = tx.Commit(); err != nil {
		err = executeRollback(tx, err)
		return common.NewDatabseError(errors.WithStack(err))
	}
	return nil
}

func executeRollback(tx *sql.Tx, err error) error {
	if txErr := tx.Rollback(); txErr != nil {
		err = errors.Wrap(err, txErr.Error())
	}
	return err
}

func (r *repo) UpdateCat(catID int, salary float64) error {
	_, err := r.db.Exec("UPDATE cats SET salary=$1 WHERE id=$2", salary, catID)
	if err != nil {
		return common.NewDatabseError(errors.WithStack(err))
	}
	return nil
}

func (r *repo) GetCats() ([]models.Cat, error){
	rows, err := r.db.Query("SELECT id, experience, breed, salary FROM cats")
	if err == sql.ErrNoRows {
		return nil, common.NewNoRowsError(errors.WithStack(err))
	}
	if err != nil {
		return nil, common.NewDatabseError(errors.WithStack(err))
	}

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		if err := rows.Scan(&cat.ID, &cat.Experience, &cat.Breed, &cat.Salary); err != nil {
			return nil, common.NewDatabseError(errors.WithStack(err))
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

func (r *repo) GetCat(id int) (models.Cat, error) {
	row := r.db.QueryRow("SELECT id, experience, breed, salary FROM cats WHERE id=$1", id)
	
	var cat models.Cat
	if err := row.Scan(&cat.ID, &cat.Experience, &cat.Breed, &cat.Salary); err != nil {
		return models.Cat{}, common.NewDatabseError(errors.WithStack(err))
	}
	return cat, nil
}