package system

import "github.com/slyrx/gin_exam_system/server/service"

type ApiGroup struct {
	BaseApi
	JwtApi
	ExamApi
	QuestionApi
}

var (
	userService     = service.ServiceGroupApp.SystemServiceGroup.UserService
	jwtService      = service.ServiceGroupApp.SystemServiceGroup.JwtService
	examService     = service.ServiceGroupApp.SystemServiceGroup.ExamService
	questionService = service.ServiceGroupApp.SystemServiceGroup.QuestionService
)
