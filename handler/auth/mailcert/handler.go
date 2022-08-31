package mailcert

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/4-ak/sooot/db/queries"
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
	if mail.Mail == "" {
		fmt.Println("[err] empty form")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var id, pw string
	if err := queries.AccountWithPass().
		QueryRow(mail.Mail).Scan(&id, &pw); err == nil {
		return c.SendStatus(fiber.StatusConflict)
	}

	c.Cookie(&fiber.Cookie{
		Name:  "mail",
		Value: string(mail.Mail),
	})
	c.Cookie(&fiber.Cookie{
		Name:  "key",
		Value: GenerateCertKey(),
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
	want := c.Cookies("key", "")

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
