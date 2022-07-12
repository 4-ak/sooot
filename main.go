package main

import (
	"github.com/4-ak/sooot/domain"
)

func main() {
	server := NewServer(&domain.CourseList{
		[]domain.Course{
			{"집 보내줘", "김교수", domain.ReviewList{
				[]domain.Review{
					{4, "가나다"},
					{3, "다나가"},
					{4, "별로"},
				},
			}},
			{"프로그래밍", "조교수", domain.ReviewList{}},
			{"C 언어", "이교수", domain.ReviewList{}},
			{"Java 활용", "박교수", domain.ReviewList{}},
			{"나시Golang", "최교수", domain.ReviewList{}},
		},
	})

	server.App.Listen(":8080")
}
