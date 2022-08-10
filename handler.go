package main

import (
	authtoken "github.com/4-ak/sooot/handler/auth"
	"github.com/4-ak/sooot/handler/auth/login"
	"github.com/4-ak/sooot/handler/auth/mailcert"
	"github.com/4-ak/sooot/handler/auth/register"
	"github.com/4-ak/sooot/handler/lecture"
	"github.com/4-ak/sooot/handler/post"
	"github.com/4-ak/sooot/handler/review"
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

	lecture := server.App.Group("/lecture", authtoken.AuthUser)
	lecture.Get("/", lectureHandler.Lecture)
	lecture.Get("/c", lectureHandler.Create)
	lecture.Post("/c", lectureHandler.InsertData)
	lecture.Post("/u/:id", lectureHandler.UpdateData)
	lecture.Get("/u/:id", lectureHandler.Update)
	lecture.Get("/d/:id", lectureHandler.DeleteData)

	reviewHandler := review.Handler{}

	review := server.App.Group("/review", authtoken.AuthUser)
	review.Get("/:id", reviewHandler.Review)
	review.Get("/:id/c", reviewHandler.Create)
	review.Post("/:id/c", reviewHandler.InsertData)
	review.Get("/:lectid/:uid/u", reviewHandler.Update)
	review.Post("/:lectid/:uid/u", reviewHandler.UpdateData)
	review.Get("/:lectid/:uid/d", reviewHandler.DeleteData)

	postHandler := post.Handler{}

	post := server.App.Group("/post", authtoken.AuthUser)
	post.Get("/", postHandler.Post)
	post.Get("/c", postHandler.Create)
	post.Post("/c", postHandler.InsertData)
	post.Get("/u/:id", postHandler.Update)
	post.Post("/u/:id", postHandler.UpdateData)
	post.Get("/d/:id", postHandler.DeleteData)

	return &server
}

func (s *Server) IndexPage(c *fiber.Ctx) error {
	return c.Render("main", nil)
}
