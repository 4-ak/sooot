package model

import (
	"fmt"

	"github.com/4-ak/sooot/db"
)

type Major struct {
	Uid  int
	Name string
}

func NewMajor() Major {
	major := Major{}
	return major
}

func (s *Major) Major() []Major {
	rows, err := db.MajorAll().Query()
	arr := make([]Major, 0)
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
