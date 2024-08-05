package response

import (
	"time"

	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
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

type ExamPaper struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	QuestionCount      int       `json:"questionCount"`
	Score              int       `json:"score"`
	CreateTime         time.Time `json:"createTime"`
	CreateUser         int       `json:"createUser"`
	SubjectID          int       `json:"subjectId"`
	SubjectName        string    `json:"subjectName"`
	PaperType          int       `json:"paperType"`
	FrameTextContentID int       `json:"frameTextContentId"`
}

type CreateErrorQuestionPaperResponse struct {
	PaperID int `json:"paperId"`
}

type PageListResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}

type ResponseData struct {
	Total            int         `json:"total"`
	List             []ExamPaper `json:"list"`
	PageNum          int         `json:"pageNum"`
	PageSize         int         `json:"pageSize"`
	Size             int         `json:"size"`
	StartRow         int         `json:"startRow"`
	EndRow           int         `json:"endRow"`
	Pages            int         `json:"pages"`
	PrePage          int         `json:"prePage"`
	NextPage         int         `json:"nextPage"`
	IsFirstPage      bool        `json:"isFirstPage"`
	IsLastPage       bool        `json:"isLastPage"`
	HasPreviousPage  bool        `json:"hasPreviousPage"`
	HasNextPage      bool        `json:"hasNextPage"`
	NavigatePages    int         `json:"navigatePages"`
	NavigatePageNums []int       `json:"navigatepageNums"`
}

type CustomExamPaper struct {
	systemMod.ExamPaper_1
	SubjectName string `gorm:"column:subject_name"`
}
