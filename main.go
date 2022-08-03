package main

import (
	"github.com/4-ak/sooot/db"
)

func main() {
	db.NewDB()

	db.NewPG()
	db.InsertData()
	db.SelectData()
	defer db.DB.Close()

	server := NewServer()

	server.App.Listen(":8080")
}
