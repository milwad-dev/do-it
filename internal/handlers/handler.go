package handlers

import "database/sql"

type DBHandler struct {
	*sql.DB
}

// NewDBHandler => create a new instance from DBHandler and return it
func NewDBHandler(db *sql.DB) *DBHandler {
	return &DBHandler{db}
}
