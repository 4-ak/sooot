package model

import (
	"fmt"

	"github.com/4-ak/sooot/db/queries"
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
	rows, err := queries.MajorAll().Query()
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
