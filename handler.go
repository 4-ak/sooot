package main

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/domain"
	"github.com/4-ak/sooot/handler/auth/login"
	"github.com/4-ak/sooot/handler/auth/register"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Server struct {
	*domain.CourseList
	App *fiber.App
}

func NewServer(courses *domain.CourseList) *Server {
	security.KeyGen()

	engine := html.New("./tmpl", ".html")
	server := Server{
		CourseList: courses,
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	server.App.Get("/", server.ViewCourses())

	loginHandler := login.Handler{}
	registerHandler := register.Handler{}
	auth := server.App.Group("/")
	auth.Get("/login", loginHandler.LoginPage)
	auth.Post("/login", loginHandler.Login)
	auth.Get("/mail-cert", server.MailCertPage)
	auth.Post("/mail-cert", server.MailSend)
	auth.Post("/key-cert", server.KeyCert)
	auth.Get("/registration", registerHandler.RegistrationPage)
	auth.Post("/registration", registerHandler.Register)

	course := server.App.Group("/course", server.AuthUser)
	course.Get("/", server.Course())
	course.Get("/1", server.CreateCourse())
	course.Post("/1", server.InsertCourseDB())
	course.Post("/2/:id", server.UpdateCourseDB())
	course.Get("/2/:id", server.UpdateCourse())
	course.Get("/d/:id", server.DeleteCourseDB())

	review := server.App.Group("/review", server.AuthUser)
	review.Get("/:id", server.Review)
	review.Get("/:id/c", server.CreateReview)
	review.Post("/:id/c", server.InsertReview)
	review.Get("/:lectid/:uid/u", server.UpdateReview)
	review.Post("/:lectid/:uid/u", server.UpdateReviewDB)
	review.Get("/:lectid/:uid/d", server.DeleteReviewDB)

	return &server
}

func (s *Server) ViewCourses() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return c.Render("main", fiber.Map{
			"Courses": s.List,
		})
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
	User_id          int
}

func (s *Server) Review(c *fiber.Ctx) error {
	lectid := (c.Params("id"))
	userid, ok := c.Locals("uid").(float64)
	if !ok {
		userid = -1
		fmt.Print("userid error")
	}
	result := s.SelectReviewDB(lectid)
	return c.Render("review", fiber.Map{
		"ReviewData": result,
		"Lectid":     lectid,
		"Userid":     int(userid),
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
	INSERT INTO review(lecture_id, beneficial_point, honey_point, professor_point, is_team, is_presentation, user_id) 
		VALUES(?, ?, ?, ?, ?, ?, ?)`,
		lect_id, rev.Beneficial_point, rev.Honey_point, rev.Professor_point, rev.Is_team, rev.Is_presentation, c.Locals("uid"))
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
		row.Scan(&rev.Uid, &rev.Lecture_id, &rev.Beneficial_point, &rev.Honey_point, &rev.Professor_point, &rev.Is_team, &rev.Is_presentation, &rev.User_id)
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
