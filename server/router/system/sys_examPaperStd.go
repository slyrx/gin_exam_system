package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

type ExamPaperStdRouter struct{}

func (s *ExamPaperRouter) InitExamPaperStdRouter(Router *gin.RouterGroup) {
	examPaperStdRouter := Router.Group("api/student/exam")
	examPaperStdApi := v1.ApiGroupApp.SystemApiGroup.ExamPaperStdApi
	{
		examPaperStdRouter.POST("paper/pageList", examPaperStdApi.GetExamPaperPageList_1)
	}
}
