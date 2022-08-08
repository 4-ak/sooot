package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST        = "database-1.cgyed2b704ns.ap-northeast-2.rds.amazonaws.com"
	DB_user     = "postgres"
	DB_password = "12345678"
	DB_name     = "postgres"
)

var DB *sql.DB

func NewDB() {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, DB_user, DB_password, DB_name)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	DB = db
}
