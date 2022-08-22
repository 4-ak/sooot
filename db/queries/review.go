package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func ReviewAll() *sql.Stmt {
	query := `
	SELECT * 
	FROM review 
	WHERE lecture_id = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func InsertReview() *sql.Stmt {
	query := `
	INSERT INTO 
	review(
		lecture_id, 
		writer, 
		beneficial_point, 
		honey_point, 
		assignment, 
		team_project, 
		presentation, 
		comment) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func UpdateReview() *sql.Stmt {
	query := `
	UPDATE review
	SET beneficial_point = $1, honey_point = $2, assignment = $3, team_project = $4, presentation = $5, comment = $6 
	WHERE uid = $7
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func Review() *sql.Stmt {
	query := `
	SELECT beneficial_point, honey_point, assignment, team_project, presentation, comment
	FROM review
	WHERE uid = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func DeleteReview() *sql.Stmt {
	query := `
	DELETE FROM review WHERE uid = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
