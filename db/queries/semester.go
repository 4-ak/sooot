package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func SemesterAll() *sql.Stmt {
	query := `
	SELECT name
	FROM semester;
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
