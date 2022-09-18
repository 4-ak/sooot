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

func (h *Handler) MailCertForCreateAccount(c *fiber.Ctx) error {
	return h.SendMail(c, false)
}

func (h *Handler) MailCertForResetPassword(c *fiber.Ctx) error {
	return h.SendMail(c, true)
}

func (h *Handler) SendMail(c *fiber.Ctx, doAccuntMustExist bool) error {
	var mail struct {
		Mail string
	}
	if err := json.Unmarshal(c.Body(), &mail); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	if mail.Mail == "" {
		fmt.Println("[err] empty form")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var uid, pw string
	err := queries.AccountWithPass().
		QueryRow(mail.Mail).Scan(&uid, &pw)

	if (doAccuntMustExist && err != nil) ||
		(!doAccuntMustExist && err == nil) {
		return c.SendStatus(fiber.StatusConflict)
	}

	if doAccuntMustExist {
		c.Cookie(&fiber.Cookie{
			Name:  "uid",
			Value: uid,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:   "mail",
		Value:  string(mail.Mail),
		MaxAge: int(60 * 15),
	})

	c.Cookie(&fiber.Cookie{
		Name:   "key",
		Value:  GenerateCertKey(),
		MaxAge: int(60 * 15),
	})

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) KeyCertForCreateAccount(c *fiber.Ctx) error {
	return h.KeyCert(c, "/registration")
}

func (h *Handler) KeyCertForResetPassword(c *fiber.Ctx) error {
	return h.KeyCert(c, "/reset-password")
}
func (h *Handler) KeyCert(c *fiber.Ctx, redirectURL string) error {
	var key struct {
		Key string
	}
	if err := json.Unmarshal(c.Body(), &key); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	want := c.Cookies("key", "")

	if string(want) == key.Key {
		return c.Redirect(redirectURL, fiber.StatusFound)
	} else {
		return c.SendStatus(fiber.StatusNotFound)
	}
}

func GenerateCertKey() string {
	key := rand.Int()%1000000 + 100000
	fmt.Println(key)
	return fmt.Sprintf("%d", key)
}
