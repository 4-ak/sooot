package login

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db/queries"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ForgotPasswordPage(c *fiber.Ctx) error {
	return c.Render("forgot_pass", nil)
}

func (h *Handler) ResetPasswordPage(c *fiber.Ctx) error {
	uid := c.Cookies("uid", "")
	if uid == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.Render("reset_pass", fiber.Map{
		"Uid": uid,
	})
}

// TODO : Authentication!!!!! anyone can reset password!!!
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid", ""))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var plain struct {
		Pass string
	}
	if err := json.Unmarshal(c.Body(), &plain); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mail := c.Cookies("mail", "")
	if mail == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	pass, err := security.CreatePass(mail, plain.Pass)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusConflict)
	}
	_, err = queries.SetPasswordOfAccount().Exec(uid, string(pass))
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.SendStatus(fiber.StatusOK)
}
