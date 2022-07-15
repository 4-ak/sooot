package main

import (
	"fmt"
	"strconv"

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

	server.App.Get("/course", server.EditCourse())
	server.App.Get("/course/1", server.CreateCourse())
	server.App.Post("/course/1", server.InsertDB())
	server.App.Post("/course/2/:id", server.UpdateDB())
	server.App.Get("/course/2/:id", server.UpdateCourse())
	server.App.Get("/course/d/:id", server.DeleteDB())

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
					window.location = "/registration";
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
}

type lecture struct {
	Uid            int
	Name           string
	Professor_name string
	Season         string
	Grade          string
	Credit         string
	Category       string
}

func (s *Server) CreateCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("createcourse", fiber.Map{
			"DB": db.DB,
		})
	}
}

func (s *Server) EditCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("editcourse", fiber.Map{
			"DB": s.SelectDB(),
		})
	}
}

func (s *Server) UpdateCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		return c.Render("updatecourse", fiber.Map{
			"DB": s.SendDB(uid),
		})
	}
}

func (s *Server) SelectDB() []lecture {

	row, err := db.DB.Query("SELECT * from lecture")
	arr := make([]lecture, 0)
	for row.Next() {
		var lect lecture
		row.Scan(&lect.Uid, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
		arr = append(arr, lect)
	}
	if err != nil {
		fmt.Println(err)
	}
	return arr
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

func (s *Server) UpdateDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`UPDATE lecture 
		SET name = ?, professor_name = ?, season = ?, grade = ?, credit = ?, category = ?  WHERE uid = ?`, lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category, uid)
		if err != nil {
			fmt.Print(err)
			return c.SendString("UPDATE ERROR")
		}
		return c.Redirect("/course")
	}
}

func (s *Server) SendDB(uid int) lecture {
	rows := db.DB.QueryRow("SELECT * FROM lecture WHERE uid = ?", uid)
	var lect lecture
	rows.Scan(&lect.Uid, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
	return lect

}

func (s *Server) DeleteDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		_, err := db.DB.Exec("DELETE FROM lecture WHERE uid = ?", uid)
		if err != nil {
			fmt.Print(err)
			return c.SendString("DELETE ERROR")
		}
		return c.Redirect("/course")
	}
}

func (s *Server) CreateReview() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("test1")
	}
}
