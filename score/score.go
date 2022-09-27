package score

import (
	"math"

	"github.com/4-ak/sooot/db/model"
)

type Score_Avg struct {
	Assignment_Avg   float64
	Team_project_Avg float64
	Presentation_Avg float64
	Beneficial_Avg   float64
	Honey_Avg        float64
}

type Score_text struct {
	Assignment   string
	Team_project string
	Presentation string
	Beneficial   float64
	Honey        float64
}

func Score(review_data []model.Review) Score_text {
	score_avg := Score_Avg{}
	score_text := Score_text{}
	var count float64
	for i, score := range review_data {
		score_avg.Assignment_Avg += float64(score.Assignment)
		score_avg.Team_project_Avg += float64(score.Team_project)
		score_avg.Presentation_Avg += float64(score.Presentation)
		score_avg.Beneficial_Avg += float64(score.Beneficial)
		score_avg.Honey_Avg += float64(score.Honey)
		count = float64(i)
	}
	score_avg.Assignment_Avg = math.Round((score_avg.Assignment_Avg/count)*100) / 100
	score_avg.Team_project_Avg = math.Round((score_avg.Team_project_Avg)/count*100) / 100
	score_avg.Presentation_Avg = math.Round((score_avg.Presentation_Avg/count)*100) / 100
	score_avg.Beneficial_Avg = math.Round((score_avg.Beneficial_Avg/count)*100) / 100
	score_avg.Honey_Avg = math.Round((score_avg.Honey_Avg/count)*100) / 100

	score_text.Beneficial = score_avg.Beneficial_Avg
	score_text.Honey = score_avg.Honey_Avg
	score_text.Assignment = compare(score_avg.Assignment_Avg)
	score_text.Team_project = compare(score_avg.Team_project_Avg)
	score_text.Presentation = compare(score_avg.Presentation_Avg)

	return score_text
}

func compare(score float64) string {
	if score > 2.34 {
		return "많음"
	} else if score <= 2.33 && score > 1.67 {
		return "보통"
	} else if score <= 1.66 {
		return "적음"
	} else {
		return "확인요망"
	}
}
