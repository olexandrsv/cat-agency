package models

type Mission struct {
	ID       int      `json:"id"`
	CatID    int      `json:"cat_id"`
	Complete bool     `json:"complete"`
	Targets  []Target `json:"targets"`
}

type Target struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Complete bool   `json:"complete"`
	Notes    []Note `json:"notes"`
}

type Note struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
