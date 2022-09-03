package authtoken

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v4"
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
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		fmt.Println(err)
		c.ClearCookie("token")
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

	err = queries.AccountExists().QueryRow(uid, mail).Scan(&uid, &mail)
	if err != nil {
		fmt.Println(err)
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}
	c.Locals("uid", uid)
	return c.Next()
}
