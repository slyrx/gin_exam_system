package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

func (s *JavaProxyRouter) InitJavaProxyEducationRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/education")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyExamRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/exam/paper/page")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyTaskRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/task")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyMessageRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/message")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyDashboardRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/dashboard")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyQuestionSelectRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/question/select")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyQuestionEditRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/question/edit")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyQuestionDeleteRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/question/delete")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyUploadRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/upload")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyExamPaperSelectRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/exam/paper/select")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyExamPaperDeleteRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/exam/paper/delete")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyExamPaperTaskExamPageRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/exam/paper/taskExamPage")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}

func (s *JavaProxyRouter) InitJavaProxyExamPaperAnswerRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/examPaperAnswer")
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi
	{
		javaProxyRouter.Any("/*proxyPath", javaProxyApi.ReverseProxy("http://127.0.0.1:8003/"))
	}

}
