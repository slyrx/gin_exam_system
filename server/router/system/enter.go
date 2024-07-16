package system

type RouterGroup struct {
	BaseRouter
	JwtRouter
	ExamRouter
	QuestionRouter
}
