package lecture

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

type lecture struct {
	Uid            int
	Name           string
	Professor_name string
	Season         int
	Grade          int
	Credit         int
	Category       int
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return c.Render("editlecture", fiber.Map{
		"isUpdate": false,
	})
}

func (h *Handler) Lecture(c *fiber.Ctx) error {
	return c.Render("lecture", fiber.Map{
		"LectureData": h.SelectData(),
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	return c.Render("editlecture", fiber.Map{
		"LectureData": h.RowData(uid),
		"isUpdate":    true,
	})
}

func (h *Handler) SelectData() []lecture {
	row, err := db.TakeAllLecture()
	arr := make([]lecture, 0)
	for row.Next() {
		var lect lecture
		row.Scan(
			&lect.Uid,
			&lect.Name,
			&lect.Professor_name,
			&lect.Season,
			&lect.Grade,
			&lect.Credit,
			&lect.Category)
		arr = append(arr, lect)
	}
	if err != nil {
		fmt.Println(err)
	}
	return arr
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	var lect lecture
	c.BodyParser(&lect)
	err := db.InsertLecture(
		lect.Name,
		lect.Professor_name,
		lect.Season,
		lect.Grade,
		lect.Credit,
		lect.Category)
	if err != nil {
		return c.SendString("INSERT ERROR")
	}
	return c.Redirect("/lecture")
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	var lect lecture
	c.BodyParser(&lect)
	err := db.UpdateLecture(
		lect.Name,
		lect.Professor_name,
		lect.Season,
		lect.Grade,
		lect.Credit,
		lect.Category,
		uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("UPDATE ERROR")
	}
	return c.Redirect("/lecture")
}

func (h *Handler) RowData(uid int) lecture {
	row := db.RowLecture(uid)
	var lect lecture
	row.Scan(
		&lect.Uid,
		&lect.Name,
		&lect.Professor_name,
		&lect.Season,
		&lect.Grade,
		&lect.Credit,
		&lect.Category)
	return lect
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	err := db.DeleteLecture(uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/lecture")
}
