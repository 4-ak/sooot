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

	server.App.Get("/course", server.Course())
	server.App.Get("/course/1", server.CreateCourse())
	server.App.Post("/course/1", server.InsertCourseDB())
	server.App.Post("/course/2/:id", server.UpdateCourseDB())
	server.App.Get("/course/2/:id", server.UpdateCourse())
	server.App.Get("/course/d/:id", server.DeleteCourseDB())

	server.App.Get("/review/:id", server.Review)
	server.App.Get("/review/:id/c", server.CreateReview)
	server.App.Post("/review/:id/c", server.InsertReview)
	server.App.Get("/review/:lectid/:uid/u", server.UpdateReview)
	server.App.Post("/review/:lectid/:uid/u", server.UpdateReviewDB)
	server.App.Get("/review/:lectid/:uid/d", server.DeleteReviewDB)

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
		return c.Render("editcourse", fiber.Map{
			"CourseData": db.DB,
			"isUpdate":   false,
		})
	}
}

func (s *Server) Course() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("course", fiber.Map{
			"CourseData": s.SelectCourseDB(),
		})
	}
}

func (s *Server) UpdateCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		return c.Render("editcourse", fiber.Map{
			"CourseData": s.SendCourseDB(uid),
			"isUpdate":   true,
		})
	}
}

func (s *Server) SelectCourseDB() []lecture {

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

func (s *Server) InsertCourseDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`
		INSERT INTO lecture(name, professor_name, season, grade, credit, category) 
		VALUES(?, ?, ?, ?, ?, ?)`,
			lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category)
		if err != nil {
			return c.SendString("INSERT ERROR")
		}
		return c.Redirect("/course")
	}
}

func (s *Server) UpdateCourseDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`
		UPDATE lecture 
		SET name = ?, professor_name = ?, season = ?, grade = ?, credit = ?, category = ?  
		WHERE uid = ?`,
			lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category,
			uid)
		if err != nil {
			fmt.Print(err)
			return c.SendString("UPDATE ERROR")
		}
		return c.Redirect("/course")
	}
}

func (s *Server) SendCourseDB(uid int) lecture {
	rows := db.DB.QueryRow("SELECT * FROM lecture WHERE uid = ?", uid)
	var lect lecture
	rows.Scan(&lect.Uid, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
	return lect
}

func (s *Server) DeleteCourseDB() fiber.Handler {
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

type review struct {
	Uid              int
	Lecture_id       int
	Beneficial_point int  //1~5
	Honey_point      int  //1~5
	Professor_point  int  //1~5
	Is_team          bool //0, 1
	Is_presentation  bool //0, 1
}

func (s *Server) Review(c *fiber.Ctx) error {
	lectid := (c.Params("id"))
	result := s.SelectReviewDB(lectid)
	return c.Render("review", fiber.Map{
		"ReviewData": result,
		"Lectid":     lectid,
	})
}

func (s *Server) CreateReview(c *fiber.Ctx) error {
	return c.Render("editreview", fiber.Map{
		"ReviewData": db.DB,
		"isUpdate":   false,
	})
}

func (s *Server) UpdateReview(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	return c.Render("editreview", fiber.Map{
		"ReviewData": s.SendReviewDB(uid),
		"isUpdate":   true,
	})
}

func (s *Server) SendReviewDB(uid int) review {
	rows := db.DB.QueryRow("SELECT * FROM review WHERE lecture_id = ?", uid)
	var rev review
	rows.Scan(&rev.Uid, &rev.Beneficial_point, &rev.Honey_point, &rev.Professor_point, &rev.Is_team, &rev.Is_presentation)
	return rev
}

func (s *Server) InsertReview(c *fiber.Ctx) error {
	var rev review
	lect_id := (c.Params("id"))
	c.BodyParser(&rev)
	_, err := db.DB.Exec(`
	INSERT INTO review(lecture_id, beneficial_point, honey_point, professor_point, is_team, is_presentation) 
		VALUES(?, ?, ?, ?, ?, ?)`,
		lect_id, rev.Beneficial_point, rev.Honey_point, rev.Professor_point, rev.Is_team, rev.Is_presentation)
	if err != nil {
		return c.SendString(err.Error())
	}
	return c.Redirect("/review/" + lect_id)
}

func (s *Server) SelectReviewDB(lectid string) []review {
	row, err := db.DB.Query("SELECT * from review WHERE lecture_id = ?", lectid)
	arr := make([]review, 0)
	for row.Next() {
		var rev review
		row.Scan(&rev.Uid, &rev.Lecture_id, &rev.Beneficial_point, &rev.Honey_point, &rev.Professor_point, &rev.Is_team, &rev.Is_presentation)
		arr = append(arr, rev)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
		return nil
	}
	return arr
}

func (s *Server) UpdateReviewDB(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	var rev review
	c.BodyParser(&rev)
	_, err := db.DB.Exec(`
	UPDATE review
	SET beneficial_point = ?, honey_point = ?, professor_point = ?, is_team = ?, is_presentation = ? 
	WHERE uid = ?`,
		rev.Beneficial_point, rev.Honey_point, rev.Professor_point, rev.Is_team, rev.Is_presentation,
		uid)
	if err != nil {
		fmt.Print(err)
		return c.Format(fmt.Sprintf(`
		<head>
			<meta charset="UTF-8">
			<script>
				if(!alert("값을 다시 입력해주세요")) {
					window.location="/review/%v/%v/u";
				}
			</script>
		</head>
		`, lect_id, uid))
		// return c.Redirect("review/" + lect_id + "/" + uid + "/u")
	}
	return c.Redirect("/review/" + lect_id)
}

func (s *Server) DeleteReviewDB(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	_, err := db.DB.Exec("DELETE FROM review WHERE uid = ?", uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/review/" + lect_id)
}
