package post

import (
	"fmt"
	"strconv"

	"github.com/4-ak/sooot/db"
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

type post struct {
	Uid       int
	Title     string
	Contents  string
	Recommand int
	Writer    int
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return c.Render("editpost", fiber.Map{
		"isUpdate": false,
	})
}

func (h *Handler) Post(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok {
		uid = "-1"
		fmt.Print("userid error")
	}
	userid, _ := strconv.Atoi(uid)
	return c.Render("post", fiber.Map{
		"PostData": h.SelectData(),
		"Userid":   userid,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	return c.Render("editpost", fiber.Map{
		"isUpdate": true,
	})
}

func (h *Handler) SelectData() []post {
	row, err := db.PostAll()
	arr := make([]post, 0)
	for row.Next() {
		var post post
		row.Scan(
			&post.Uid,
			&post.Title,
			&post.Contents,
			&post.Recommand,
			&post.Writer)
		arr = append(arr, post)
	}
	if err != nil {
		fmt.Println(err)
	}
	return arr
}

func (h *Handler) InsertData(c *fiber.Ctx) error {
	var post post
	c.BodyParser(&post)
	err := db.InsertPost(
		post.Title,
		post.Contents,
		c.Locals("uid").(string),
		post.Recommand)
	if err != nil {
		fmt.Println(err)
		return c.SendString("INSERT ERROR")
	}
	return c.Redirect("/post")
}

func (h *Handler) UpdateData(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	var post post
	c.BodyParser(&post)
	err := db.UpdatePost(
		post.Title,
		post.Contents,
		uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("UPDATE ERROR")
	}
	return c.Redirect("/post")
}

func (h *Handler) DeleteData(c *fiber.Ctx) error {
	uid, _ := strconv.Atoi(c.Params("id"))
	err := db.DeletePost(uid)
	if err != nil {
		fmt.Print(err)
		return c.SendString("DELETE ERROR")
	}
	return c.Redirect("/post")
}
