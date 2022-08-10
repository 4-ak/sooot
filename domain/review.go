package domain

type Review struct {
	Score   int
	Comment string
}

type ReviewList struct {
	List []Review
}

func NewReview(score int, comment string) *Review {
	if score > 0 && score <= 5 {
		return &Review{score, comment}
	}
	return nil
}

func (r *ReviewList) Input(sc int, cm string) bool {
	if review := NewReview(sc, cm); review != nil {
		r.List = append(r.List, *review)
		return true
	}
	return false
}

func (r *ReviewList) Update(sc, num int, cm string) bool {
	if review := NewReview(sc, cm); review != nil {
		r.List[num-1] = *review
		return true
	}
	return false
}

func (r *ReviewList) Delete(num int) {
	aa := r.List
	r.List = aa[:num-1]
	aa = aa[num:]
	r.List = append(r.List, aa...)
}

func (r *ReviewList) Average() float32 {
	var sum float32
	for i := 0; i < len(r.List); i++ {
		sum += float32(r.List[i].Score)
	}
	return sum / float32(len(r.List))
}
