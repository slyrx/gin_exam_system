package response

import "time"

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type PageQuestionResult struct {
	Total             int        `json:"total"`
	List              []Question `json:"list"` // 使用 any 类型，因为列表内容未指定
	PageNum           int        `json:"pageNum"`
	PageSize          int        `json:"pageSize"`
	Size              int        `json:"size"`
	StartRow          int        `json:"startRow"`
	EndRow            int        `json:"endRow"`
	Pages             int        `json:"pages"`
	PrePage           int        `json:"prePage"`
	NextPage          int        `json:"nextPage"`
	IsFirstPage       bool       `json:"isFirstPage"`
	IsLastPage        bool       `json:"isLastPage"`
	HasPreviousPage   bool       `json:"hasPreviousPage"`
	HasNextPage       bool       `json:"hasNextPage"`
	NavigatePages     int        `json:"navigatePages"`
	NavigatePageNums  []int      `json:"navigatepageNums"`
	NavigateFirstPage int        `json:"navigateFirstPage"`
	NavigateLastPage  int        `json:"navigateLastPage"`
}

type Question struct {
	ID                   int       `json:"id"`
	QuestionType         int       `json:"questionType"`
	TextContentID        *int      `json:"textContentId"`
	CreateTime           time.Time `json:"createTime"`
	SubjectID            int       `json:"subjectId"`
	CreateUser           int       `json:"createUser"`
	Score                string    `json:"score"`
	Status               int       `json:"status"`
	Correct              string    `json:"correct"`
	AnalyzeTextContentID *int      `json:"analyzeTextContentId"`
	Difficult            int       `json:"difficult"`
	ShortTitle           string    `json:"shortTitle"`
}
