package main

import (
	"fmt"

	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Server struct {
	*domain.CourseList
	App *fiber.App
}

func NewServer(courses *domain.CourseList) *Server {
	engine := html.New("./tmpl", ".html")
	server := Server{
		CourseList: courses,
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	server.App.Get("/", server.ViewCourses())
	server.App.Get("/login", server.LoginPage())
	server.App.Get("/course/:num?", server.EditCourse())
	server.App.Post("/course/:1", server.InsertDB())
	return &server
}

func (s *Server) ViewCourses() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return c.Render("main", fiber.Map{
			"Courses": s.List,
		})
	}
}

func (s *Server) LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	}
}

func (s *Server) EditCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch c.Params("num") {
		case "1":
			return c.Render("createcourse", fiber.Map{
				"DB": db.DB,
			})
		case "2":
			return c.SendString("Update")
		case "3":
			return c.SendString("Delete")
		}
		return c.Render("editcourse", fiber.Map{
			"Courses": s.List,
			"DB":      s.SelectDB(1),
		})
	}
}

type lecture struct {
	Name           string
	Professor_name string
	Season         string
	Grade          string
	Credit         string
	Category       string
}

func (s *Server) SelectDB(uid int) lecture {
	var lect lecture
	rows := db.DB.QueryRow("SELECT * from lecture where uid = ?", uid)
	n := 0
	err := rows.Scan(&n, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
	if err != nil {
		fmt.Println(err)
	}
	return lect
}

func (s *Server) InsertDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`INSERT INTO lecture(name, professor_name, season, grade, credit, category) 
		VALUES(?, ?, ?, ?, ?, ?)`, lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category)
		if err != nil {
			return c.SendString("INSERT ERROR")
		}
		return c.Redirect("/course")
	}

}
