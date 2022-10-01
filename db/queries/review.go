package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func ReviewAll() *sql.Stmt {
	query := `
	SELECT r.uid, r.lecture , r.writer, r.beneficial, r.honey, r.assignment, r.team_project, r.presentation, r.comment, r.created_at, lb.name
	FROM review r, lecture_base lb, lecture l
	WHERE r.lecture = $1 AND l.uid = r.lecture AND l.base = lb.uid;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func ReviewList() *sql.Stmt {
	query := `
	SELECT writer, beneficial, honey, assignment. team_project, presentation
	FROM review
	WHERE lecture = $1
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
	//모든 별점을 더해서 보여주는 방식은?
	//이게 reviewall로 대체?
	//lecture 외래키를 lecture_base의 uid로 옮겨야 할듯?
}

func ReviewContent() *sql.Stmt {
	query := `
	SELECT writer, beneficial, honey, assignment. team_project, presentation, comment, created_at
	FROM review
	WHERE lecture = $1
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
		lecture,
		writer, 
		beneficial, 
		honey, 
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
	SET beneficial = $1, 
		honey = $2, 
		assignment = $3,
		team_project = $4, 
		presentation = $5,
		comment = $6 
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
	SELECT beneficial, honey, assignment, team_project, presentation, comment
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
