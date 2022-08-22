package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func InsertLecture_base() *sql.Stmt {
	query := `
	INSERT INTO lecture_base(name, professor)
	VALUES($1, $2)
	RETURNING uid
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func UpdateLecture_base() *sql.Stmt {
	query := `
	UPDATE lecture_base
	SET name = $1, professor = $2
	WHERE uid = $3
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func DeleteLecture_base() *sql.Stmt {
	query := `
	DELETE FROM lecture_base WHERE uid = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func Lecture_baseAll() *sql.Stmt {
	query := `
	SELECT lb.name, p.name, m.name
	FROM lecture_base lb, professor p, major m
    WHERE lb.professor = p.uid AND p.major = (
        SELECT m.uid
    )
    ORDER BY lb.name ASC;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
