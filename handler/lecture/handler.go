package lecture

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/4-ak/sooot/db/model"
	"github.com/4-ak/sooot/handler/lectbrowser"
	"github.com/gofiber/fiber/v2"
)

var browser lectbrowser.LectureBrowser

type Handler struct{}

func (h *Handler) Create(c *fiber.Ctx) error {
	major := model.NewMajor()
	semester := model.NewSemester()
	return c.Render("addreview", fiber.Map{
		"Major":    major.Major(),
		"Semester": semester.Semester(),
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
	review := model.NewReview()
	c.BodyParser(&lect.Base)
	c.BodyParser(&lect.Data)
	c.BodyParser(&review)
	lect.CompareData()
	if lect.Data.Uid == 0 {
		fmt.Println(lect.Data)
		lect.Data.Insert(lect.Base.Name, lect.Base.Professor)
		review.Insert(lect.Data.Uid, c.Locals("uid").(string))
	} else {
		fmt.Println(lect.Data)
		review.Insert(lect.Data.Uid, c.Locals("uid").(string))
	}
	return c.Redirect("/lecture")
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	lect := model.NewLecture()
	uid, _ := strconv.Atoi(c.Params("id"))
	lect.Data.Delete(uid)
	return c.Redirect("/lecture")
}

func (h *Handler) GetLectures(c *fiber.Ctx) error {
	query := c.Query("query", "")
	length := utf8.RuneCountInString(query)

	switch length {
	case 0:
		fmt.Println("잘못된 값이 들어왔음")
		return c.SendStatus(fiber.StatusBadRequest)
	case 1:
		return c.JSON(browser.FindByFirstInitial([]rune(query)[0]))
	default:
		return c.JSON(browser.FindPriorityByConsonants(query))
	}
}

func (h *Handler) GetMajor(c *fiber.Ctx) error {
	major := model.NewMajor()
	query := c.Query("query", "")
	if query == "" {
		fmt.Println("잘못된 값이 들어왔음")
		return c.SendStatus(fiber.StatusNotFound)
	}
	majors := major.Major()
	data, err := json.Marshal(majors)
	if err != nil {
		fmt.Println("json 변환 실패")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString(string(data))
}

func (h *Handler) CacheLecture() {
	lect := model.NewLecture()
	lectures := lect.Base.Lecture_base()
	browser.Init(lectures)
}
