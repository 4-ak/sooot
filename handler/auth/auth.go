package authtoken

import (
	"fmt"

	"github.com/4-ak/sooot/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthUser(c *fiber.Ctx) error {
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

	err = db.AccountExists(uid, mail).Scan(&uid, &mail)
	if err != nil {
		fmt.Println(err)
		c.ClearCookie("token")
		return c.Redirect("/login", 302)
	}
	c.Locals("uid", uid)
	return c.Next()
}
