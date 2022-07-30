package main

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db"
	authtoken "github.com/4-ak/sooot/handler/auth"
	"github.com/4-ak/sooot/handler/auth/login"
	"github.com/4-ak/sooot/handler/auth/mailcert"
	"github.com/4-ak/sooot/handler/auth/register"
	"github.com/4-ak/sooot/handler/lecture"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Server struct {
	App *fiber.App
}

func NewServer() *Server {
	security.KeyGen()

	engine := html.New("./tmpl", ".html")
	server := Server{
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	server.App.Get("/", server.IndexPage)

	loginHandler := login.Handler{}
	registerHandler := register.Handler{}
	mailCertHandler := mailcert.Handler{}
	auth := server.App.Group("/")
	auth.Get("/login", loginHandler.Page)
	auth.Post("/login", loginHandler.Login)
	auth.Get("/mail-cert", mailCertHandler.Page)
	auth.Post("/mail-cert", mailCertHandler.SendMail)
	auth.Post("/key-cert", mailCertHandler.KeyCert)
	auth.Get("/registration", registerHandler.RegistrationPage)
	auth.Post("/registration", registerHandler.Register)

	lectureHandler := lecture.Handler{}

	course := server.App.Group("/lecture", authtoken.AuthUser)
	course.Get("/", lectureHandler.Lecture)
	course.Get("/c", lectureHandler.CreateLecture)
	course.Post("/c", lectureHandler.InsertLectureDB)
	course.Post("/u/:id", lectureHandler.UpdateLectureDB)
	course.Get("/u/:id", lectureHandler.UpdateLecture)
	course.Get("/d/:id", lectureHandler.DeleteLectureDB)

	review := server.App.Group("/review", authtoken.AuthUser)
	review.Get("/:id", server.Review)
	review.Get("/:id/c", server.CreateReview)
	review.Post("/:id/c", server.InsertReview)
	review.Get("/:lectid/:uid/u", server.UpdateReview)
	review.Post("/:lectid/:uid/u", server.UpdateReviewDB)
	review.Get("/:lectid/:uid/d", server.DeleteReviewDB)

	return &server
}

func (s *Server) IndexPage(c *fiber.Ctx) error {
	return c.Render("main", nil)
}

type review struct {
	Uid              int
	Lecture_id       int
	Beneficial_point int  //1~5
	Honey_point      int  //1~5
	Professor_point  int  //1~5
	Is_team          bool //0, 1
	Is_presentation  bool //0, 1
	User_id          string
}

func (s *Server) Review(c *fiber.Ctx) error {
	lectid := (c.Params("id"))
	userid, ok := c.Locals("uid").(string)
	if !ok {
		userid = "-1"
		fmt.Print("userid error")
	}
	result := s.SelectReviewDB(lectid)
	return c.Render("review", fiber.Map{
		"ReviewData": result,
		"Lectid":     lectid,
		"Userid":     userid,
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
		lect_id, rev.Beneficial_point, rev.Honey_point, rev.Professor_point, rev.Is_team, rev.Is_presentation, c.Locals("uid").(string))
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
