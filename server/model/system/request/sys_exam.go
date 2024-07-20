package request

// 定义请求和响应结构
type AnswerSubmitRequest struct {
	ID          int          `json:"id"`
	DoTime      int          `json:"doTime"`
	AnswerItems []AnswerItem `json:"answerItems"`
}

type AnswerItem struct {
	QuestionID   int      `json:"questionId"`
	Content      string   `json:"content"`
	ContentArray []string `json:"contentArray"`
	Completed    bool     `json:"completed"`
	ItemOrder    int      `json:"itemOrder"`
}
