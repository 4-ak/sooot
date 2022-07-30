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
	Season         string
	Grade          string
	Credit         string
	Category       string
}

func (h *Handler) CreateCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("editcourse", fiber.Map{
			"CourseData": db.DB,
			"isUpdate":   false,
		})
	}
}

func (h *Handler) Course() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("course", fiber.Map{
			"CourseData": h.SelectCourseDB(),
		})
	}
}

func (h *Handler) UpdateCourse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		return c.Render("editcourse", fiber.Map{
			"CourseData": h.SendCourseDB(uid),
			"isUpdate":   true,
		})
	}
}

func (h *Handler) SelectCourseDB() []lecture {

	row, err := db.DB.Query("SELECT * from lecture")
	arr := make([]lecture, 0)
	for row.Next() {
		var lect lecture
		row.Scan(&lect.Uid, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
		arr = append(arr, lect)
	}
	if err != nil {
		fmt.Println(err)
	}
	return arr
}

func (h *Handler) InsertCourseDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`
		INSERT INTO lecture(name, professor_name, season, grade, credit, category) 
		VALUES(?, ?, ?, ?, ?, ?)`,
			lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category)
		if err != nil {
			return c.SendString("INSERT ERROR")
		}
		return c.Redirect("/course")
	}
}

func (h *Handler) UpdateCourseDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		var lect lecture
		c.BodyParser(&lect)
		_, err := db.DB.Exec(`
		UPDATE lecture 
		SET name = ?, professor_name = ?, season = ?, grade = ?, credit = ?, category = ?  
		WHERE uid = ?`,
			lect.Name, lect.Professor_name, lect.Season, lect.Grade, lect.Credit, lect.Category,
			uid)
		if err != nil {
			fmt.Print(err)
			return c.SendString("UPDATE ERROR")
		}
		return c.Redirect("/course")
	}
}

func (h *Handler) SendCourseDB(uid int) lecture {
	rows := db.DB.QueryRow("SELECT * FROM lecture WHERE uid = ?", uid)
	var lect lecture
	rows.Scan(&lect.Uid, &lect.Name, &lect.Professor_name, &lect.Season, &lect.Grade, &lect.Credit, &lect.Category)
	return lect
}

func (h *Handler) DeleteCourseDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := strconv.Atoi(c.Params("id"))
		_, err := db.DB.Exec("DELETE FROM lecture WHERE uid = ?", uid)
		if err != nil {
			fmt.Print(err)
			return c.SendString("DELETE ERROR")
		}
		return c.Redirect("/course")
	}
}
