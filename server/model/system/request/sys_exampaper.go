package request

type CreateExamPaperRequest struct {
	ID            *int        `json:"id"`
	Level         int         `json:"level"`
	SubjectID     int         `json:"subjectId"`
	PaperType     int         `json:"paperType"`
	LimitDateTime []string    `json:"limitDateTime"`
	Name          string      `json:"name"`
	SuggestTime   interface{} `json:"suggestTime"`
	TitleItems    []TitleItem `json:"titleItems"`
}

type TitleItem struct {
	Name          string           `json:"name"`
	QuestionItems []QuestionItem_1 `json:"questionItems"`
}

type QuestionItem_1 struct {
	ID           int    `json:"id"`
	QuestionType int    `json:"questionType"`
	SubjectID    int    `json:"subjectId"`
	Title        string `json:"title"`
	GradeLevel   int    `json:"gradeLevel"`
	Items        []Item `json:"items"`
	Analyze      string `json:"analyze"`
	Correct      string `json:"correct"`
	Score        string `json:"score"`
	Difficult    int    `json:"difficult"`
	ItemOrder    *int   `json:"itemOrder"`
}

type Item struct {
	Prefix   string  `json:"prefix"`
	Content  string  `json:"content"`
	Score    *int    `json:"score"`
	ItemUUID *string `json:"itemUuid"`
}

type TextContent struct {
	TitleContent        string               `json:"titleContent"`
	Analyze             string               `json:"analyze"`
	QuestionItemObjects []QuestionItemObject `json:"questionItemObjects"`
	Correct             string               `json:"correct"`
}

type QuestionItemObject struct {
	Prefix   string  `json:"prefix"`
	Content  string  `json:"content"`
	Score    *int    `json:"score"`
	ItemUUID *string `json:"itemUuid"`
}

type CreateErrorQuestionPaperRequest struct {
	SubjectID  int `json:"subjectId" binding:"required"`
	GradeLevel int `json:"gradeLevel" binding:"required"`
	UserID     int `json:"userId" binding:"required"`
	Limit      int `json:"limit" binding:"required"` // 限制题目数量
}
