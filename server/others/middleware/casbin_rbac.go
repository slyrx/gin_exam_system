package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	"github.com/slyrx/gin_exam_system/server/others/global"

	// "github.com/slyrx/gin_exam_system/server/others/utils"
	"github.com/slyrx/gin_exam_system/server/others/utils"
	"github.com/slyrx/gin_exam_system/server/service"
	"go.uber.org/zap"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	fmt.Println("ccc0")
	return func(c *gin.Context) {
		global.GES_LOG.Info("CasbinHandler")
		// waitUse, _ := utils.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GES_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		// sub := strconv.Itoa(888)
		sub := utils.GetCurrentUser(c).UserName
		global.GES_LOG.Info("CasbinHandler", zap.Any("sub", sub), zap.Any("obj", obj), zap.Any("act", act))
		e := casbinService.Casbin() // 判断策略中是否存在
		e.EnableLog(true)
		policies := e.GetPolicy()
		fmt.Println("Loaded policies:")
		for _, policy := range policies {
			fmt.Println(policy)
		}
		success, err := e.Enforce(sub, obj, act)
		global.GES_LOG.Info("CasbinHandler", zap.Any("success", success), zap.Any("err", err))
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
