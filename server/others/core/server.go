package core

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/initialize"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Routers := initialize.Routers()
	Routers.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GES_CONFIG.System.Addr)
	s := initServer(address, Routers)
	global.GES_LOG.Info("server run success on ", zap.String("address", address))
	global.GES_LOG.Error(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) server {
	// 创建一个新的 Endless Server 实例，传入监听地址和 Gin 路由器
	s := endless.NewServer(address, router)

	// 设置读取请求头的超时时间为 20 秒
	s.ReadHeaderTimeout = 20 * time.Second

	// 设置写响应的超时时间为 20 秒
	s.WriteTimeout = 20 * time.Second

	// 设置请求头的最大字节数为 1MB
	s.MaxHeaderBytes = 1 << 20

	// 返回配置好的 Endless Server 实例
	return s
}
