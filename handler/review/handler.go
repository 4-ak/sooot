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
	Beneficial_point int  //1~5
	Honey_point      int  //1~5
	Professor_point  int  //1~5
	Is_team          bool //0, 1
	Is_presentation  bool //0, 1
	User_id          string
}

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
		"ReviewData": db.DB,
		"isUpdate":   false,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	return c.Render("editreview", fiber.Map{
		"ReviewData": h.RowsData(uid),
		"isUpdate":   true,
	})
}

func (h *Handler) RowsData(uid int) review {
	rows := db.DB.QueryRow("SELECT * FROM review WHERE lecture_id = ?", uid)
	var rev review
	rows.Scan(
		&rev.Uid,
		&rev.Beneficial_point,
		&rev.Honey_point,
		&rev.Professor_point,
		&rev.Is_team,
		&rev.Is_presentation)
	return rev
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	var rev review
	lect_id := (c.Params("id"))
	c.BodyParser(&rev)
	_, err := db.DB.Exec(`
	INSERT INTO review(lecture_id, beneficial_point, honey_point, professor_point, is_team, is_presentation, user_id) 
		VALUES(?, ?, ?, ?, ?, ?, ?)`,
		lect_id,
		rev.Beneficial_point,
		rev.Honey_point,
		rev.Professor_point,
		rev.Is_team,
		rev.Is_presentation,
		c.Locals("uid").(string))
	if err != nil {
		return c.SendString(err.Error())
	}
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) SelectData(lectid string) []review {
	row, err := db.DB.Query("SELECT * from review WHERE lecture_id = ?", lectid)
	arr := make([]review, 0)
	for row.Next() {
		var rev review
		row.Scan(
			&rev.Uid,
			&rev.Lecture_id,
			&rev.Beneficial_point,
			&rev.Honey_point,
			&rev.Professor_point,
			&rev.Is_team,
			&rev.Is_presentation,
			&rev.User_id)
		arr = append(arr, rev)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
		return nil
	}
	return arr
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	var rev review
	c.BodyParser(&rev)
	_, err := db.DB.Exec(`
	UPDATE review
	SET beneficial_point = ?, honey_point = ?, professor_point = ?, is_team = ?, is_presentation = ? 
	WHERE uid = ?`,
		rev.Beneficial_point,
		rev.Honey_point,
		rev.Professor_point,
		rev.Is_team,
		rev.Is_presentation,
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
		// return c.Redirect("review/" + lect_id + "/" + uid + "/u")
	}
	return c.Redirect("/review/" + lect_id)
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	uid := c.Params("uid")
	lect_id := c.Params("lectid")
	_, err := db.DB.Exec("DELETE FROM review WHERE uid = ?", uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/review/" + lect_id)
}
