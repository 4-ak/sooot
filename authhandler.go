package main

import (
	"fmt"

	"github.com/4-ak/sooot/db"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
func createPlainPass(id, pass, salt string) []byte {
	data := append([]byte(id), []byte(pass)...)
	return append(data, []byte(salt)...)
}

func createPass(id, pass, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(createPlainPass(id, pass, salt), 11)
}

func comparePass(id, pass, salt string, hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, createPlainPass(id, pass, salt))
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
		if err := comparePass(user.ID, user.Pass, "123", []byte(hashed)); err != nil {
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

		if err := c.BodyParser(&user); err != nil {
			fmt.Println(err)
			return nil // TODO :500번대 메시지를 전송?
		}
		c.Accepts("html")

		hashed, err := createPass(user.ID, user.Pass, "123")
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

		c.Format(registerFailCode)
		return c.SendStatus(201)
	}
}
