package login

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ForgotPasswordPage(c *fiber.Ctx) error {
	return c.Render("forgot_pass", nil)
}

func (h *Handler) ResetPasswordPage(c *fiber.Ctx) error {
	return c.Render("reset_pass", nil)
}
