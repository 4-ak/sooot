package lecture

import (
	"strconv"

	"github.com/4-ak/sooot/db/model"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func (h *Handler) Create(c *fiber.Ctx) error {
	lect := model.NewLecture()
	major := model.NewMajor()
	semester := model.NewSemester()
	return c.Render("editlecture", fiber.Map{
		"isUpdate":     false,
		"Lecture_base": lect.Base.Lecture_base(),
		"Major":        major.Major(),
		"Semester":     semester.Semester(),
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
	c.BodyParser(&lect.Base)
	c.BodyParser(&lect.Data)
	lect.Data.Update(uid)
	return c.Redirect("/lecture")
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	c.BodyParser(&lect.Base)
	c.BodyParser(&lect.Data)
	lect.Data.Insert(lect.Base.Name, lect.Base.Professor)
	return c.Redirect("/lecture")
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	uid, _ := strconv.Atoi(c.Params("id"))
	lect.Data.Delete(uid)
	return c.Redirect("/lecture")
}
