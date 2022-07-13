package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func NewDB() {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		panic(err)
	}
	DB = db
}
