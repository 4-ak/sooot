package review

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db/model"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func (h *Handler) Review(c *fiber.Ctx) error {
	review := model.NewReview()
	lect_id := (c.Params("id"))
	userid, ok := c.Locals("uid").(string)
	if !ok {
		userid = "-1"
		fmt.Print("userid error")
	}
	review.SelectData(lect_id)
	return c.Render("review", fiber.Map{
		"ReviewData": review,
		"Lectid":     lect_id,
		"Userid":     userid,
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return c.Render("editreview", fiber.Map{
		"isUpdate": false,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	review := model.NewReview()
	uid, _ := strconv.Atoi(c.Params("uid"))
	review.RowData(uid)
	return c.Render("editreview", fiber.Map{
		"ReviewData": review,
		"isUpdate":   true,
	})
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	review := model.NewReview()
	lect_id := (c.Params("id"))
	c.BodyParser(&review)
	review.Insert(lect_id, c.Locals("uid").(string))
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	review := model.NewReview()
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	c.BodyParser(&review)
	review.Update(uid)
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	review := model.NewReview()
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	review.Delete(uid)
	return c.Redirect("/review/" + lect_id)
}
