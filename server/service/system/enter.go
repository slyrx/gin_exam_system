package system

type ServiceGroup struct {
	UserService
	JwtService
	CasbinService
	ExamService
	QuestionService
	ExamPaperService
	ExamPaperStdService
}
