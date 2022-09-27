package model

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
)

type Lecture struct {
	Base Lecture_base
	Data Lecture_data
}

func NewLecture() Lecture {
	lecture := Lecture{}
	return lecture
}

func (l *Lecture) SelectData() []Lecture {
	rows, err := queries.LectureAll().Query()
	arr := make([]Lecture, 0)
	for rows.Next() {
		rows.Scan(
			&l.Data.Uid,
			&l.Data.Year,
			&l.Data.Base,
			&l.Data.Semester,
			&l.Data.Credit,
			&l.Data.Major,
			&l.Base.Uid,
			&l.Base.Name,
			&l.Base.Professor)
		arr = append(arr, *l)
	}
	if err != nil || len(arr) == 0 {
		fmt.Println(err)
	}
	return arr
}

func (l *Lecture) RowData(uid string) Lecture {
	row := queries.Lecture().QueryRow(uid)
	row.Scan(
		&l.Data.Year,
		&l.Data.Semester,
		&l.Data.Credit,
		&l.Data.Major,
		&l.Base.Name,
		&l.Base.Professor)
	return *l
}

func (l *Lecture) CompareData() {
	row := queries.CompareLecture().QueryRow(
		l.Base.Name,
		l.Base.Professor,
		l.Data.Year,
		l.Data.Semester,
		l.Data.Credit,
		l.Data.Major)
	row.Scan(&l.Data.Uid)
}
