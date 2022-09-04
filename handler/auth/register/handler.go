package register

import (
	"errors"
	"fmt"

	"github.com/4-ak/sooot/db/queries"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

type account struct {
	ID   string `form:"email_id"`
	Pass string `form:"pass"`
}

func (h *Handler) RegistrationPage(c *fiber.Ctx) error {
	return c.Render("registration", nil)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var user account

	mail := c.Cookies("mail", "")
	if mail == "" {
		return h.failure(c, fmt.Errorf("mail cookie is empty"), 0)
	}

	if err := c.BodyParser(&user); err != nil {
		return h.failure(c, err, 2)
	}
	user.ID = string(mail)

	if user.ID == "" || user.Pass == "" {
		return h.failure(c, errors.New("empty form"), 2)
	}

	hashed, err := security.CreatePass(user.ID, user.Pass)
	if err != nil {
		return h.failure(c, err, 3)
	}

	if _, err := queries.RegisterAccount().Exec(user.ID, hashed); err != nil {
		return h.failure(c, err, 3)
	}

	fmt.Printf("[정보] 계정 생성 성공 : %v\n", user.ID)
	c.Status(fiber.StatusCreated)
	return c.Render("redirect_alert", fiber.Map{
		"Msg":      "가입이 완료되었습니다.",
		"Location": "/login",
	})
}

func (h *Handler) failure(c *fiber.Ctx, err error, code int) error {
	fmt.Printf("[error] register(%v):%v", code, err)
	c.Status(fiber.StatusConflict)
	return c.Render("redirect_alert", fiber.Map{
		"Msg":      "이미 가입되어 있거나 잘못된 요청입니다.",
		"Location": "/registration",
	})
}
