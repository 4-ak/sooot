package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func LectureAll() *sql.Stmt {
	query := `
	SELECT *
	FROM lecture l, lecture_base lb
	WHERE l.base = lb.uid;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func InsertLecture() *sql.Stmt {
	query := `
	INSERT INTO lecture(base, year, semester, credit, major)
	SELECT lb.uid, $3, s.uid, $5, m.uid
	FROM lecture_base lb, semester s, major m
	WHERE lb.name = $1 AND lb.professor = (
		SELECT p.uid
		FROM professor p
		WHERE p.name = $2
	) AND s.name= $4 AND m.name = $6;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
