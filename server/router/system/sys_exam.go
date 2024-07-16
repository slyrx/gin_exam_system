package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

type ExamRouter struct{}

func (s *ExamRouter) InitExamRouter(Router *gin.RouterGroup) {
	examRouter := Router.Group("api/student")
	examApi := v1.ApiGroupApp.SystemApiGroup.ExamApi
	{
		examRouter.POST("exampaper/answer/answerSubmit", examApi.AnswerSubmit)
	}
}
