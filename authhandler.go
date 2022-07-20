package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
)

const (
	loginFailCode = `
	<head>
		<meta charset="UTF-8">
		<script>
			if(!alert("존재하지 않는 계정이거나, 비밀번호가 틀렸습니다.")) {
				window.location = "/login";
			}
		</script>
	</head>
	`
	registerFailCode = `
	<head>
		<meta charset="UTF-8">
		<script>
		if(!alert("이미 가입되어 있습니다.")) {
			window.location = "/registration";
		}
		</script>
	</head>
	`
)

type Account struct {
	ID   string `form:"email_id"`
	Pass string `form:"pass"`
	Salt string
}

func (s *Server) LoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	}
}

func (s *Server) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user Account

		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil // TODO :500번대 메시지를 전송?
		}

		var uid int = -1
		var hashed string

		row := db.DB.QueryRow("SELECT uid, pass FROM user WHERE id=?;", user.ID)
		if err := row.Scan(&uid, &hashed); err != nil {
			fmt.Println(err)
			c.Accepts("html")
			c.Format(loginFailCode)
			return c.SendStatus(200)
		}

		//pass vaildation
		if err := security.ComparePass(user.ID, user.Pass, "123", []byte(hashed)); err != nil {
			fmt.Println(err)
			c.Accepts("html")
			c.Format(loginFailCode)
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
		var user Account

		mail, err := security.DecrptionWithBase64(c.Cookies("mail", ""))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil // TODO :500번대 메시지를 전송?
		}

		user.ID = string(mail)
		c.Accepts("html")

		hashed, err := security.CreatePass(user.ID, user.Pass, "123")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if _, err := db.DB.Exec(`INSERT INTO user(id, pass, is_cert) VALUES(?,?,1)`, user.ID, hashed); err != nil {
			fmt.Printf("[정보] 계정 생성 실패 : %v\n", err.Error())
			fmt.Println(err)

			c.Format(registerFailCode)
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

func (s *Server) MailCertPage(c *fiber.Ctx) error {
	return c.Render("mail_cert", nil)
}

func (s *Server) MailSend(c *fiber.Ctx) error {
	var mail struct {
		Mail string
	}
	if err := json.Unmarshal(c.Body(), &mail); err != nil {
		fmt.Println(err)
		return c.SendStatus(404)
	}
	cypherMail, err := security.EncrpytionWithBase64([]byte(mail.Mail))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cypherKey, err := security.EncrpytionWithBase64([]byte(GenerateCertKey()))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	c.Cookie(&fiber.Cookie{
		Name:  "mail",
		Value: string(cypherMail),
	})
	c.Cookie(&fiber.Cookie{
		Name:  "key",
		Value: cypherKey,
	})

	return c.SendString("전송된 메일 :" + mail.Mail)
}

func GenerateCertKey() string {
	key := rand.Int()%1000000 + 100000
	fmt.Println(key)
	return fmt.Sprintf("%d", key)
}

func (s *Server) KeyCert(c *fiber.Ctx) error {
	var key struct {
		Key string
	}
	if err := json.Unmarshal(c.Body(), &key); err != nil {
		fmt.Println(err)
		return c.SendStatus(404)
	}
	keyChyper := c.Cookies("key", "")
	want, err := security.DecrptionWithBase64(keyChyper)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if string(want) == key.Key {
		return c.Redirect("/registration", 302)
	} else {
		return c.SendStatus(404)
	}
}
