package db

import (
	"database/sql"
	"fmt"
)

func AccountWithPass() *sql.Stmt {
	query := `
	SELECT uid, pass
	FROM account
	WHERE id=$1
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func RegisterAccount() *sql.Stmt {
	query := `
	INSERT INTO account(id, pass, is_cert)
	VALUES($1,$2,1)
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func AccountExists() *sql.Stmt {
	query := `
	SELECT uid, id 
	FROM account 
	WHERE uid=$1 AND id=$2
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func LectureAll() *sql.Stmt {
	query := `
	SELECT *
	FROM lecture
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func InsertLecture() *sql.Stmt {
	query := `
	INSERT INTO lecture(department, name, professor_name, semester, credit, parent) 
	VALUES($1, $2, $3, $4, $5, 0)
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func UpdateLecture() *sql.Stmt {
	query := `
	UPDATE lecture 
	SET name = $1, professor_name = $2, semester = $3, credit = $4  
	WHERE uid = $5
	`
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func ReviewAll() *sql.Stmt {
	query := `
	SELECT * 
	FROM review 
	WHERE lecture_id = $1
	`
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func DeleteReview() *sql.Stmt {
	query := `
	DELETE FROM review WHERE uid = $1
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
