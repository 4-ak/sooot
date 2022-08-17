package lecture

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db/model"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func (h *Handler) Create(c *fiber.Ctx) error {
	return c.Render("editlecture", fiber.Map{
		"isUpdate": false,
	})
}

func (h *Handler) Lecture(c *fiber.Ctx) error {
	lect := model.NewLecture()
	return c.Render("lecture", fiber.Map{
		"LectureData": lect.SelectData(),
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	lect := model.NewLecture()
	uid, _ := strconv.Atoi(c.Params("id"))
	lect.RowData(uid)
	return c.Render("editlecture", fiber.Map{
		"LectureData": lect,
		"isUpdate":    true,
	})
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	uid, _ := strconv.Atoi(c.Params("id"))
	c.BodyParser(&lect)
	fmt.Println(lect)
	lect.Update(uid)
	return c.Redirect("/lecture")
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	c.BodyParser(&lect)
	lect.Insert()
	return c.Redirect("/lecture")
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	uid, _ := strconv.Atoi(c.Params("id"))
	lect.Delete(uid)
	return c.Redirect("/lecture")
}
