package core

import (
	"fmt"
	"github.com/slyrx/gin_exam_system/server/others/core/internal"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Zap() (logger *zap.Logger) {
	// 检查配置中的 Director 文件夹是否存在
	if ok, _ := utils.PathExists(global.GES_CONFIG.Zap.Director); !ok {
		// 如果不存在，则打印创建目录的消息，并创建该目录
		fmt.Printf("create %v directory\n", global.GES_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GES_CONFIG.Zap.Director, os.ModePerm)
	}

	// 获取配置中的日志级别
	levels := global.GES_CONFIG.Zap.Levels()
	length := len(levels)
	// 创建一个存放 zapcore.Core 的切片，长度为日志级别的数量
	cores := make([]zapcore.Core, 0, length)

	// 遍历日志级别，为每个级别创建一个 zapcore.Core，并添加到 cores 切片中
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}

	// 创建一个新的 zap.Logger，使用 zapcore.NewTee 组合多个核心
	logger = zap.New(zapcore.NewTee(cores...))

	// 如果配置中指定要显示代码行号，则添加该选项到 logger 中
	if global.GES_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	// 返回创建好的 logger
	return logger
}
