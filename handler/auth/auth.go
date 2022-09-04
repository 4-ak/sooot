package authtoken

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func IdentifyClient(c *fiber.Ctx) error {
	clientID := c.Cookies("ident", "none")
	if clientID == "none" {
		clientID = utils.UUID()
		c.Cookie(&fiber.Cookie{
			Name:     "ident",
			Value:    clientID,
			HTTPOnly: true,
			MaxAge:   60 * 60 * 24 * 31 * 12,
		})
	}

	return c.Next()
}

func AuthUser(c *fiber.Ctx) error {
	cookie := c.Cookies("token", "none")
	if cookie == "none" {
		return c.Redirect("/login", 302)
	}

	token, err := NewUserToken(cookie)
	if err != nil {
		fmt.Println(err)
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}

	cid := c.Cookies("ident", "")

	if token.ClientID != cid {
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}
	if !token.isExpiration() {
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}
	c.Locals("user", token)
	return c.Next()
}
