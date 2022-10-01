package review

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db/model"
	authtoken "github.com/4-ak/sooot/handler/auth"
	"github.com/4-ak/sooot/score"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func (h *Handler) Review(c *fiber.Ctx) error {
	review := model.NewReview()
	lecture := model.NewLecture()
	lect_id := (c.Params("lectid"))
	userid, ok := strconv.Atoi(c.Locals("user").(authtoken.UserToken).ID)
	if ok != nil {
		userid = -1
		fmt.Print("userid error")
	}
	lecture_data := lecture.RowData(lect_id)
	review_data := review.SelectData(lect_id)
	score_text := score.Score(review_data)
	return c.Render("review", fiber.Map{
		"LectureData": lecture_data,
		"ReviewData":  review_data,
		"Lectid":      lect_id,
		"Userid":      userid,
		"Scale_5":     make([]int, 5),
		"Scale_3":     make([]int, 3),
		"Score_Avg":   score_text,
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	lecture := model.NewLecture()
	lect_id := (c.Params("lectid"))
	lecture_data := lecture.RowData(lect_id)
	return c.Render("editreview", fiber.Map{
		"LectureData": lecture_data,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	review := model.NewReview()
	lecture := model.NewLecture()
	lect_id := (c.Params("lectid"))
	uid, _ := strconv.Atoi(c.Params("uid"))
	review.RowData(uid)
	lecture_data := lecture.RowData(lect_id)
	return c.Render("editreview", fiber.Map{
		"LectureData": lecture_data,
		"ReviewData":  review,
		"isUpdate":    true,
	})
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	review := model.NewReview()
	lect_id := c.Params("lectid")
	c.BodyParser(&review)
	base_id, _ := strconv.Atoi(lect_id)
	review.Insert(base_id, c.Locals("user").(authtoken.UserToken).ID)
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
