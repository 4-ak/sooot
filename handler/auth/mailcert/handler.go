package mailcert

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func (h *Handler) Page(c *fiber.Ctx) error {
	return c.Render("mail_cert", nil)
}

func (h *Handler) SendMail(c *fiber.Ctx) error {
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

func (h *Handler) KeyCert(c *fiber.Ctx) error {
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

func GenerateCertKey() string {
	key := rand.Int()%1000000 + 100000
	fmt.Println(key)
	return fmt.Sprintf("%d", key)
}
