package model

import (
	"database/sql"
)

type Model struct {
	db *sql.DB
}

func NewModel(db *sql.DB) Model {
	if db == nil {
		panic("DB handle is nil")
	}

	return Model{
		db: db,
	}
}
