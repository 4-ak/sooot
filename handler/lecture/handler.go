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
		"LectureData": db.DB,
		"isUpdate":    false,
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
		"LectureData": h.RowsData(uid),
		"isUpdate":    true,
	})
}

func (h *Handler) SelectData() []lecture {
	row, err := db.DB.Query("SELECT * from lecture")
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
	_, err := db.DB.Exec(`
	INSERT INTO lecture(name, professor_name, season, grade, credit, category) 
	VALUES($1, $2, $3, $4, $5, $6)`,
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
	_, err := db.DB.Exec(`
	UPDATE lecture 
	SET name = $1, professor_name = $2, season = $3, grade = $4, credit = $5, category = $6  
	WHERE uid = $7`,
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

func (h *Handler) RowsData(uid int) lecture {
	rows := db.DB.QueryRow("SELECT * FROM lecture WHERE uid = $1", uid)
	var lect lecture
	rows.Scan(
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
	_, err := db.DB.Exec("DELETE FROM lecture WHERE uid = $1", uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/lecture")
}
