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
	lect_id := (c.Params("lectid"))
	userid, ok := strconv.Atoi(c.Locals("uid").(string))
	if ok != nil {
		userid = -1
		fmt.Print("userid error")
	}
	return c.Render("review", fiber.Map{
		"ReviewData": review.SelectData(lect_id),
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
	lect_id := c.Params("lectid")
	c.BodyParser(&review)
	base_id, _ := strconv.Atoi(lect_id)
	review.Insert(base_id, c.Locals("uid").(string))
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
