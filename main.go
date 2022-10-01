package main

import (
	"github.com/4-ak/sooot/db"
)

func main() {
	db.NewDB()
	defer db.DB.Close()

	server := NewServer()

	server.App.Listen(":8080")
}
