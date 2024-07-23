package response

import (
	"github.com/slyrx/gin_exam_system/server/model/system/request"
)

type CreateExamPaperResponse struct {
	ID            int                 `json:"id"`
	Level         int                 `json:"level"`
	SubjectID     int                 `json:"subjectId"`
	PaperType     int                 `json:"paperType"`
	Name          string              `json:"name"`
	SuggestTime   int                 `json:"suggestTime"`
	LimitDateTime []string            `json:"limitDateTime"`
	TitleItems    []request.TitleItem `json:"titleItems"`
	Score         string              `json:"score"`
}

type CreateErrorQuestionPaperResponse struct {
	PaperID int `json:"paperId"`
}
