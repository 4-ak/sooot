package login

import (
	"errors"
	"fmt"
	"time"

	"github.com/4-ak/sooot/db/queries"
	authtoken "github.com/4-ak/sooot/handler/auth"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct{}

type account struct {
	UID  string
	Mail string `form:"email_id"`
	Pass string `form:"pass"`
	Salt string
}

func (h *Handler) Page(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var stored, form account

	if err := c.BodyParser(&form); err != nil {
		return h.LoginFailed(c, err, 1)
	}
	if form.Mail == "" || form.Pass == "" {
		return h.LoginFailed(c, errors.New("empty form"), 1)
	}

	row := queries.AccountWithPass().QueryRow(c.FormValue("email_id", ""))

	if err := row.Scan(&stored.UID, &stored.Pass); err != nil {
		return h.LoginFailed(c, err, 2)
	}
	if err := security.ComparePass(
		[]byte(stored.Pass), form.Mail, form.Pass); err != nil {
		return h.LoginFailed(c, err, 3)
	}
	token := authtoken.UserToken{
		ID:         stored.UID,
		ClientID:   c.Cookies("ident", ""),
		Permission: "user",
		Expires:    time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	jwtoken, err := token.JWT()
	if err != nil {
		return h.LoginFailed(c, err, 4)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    jwtoken,
		HTTPOnly: true,
		//Secure: true,
	})

	return c.Redirect("/", fiber.StatusFound)
}

func (h *Handler) LoginFailed(c *fiber.Ctx, err error, code int) error {
	fmt.Printf("[error] login(%v):%v", code, err)
	c.Status(fiber.StatusConflict)
	return c.Render("redirect_alert", fiber.Map{
		"Msg":      "존재하지 않은 계정이거나, 입력된 정보가 틀렸습니다.",
		"Location": "/login",
	})
}

func (h *Handler) MakeToken(uid, mail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"mail": mail,
	})
	return token.SignedString([]byte("secret"))
}
