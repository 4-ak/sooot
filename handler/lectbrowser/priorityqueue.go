package lectbrowser

import (
	"github.com/4-ak/sooot/db/model"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

type PriorityLecture struct {
	priority int
	lecture  *model.Lecture_base
}

func NewPriorityQueue(isMax bool) *pq.Queue {
	if isMax {
		return pq.NewWith(byPriorityMax)
	}
	return pq.NewWith(byPriorityMin)
}

func byPriorityMax(a, b interface{}) int {
	pa := a.(PriorityLecture).priority
	pb := b.(PriorityLecture).priority
	return utils.IntComparator(pa, pb)
}

func byPriorityMin(a, b interface{}) int {
	pa := a.(PriorityLecture).priority
	pb := b.(PriorityLecture).priority
	return -utils.IntComparator(pa, pb)
}
