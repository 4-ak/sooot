package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func UpdateLecture() *sql.Stmt {
	query := `
	UPDATE lecture 
	SET year = $1, semester = $2, credit = $3, major = $4  
	WHERE uid = $5
	RETURNING base
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func Lecture() *sql.Stmt {
	query := `
	SELECT * 
	FROM lecture 
	WHERE uid = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func DeleteLecture() *sql.Stmt {
	query := `
	DELETE FROM lecture 
	WHERE uid = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
