package model

import (
	"fmt"

	"github.com/4-ak/sooot/db"
)

type Lecture_base struct {
	Uid       int
	Name      string
	Professor string
	Major     string
}

func (lb *Lecture_base) Insert() {
	err := db.InsertLecture_base().QueryRow(
		lb.Name,
		lb.Professor).Scan(
		&lb.Uid)
	if err != nil {
		fmt.Println(err)
	}
}

func (lb *Lecture_base) Update(uid int) {
	_, err := db.UpdateLecture_base().Exec(
		lb.Name,
		lb.Professor,
		uid)
	if err != nil {
		fmt.Print(err)
	}

}

func (lb *Lecture_base) Delete(uid int) {
	_, err := db.DeleteLecture_base().Exec(uid)
	if err != nil {
		fmt.Print(err)
	}
}

func (lb *Lecture_base) Lecture_base() []Lecture_base {
	rows, err := db.Lecture_baseAll().Query()
	arr := make([]Lecture_base, 0)
	for rows.Next() {
		rows.Scan(
			&lb.Name,
			&lb.Professor,
			&lb.Major)
		arr = append(arr, *lb)
	}
	if err != nil {
		fmt.Print(err)
	}
	return arr
}
