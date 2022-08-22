package model

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
)

type Semester struct {
	Uid  int
	Name string
}

func NewSemester() Semester {
	semester := Semester{}
	return semester
}

func (s *Semester) Semester() []Semester {
	rows, err := queries.SemesterAll().Query()
	arr := make([]Semester, 0)
	for rows.Next() {
		rows.Scan(
			&s.Name)
		arr = append(arr, *s)
	}
	if err != nil {
		fmt.Print(err)
	}
	return arr
}
