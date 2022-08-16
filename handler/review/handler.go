package review

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

type review struct {
	Uid              int
	Lecture_id       int
	Writer           int
	Beneficial_point int //1~5
	Honey_point      int //1~5
	Assignment       int //1~3
	Team_project     int //1~3
	Pressentation    int //1~3
	Comment          string
}

var rev review

func (h *Handler) Review(c *fiber.Ctx) error {
	lectid := (c.Params("id"))
	userid, ok := c.Locals("uid").(string)
	if !ok {
		userid = "-1"
		fmt.Print("userid error")
	}
	result := h.SelectData(lectid)
	return c.Render("review", fiber.Map{
		"ReviewData": result,
		"Lectid":     lectid,
		"Userid":     userid,
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return c.Render("editreview", fiber.Map{
		"isUpdate": false,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("uid"))
	return c.Render("editreview", fiber.Map{
		"ReviewData": h.RowData(uid),
		"isUpdate":   true,
	})
}

func (h *Handler) RowData(uid int) review {
	row := db.Review().QueryRow(uid)
	row.Scan(
		&rev.Beneficial_point,
		&rev.Honey_point,
		&rev.Assignment,
		&rev.Team_project,
		&rev.Pressentation,
		&rev.Comment)
	return rev
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	lect_id := (c.Params("id"))
	c.BodyParser(&rev)
	_, err := db.InsertReview().Exec(
		lect_id,
		c.Locals("uid").(string),
		rev.Beneficial_point,
		rev.Honey_point,
		rev.Assignment,
		rev.Team_project,
		rev.Pressentation,
		rev.Comment)
	if err != nil {
		return c.SendString(err.Error())
	}
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) SelectData(lectid string) []review {
	row, err := db.ReviewAll().Query(lectid)
	arr := make([]review, 0)
	for row.Next() {
		row.Scan(
			&rev.Uid,
			&rev.Writer,
			&rev.Lecture_id,
			&rev.Beneficial_point,
			&rev.Honey_point,
			&rev.Assignment,
			&rev.Team_project,
			&rev.Pressentation,
			&rev.Comment)
		arr = append(arr, rev)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
	}
	return arr
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	c.BodyParser(&rev)
	_, err := db.UpdateReview().Exec(
		rev.Beneficial_point,
		rev.Honey_point,
		&rev.Assignment,
		&rev.Team_project,
		&rev.Pressentation,
		&rev.Comment,
		uid)
	if err != nil {
		fmt.Print(err)
		return c.Format(fmt.Sprintf(`
		<head>
			<meta charset="UTF-8">
			<script>
				if(!alert("값을 다시 입력해주세요")) {
					window.location="/review/%v/%v/u";
				}
			</script>
		</head>
		`, lect_id, uid))
	}
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	_, err := db.DeleteReview().Exec(uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/review/" + lect_id)
}
