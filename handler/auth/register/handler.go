package register

import (
	"fmt"

	"github.com/4-ak/sooot/db"
	"github.com/4-ak/sooot/security"
	"github.com/gofiber/fiber/v2"
)

const (
	registerFailCode = `
	<head>
		<meta charset="UTF-8">
		<script>
		if(!alert("이미 가입되어 있거나 요청이 잘못되었습니다.")) {
			window.location = "/registration";
		}
		</script>
	</head>
	`
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

	mail, err := security.DecrptionWithBase64(c.Cookies("mail", ""))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := c.BodyParser(&user); err != nil {
		return h.failure(c, err, 2)
	}

	user.ID = string(mail)
	c.Accepts("html")

	hashed, err := security.CreatePass(user.ID, user.Pass, "123")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if _, err := db.RegisterAccount().Exec(user.ID, hashed); err != nil {
		return h.failure(c, err, 3)
	}

	fmt.Printf("[정보] 계정 생성 성공 : %v\n", user.ID)
	c.Format(`
		<head>
			<meta charset="UTF-8">
			<script>
				if(!alert("가입이 완료되었습니다!")) {
					window.location="/login";
				}
			</script>
		</head>
		`)
	return c.SendStatus(201)
}

func (h *Handler) failure(c *fiber.Ctx, err error, code int) error {
	fmt.Printf("[error] register(%v):%v", code, err)
	c.Accepts("html")
	c.Format(registerFailCode)
	return c.SendStatus(fiber.StatusUnauthorized)
}
