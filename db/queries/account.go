package queries

import (
	"database/sql"
	"fmt"

	"github.com/4-ak/sooot/db"
)

func AccountWithPass() *sql.Stmt {
	query := `
	SELECT uid, pass
	FROM account
	WHERE id=$1
	`
	stmt, err := db.DB.Prepare(query)
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
	stmt, err := db.DB.Prepare(query)
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
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	return stmt
}
