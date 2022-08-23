package model

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
)

type Review struct {
	Uid          int
	Lecture_uid  int
	Writer       int
	Beneficial   int //1~5
	Honey        int //1~5
	Assignment   int //1~3
	Team_project int //1~3
	Presentation int //1~3
	Comment      string
	Created_at   string
	Lecture_name string
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
			&r.Lecture_uid,
			&r.Writer,
			&r.Beneficial,
			&r.Honey,
			&r.Assignment,
			&r.Team_project,
			&r.Presentation,
			&r.Comment,
			&r.Created_at,
			&r.Lecture_name)
		arr = append(arr, *r)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
	}
	return arr
}

func (r *Review) Insert(lect_uid, account_uid string) {
	_, err := queries.InsertReview().Exec(
		lect_uid,
		account_uid,
		r.Beneficial,
		r.Honey,
		r.Assignment,
		r.Team_project,
		r.Presentation,
		r.Comment)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Review) Update(account_uid string) {
	_, err := queries.UpdateReview().Exec(
		r.Beneficial,
		r.Honey,
		r.Assignment,
		r.Team_project,
		r.Presentation,
		r.Comment,
		account_uid)
	if err != nil {
		fmt.Print(err)
	}
}

func (r *Review) Delete(account_uid string) {
	_, err := queries.DeleteReview().Exec(account_uid)
	if err != nil {
		fmt.Print(err)
	}
}

func (r *Review) RowData(uid int) {
	row := queries.Review().QueryRow(uid)
	row.Scan(
		&r.Beneficial,
		&r.Honey,
		&r.Assignment,
		&r.Team_project,
		&r.Presentation,
		&r.Comment)
}
