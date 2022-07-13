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
  
	server.App.Post("/login", server.Login())
	server.App.Get("/registration", server.RegistrationPage())
	server.App.Post("/registration", server.Registration())
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


func (s *Server) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user struct {
			ID   string `form:"email_id"`
			Pass string `form:"pass"`
		}

		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil // TODO :500번대 메시지를 전송?
		}

		row := db.DB.QueryRow("SELECT uid FROM user WHERE id=? and pass=?", user.ID, user.Pass)
		uid := -1
		if err := row.Scan(&uid); err != nil {
			c.Accepts("html")
			c.Format(`
			<head>
				<meta charset="UTF-8">
				<script>
					if(!alert("존재하지 않는 계정이거나, 비밀번호가 틀렸습니다.")) {
						window.location = "/login";
					}
				</script>
			</head>
			`)
			return c.SendStatus(200)
		}

		cookie := fiber.Cookie{
			Name:  "token",
			Value: fmt.Sprintf("%v", uid),
		}
		c.Cookie(&cookie)

		return c.Redirect("/", fiber.StatusFound)
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
			ID   string `form:"email_id"`
			Pass string `form:"pass"`
		}
		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil // TODO :500번대 메시지를 전송?
		}
		c.Accepts("html")

		if _, err := db.DB.Exec(`INSERT INTO user(id, pass, is_cert) VALUES(?,?,1)`, user.ID, user.Pass); err != nil {
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

		fmt.Printf("[정보] 계정 생성 성공 : %v\n", user.ID)

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
