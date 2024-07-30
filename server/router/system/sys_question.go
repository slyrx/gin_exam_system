package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

type QuestionRouter struct{}

func (s *QuestionRouter) InitQuestionRouter(Router *gin.RouterGroup) {
	questionRouter := Router.Group("api/admin/question")
	questionApi := v1.ApiGroupApp.SystemApiGroup.QuestionApi
	{
		questionRouter.POST("page", questionApi.GetPageInfo)
		questionRouter.GET("createUserQuestionTable", questionApi.CreateUserQuestionTable)
	}

}
