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
		examPaperRouter.POST("paper/get", examPaperApi.CreateErrorQuestionPaper)                  // 以总错题自动组卷
		examPaperRouter.POST("paper/getByUserId", examPaperApi.CreateErrorQuestionPaper)          // 以用户错题自动组卷
		examPaperRouter.POST("paper/assign-paper-visibility", examPaperApi.AssignPaperVisibility) // 管理员指定试卷对用户的可见情况
	}
}
