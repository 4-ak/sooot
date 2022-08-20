package model

import (
	"fmt"

	"github.com/4-ak/sooot/db"
)

type Lecture_data struct {
	Uid      int
	Base     string
	Year     int
	Semester string
	Credit   int
	Major    string
}

func (ld *Lecture_data) Insert(name, professor string) {
	ld.Base = name
	_, err := db.InsertLecture().Exec(
		ld.Base,
		professor,
		ld.Year,
		ld.Semester,
		ld.Credit,
		ld.Major)
	if err != nil {
		fmt.Println(err)
	}
}

func (ld *Lecture_data) Update(uid int) {
	err := db.UpdateLecture().QueryRow(
		ld.Year,
		ld.Semester,
		ld.Credit,
		ld.Major,
		uid).Scan(
		&ld.Base)
	if err != nil {
		fmt.Print(err)
	}

}

func (ld *Lecture_data) Delete(uid int) {
	_, err := db.DeleteLecture().Exec(uid)
	if err != nil {
		fmt.Print(err)
	}
}
