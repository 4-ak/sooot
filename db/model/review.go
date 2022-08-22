package model

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
)

type Review struct {
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

func NewReview() Review {
	review := Review{}
	return review
}

func (r *Review) SelectData(lect_id string) []Review {
	row, err := queries.ReviewAll().Query(lect_id)
	arr := make([]Review, 0)
	for row.Next() {
		row.Scan(
			&r.Uid,
			&r.Writer,
			&r.Lecture_id,
			&r.Beneficial_point,
			&r.Honey_point,
			&r.Assignment,
			&r.Team_project,
			&r.Pressentation,
			&r.Comment)
		arr = append(arr, *r)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
	}
	return arr
}

func (r *Review) Insert(lect_id, uid string) {
	_, err := queries.InsertReview().Exec(
		lect_id,
		uid,
		r.Beneficial_point,
		r.Honey_point,
		r.Assignment,
		r.Team_project,
		r.Pressentation,
		r.Comment)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Review) Update(uid string) {
	_, err := queries.UpdateReview().Exec(
		r.Beneficial_point,
		r.Honey_point,
		r.Assignment,
		r.Team_project,
		r.Pressentation,
		r.Comment,
		uid)
	if err != nil {
		fmt.Print(err)
	}
}

func (r *Review) Delete(uid string) {
	_, err := queries.DeleteReview().Exec(uid)
	if err != nil {
		fmt.Print(err)
	}
}

func (r *Review) RowData(uid int) {
	row := queries.Review().QueryRow(uid)
	row.Scan(
		&r.Beneficial_point,
		&r.Honey_point,
		&r.Assignment,
		&r.Team_project,
		&r.Pressentation,
		&r.Comment)
}
