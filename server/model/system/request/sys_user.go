package request

// User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// QuestionQuery 定义了请求的结构
type QuestionQuery struct {
	ID           *int `json:"id"`
	QuestionType *int `json:"questionType"`
	SubjectID    int  `json:"subjectId"`
	PageIndex    int  `json:"pageIndex"`
	PageSize     int  `json:"pageSize"`
}
