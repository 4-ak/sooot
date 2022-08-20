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
	FROM lecture l, lecture_base lb
	WHERE l.base = lb.uid;
	`
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func UpdateLecture() *sql.Stmt {
	query := `
	UPDATE lecture 
	SET year = $1, semester = $2, credit = $3, major = $4  
	WHERE uid = $5
	RETURNING base
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

func InsertLecture_base() *sql.Stmt {
	query := `
	INSERT INTO lecture_base(name, professor)
	VALUES($1, $2)
	RETURNING uid
	`
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func DeleteLecture_base() *sql.Stmt {
	query := `
	DELETE FROM lecture_base WHERE uid = $1
	`
	stmt, err := DB.Prepare(query)
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
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func MajorAll() *sql.Stmt {
	query := `
	SELECT name
	FROM major;
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func ProfessorAll() *sql.Stmt {
	query := `
	SELECT p.name, m.name
	FROM professor p, major m
    WHERE p.major = m.uid;
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

func SemesterAll() *sql.Stmt {
	query := `
	SELECT name
	FROM semester;
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}

// func () *sql.Stmt {
// 	query := `

// `
// 	stmt, err := DB.Prepare(query)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return stmt
// }
