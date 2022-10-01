package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func LectureAll() *sql.Stmt {
	query := `
	SELECT l.uid, l.year, l.base, s.name, l.credit, m.name, lb.uid, lb.name, p.name
	FROM lecture l, lecture_base lb, major m, professor p, semester s
	WHERE 
		l.base = lb.uid 
		AND l.major = m.uid 
		AND lb.professor = p.uid 
		AND l.semester = s.uid
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
		) AND s.name= $4 AND m.name = $6
	RETURNING uid;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
