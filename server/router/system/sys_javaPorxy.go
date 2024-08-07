package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/slyrx/gin_exam_system/server/api/v1"
)

type JavaProxyRouter struct{}

func (s *JavaProxyRouter) InitJavaProxyUserRouter(Router *gin.RouterGroup) {
	javaProxyRouter := Router.Group("api/admin/user")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	javaProxyApi := v1.ApiGroupApp.SystemApiGroup.JavaProxyApi

	javaProxyRouter.Any("/*proxyPath", func(c *gin.Context) {
		proxyPath := c.Param("proxyPath")
		if proxyPath == "/edit" {
			baseApi.CreateUser(c)
		} else {
			javaProxyApi.ReverseProxy("http://127.0.0.1:8003/")(c)
		}
	})

}
