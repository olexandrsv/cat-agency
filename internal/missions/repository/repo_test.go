package repository

import (
	"cat-agency/internal/missions/models"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateMission(t *testing.T) {
	r := New()
	m := []models.Target{
		models.Target{
			Name: "target1",
			Country: "country1",
			Notes: []models.Note{
				models.Note{
					Message: "m1",
				},
				models.Note{
					Message: "m2",
				},
			},
		},
		models.Target{
			Name: "target2",
			Country: "country2",
		},
	}

	if err := r.CreateMission(m); err != nil {
		t.Log(err)
	}
}