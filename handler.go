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
		c.Accepts("html")

		if _, err := db.DB.Exec("INSERT INTO user(id, pass, is_cert) VALUES(?,?,1)", user.EmailID, user.Pass); err != nil {
			fmt.Printf("[정보] 계정 생성 실패 : %v\n", err.Error())
			fmt.Println(err)

			c.Format(`
			<head>
				<meta charset="UTF-8">
				<script>
				if(!alert("이미 가입되어 있습니다.")) {
					//window.location = "/registration";
				}
				</script>
			</head>
			`)
			return c.SendStatus(200)
		}

		fmt.Printf("[정보] 계정 생성 성공 : %v\n", user.EmailID)

		c.Format(`
		<head>
			<meta charset="UTF-8">
			<script>
				if(!alert("가입이 완료되었습니다!")) {
					window.location="/login";
				}
			</script>
		</head>
		`)
		return c.SendStatus(201)
	}
}
