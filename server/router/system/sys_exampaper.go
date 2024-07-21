package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

type ExamPaperRouter struct{}

func (s *ExamPaperRouter) InitExamPaperRouter(Router *gin.RouterGroup) {
	examPaperRouter := Router.Group("api/admin/exam")
	examPaperApi := v1.ApiGroupApp.SystemApiGroup.ExamPaperApi
	{
		examPaperRouter.POST("paper/edit", examPaperApi.SetExamPaper)
	}
}
