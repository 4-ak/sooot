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
	server.App.Get("/registration", server.RegistrationPage())
	server.App.Post("/registration", server.Registration())
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

func (s *Server) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}
}

func (s *Server) RegistrationPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("registration", nil)
	}
}

func (s *Server) Registration() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user struct {
			EmailID string `form:"email_id"`
			Pass    string `form:"pass"`
		}
		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil
		}

		if _, err := db.DB.Exec("INSERT INTO user(id, pass, is_cert) VALUES(?,?,1)", user.EmailID, user.Pass); err != nil {
			fmt.Println("회원가입 실패")
			fmt.Println(err)
			// 실패를 알리는 메시지
		}
		c.Location("/login")
		return c.SendStatus(201)
	}
}
