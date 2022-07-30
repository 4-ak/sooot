package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
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

	sertKey := GenerateCertKey()
	// if err := SendMail(mail.Mail+"@live.wsu.ac.kr", sertKey); err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	cypherMail, err := security.EncrpytionWithBase64([]byte(mail.Mail))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cypherKey, err := security.EncrpytionWithBase64([]byte(sertKey))
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

func (s *Server) AuthUser(c *fiber.Ctx) error {
	cookie := c.Cookies("token", "")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println(err)
		return c.Redirect("/login", 302)
	}

	uid, ok := token.Claims.(jwt.MapClaims)["uid"].(string)
	if !ok {
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}

	mail, ok := token.Claims.(jwt.MapClaims)["mail"].(string)
	if !ok {
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}

	err = db.DB.QueryRow(
		`SELECT uid, id FROM user WHERE uid=? AND id=?`,
		uid, mail).Scan(&uid, &mail)
	if err != nil {
		fmt.Println(err)
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}
	c.Locals("uid", uid)
	return c.Next()
}
