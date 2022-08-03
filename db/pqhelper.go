package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST        = "database-1.cgyed2b704ns.ap-northeast-2.rds.amazonaws.com"
	DB_user     = "postgres"
	DB_password = "12345678"
	DB_name     = "postgres"
)

var PgDB *sql.DB

func NewPG() {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, DB_user, DB_password, DB_name)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	PgDB = db

	// test := "CREATE TABLE inventory (id serial PRIMARY KEY, name VARCHAR(50), quantity INTEGER);"
	// _, err = db.Exec(test)
	// if err != nil {
	// 	panic(err)
	// }

}

type test struct {
	id       int
	name     string
	quantity int
}

func SelectData() {
	data, err := PgDB.Query("SELECT * FROM inventory")
	if err != nil {
		panic(err)
	}
	arr := make([]test, 0)
	for data.Next() {
		var t test
		data.Scan(
			&t.id,
			&t.name,
			&t.quantity)
		arr = append(arr, t)
	}
	for i, test := range arr {
		fmt.Println(i, test)
	}
}

func InsertData() {
	sql_insert := "INSERT INTO inventory (name, quantity) VALUES ($1, $2);"
	_, err := PgDB.Exec(sql_insert, "bananas", 250)
	if err != nil {
		panic(err)
	}
}
