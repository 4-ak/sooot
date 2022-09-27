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
	SELECT l.year, s.name, l.credit, m.name, lb.name, p.name
	FROM lecture l, lecture_base lb, major m, professor p, semester s
	WHERE l.uid = $1 
		AND l.base = lb.uid
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

func CompareLecture() *sql.Stmt {
	query := `
	SELECT uid
	FROM lecture
	WHERE base = (
	    SELECT uid
	    FROM lecture_base
	    WHERE name = $1 
	        AND professor = (
	            SELECT uid
	            FROM professor
	            WHERE name = $2
	        )
	    )
		AND year = $3
	    AND semester = (
	        SELECT uid
	        FROM semester
	        WHERE name = $4
	    ) 
	    AND credit = $5 
	    AND major = (
	        SELECT uid
	        FROM major
	        WHERE name = $6
	    )
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
