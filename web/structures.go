package web

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBHandler struct {
	DB *sql.DB
}

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
