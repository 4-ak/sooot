package model

import (
	"fmt"

	"github.com/4-ak/sooot/db"
)

type Lecture struct {
	Uid            int
	Department     int
	Name           string
	Professor_name string
	Semester       int
	Credit         int
	Parent         int
}

func NewLecture() Lecture {
	lecture := Lecture{}
	return lecture
}

func (l *Lecture) SelectData() []Lecture {
	row, err := db.LectureAll().Query()
	arr := make([]Lecture, 0)
	for row.Next() {
		row.Scan(
			&l.Uid,
			&l.Department,
			&l.Name,
			&l.Professor_name,
			&l.Semester,
			&l.Credit,
			&l.Parent)
		arr = append(arr, *l)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
	}
	return arr
}

func (l *Lecture) Insert() {
	_, err := db.InsertLecture().Exec(
		l.Name,
		l.Professor_name,
		l.Semester,
		l.Credit)
	if err != nil {
		fmt.Print(err)
	}

}

func (l *Lecture) Update(uid int) {
	_, err := db.UpdateLecture().Exec(
		l.Name,
		l.Professor_name,
		l.Semester,
		l.Credit,
		uid)
	if err != nil {
		fmt.Print(err)
	}

}

func (l *Lecture) Delete(uid int) {
	_, err := db.DeleteLecture().Exec(uid)
	if err != nil {
		fmt.Print(err)
	}

}

func (l *Lecture) RowData(uid int) {
	row := db.Lecture().QueryRow(uid)
	row.Scan(
		&l.Uid,
		&l.Department,
		&l.Name,
		&l.Professor_name,
		&l.Semester,
		&l.Credit,
		&l.Parent)
}
