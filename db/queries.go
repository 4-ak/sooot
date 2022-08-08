package db

import (
	"database/sql"
)

func AccountWithPass(id string) *sql.Row {
	query := `
	SELECT uid, pass
	FROM account
	WHERE id=$1
	`
	return DB.QueryRow(query, id)
}

func RegisterAccount(id string, pass []byte) (sql.Result, error) {
	query := `
	INSERT INTO account(id, pass, is_cert)
	VALUES($1,$2,1)
	`
	return DB.Exec(query, id, pass)
}

func AccountExists(uid, id string) *sql.Row {
	query := `
	SELECT uid, id 
	FROM account 
	WHERE uid=$1 AND id=$2
	`
	return DB.QueryRow(query, uid, id)
}

func LectureAll() (*sql.Rows, error) {
	query := `
	SELECT *
	FROM lecture
	`
	return DB.Query(query)
}

func InsertLecture(name, professor_name string, season, grade, credit, category int) error {
	query := `
	INSERT INTO lecture(name, professor_name, season, grade, credit, category) 
	VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := DB.Exec(query, name, professor_name, season, grade, credit, category)
	return err
}

func UpdateLecture(name, professor_name string, season, grade, credit, category, uid int) error {
	query := `
	UPDATE lecture 
	SET name = $1, professor_name = $2, season = $3, grade = $4, credit = $5, category = $6  
	WHERE uid = $7
	`
	_, err := DB.Exec(query, name, professor_name, season, grade, credit, category, uid)
	return err
}

func Lecture(uid int) *sql.Row {
	query := `
	SELECT * 
	FROM lecture 
	WHERE uid = $1
	`
	return DB.QueryRow(query)
}

func DeleteLecture(uid int) error {
	query := `
	DELETE FROM lecture 
	WHERE uid = $1
	`
	_, err := DB.Exec(query, uid)
	return err
}

func ReviewAll(lecture_id string) (*sql.Rows, error) {
	query := `
	SELECT * 
	FROM review 
	WHERE lecture_id = $1
	`
	return DB.Query(query, lecture_id)
}

func InsertReview(beneficial_point, honey_point, professor_point int, is_team, is_presentation bool, lecture_id, uid string) error {
	query := `
	INSERT INTO 
	review(lecture_id, beneficial_point, honey_point, professor_point, is_team, is_presentation, user_id) 
	VALUES($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := DB.Exec(query, lecture_id, beneficial_point, honey_point, professor_point, is_team, is_presentation, uid)
	return err
}

func UpdateReview(beneficial_point, honey_point, professor_point int, is_team, is_presentation bool, uid string) error {
	query := `
	UPDATE review
	SET beneficial_point = $1, honey_point = $2, professor_point = $3, is_team = $4, is_presentation = $5 
	WHERE uid = $6
	`
	_, err := DB.Exec(query, beneficial_point, honey_point, professor_point, is_team, is_presentation, uid)
	return err
}

func Review(uid int) *sql.Row {
	query := `
	SELECT beneficial_point, honey_point, professor_point, is_team, is_presentation 
	FROM review 
	WHERE uid = $1
	`
	return DB.QueryRow(query, uid)
}

func DeleteReview(uid string) error {
	query := `
	DELETE FROM review WHERE uid = $1
	`
	_, err := DB.Exec(query, uid)
	return err
}
