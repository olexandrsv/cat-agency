package models

type Cat struct {
	ID         int     `json:"id"`
	Experience int     `json:"experience"`
	Breed      string  `json:"breed"`
	Salary     float64 `json:"salary"`
}
