package main

import (
	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/domain"
)

func main() {
	db.NewDB()
	defer db.DB.Close()

	server := NewServer(&domain.CourseList{})

	server.App.Listen(":8080")
}
