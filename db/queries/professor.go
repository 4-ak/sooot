package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func ProfessorAll() *sql.Stmt {
	query := `
	SELECT p.name, m.name
	FROM professor p, major m
    WHERE p.major = m.uid;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
