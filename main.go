package main

import (
	"github.com/slyrx/gin_exam_system/server/others/core"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/initialize"
	"go.uber.org/zap"
)

func main() {
	global.GES_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GES_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GES_LOG)

	global.GES_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.GES_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GES_DB.DB()
		defer db.Close()
	}

	core.RunServer()
}
