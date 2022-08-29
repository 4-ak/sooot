package lectbrowser

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/4-ak/sooot/db/model"
)

type LectureNode struct {
	Data       *model.Lecture_base
	Consonants string
}

type LectureBrowser struct {
	data []LectureNode
}

type Handler struct{}

func (lb *LectureBrowser) Init(lectures []model.Lecture_base) {
	fmt.Println(lectures)
	for i, v := range lectures {
		lb.data = append(lb.data, LectureNode{
			Data:       &lectures[i],
			Consonants: HangulsToConsonants(v.Name),
		})
	}
}

func (lb *LectureBrowser) FindByFirstInitial(name rune) []model.Lecture_base {
	ch := name
	if IsHangul(name) {
		ch = HangulToConsonant(name)
	}

	consonantResult := make([]model.Lecture_base, 0, 8)

	for _, v := range lb.data {
		if []rune(v.Consonants)[0] == ch {
			consonantResult = append(consonantResult, *v.Data)
		}
	}

	if ch == name {
		return consonantResult
	}

	result := make([]model.Lecture_base, 0, 8)
	for _, v := range consonantResult {
		if []rune(v.Name)[0] == name {
			result = append(result, v)
		}
	}
	return result
}

func (lb *LectureBrowser) FindPriorityByConsonants(query string) []model.Lecture_base {
	consonants := HangulsToConsonants(query)

	temp := NewPriorityQueue(true)
	for _, v := range lb.data {

		index := strings.Index(v.Consonants, consonants)
		if index == -1 {
			continue
		}

		temp.Enqueue(PriorityLecture{
			priority: index,
			lecture:  v.Data,
		})
	}

	result := NewPriorityQueue(true)
	queryRune := []rune(query)
	iterator := temp.Iterator()
	for iterator.Next() {
		pl := iterator.Value().(PriorityLecture)
		count := 0
		nameRune := []rune(pl.lecture.Name)
		offset := utf8.RuneCountInString(pl.lecture.Name[:pl.priority])
		for i, v := range nameRune[offset:] {
			if i >= len(queryRune) {
				break
			}
			if v == queryRune[i] {
				count++
			}
		}
		result.Enqueue(PriorityLecture{
			priority: count,
			lecture:  pl.lecture,
		})
	}

	ordered := make([]model.Lecture_base, 0, 8)
	iterator = result.Iterator()
	for iterator.Next() {
		ordered = append(ordered, *iterator.Value().(PriorityLecture).lecture)
	}
	temp.Clear()
	result.Clear()
	return ordered
}
