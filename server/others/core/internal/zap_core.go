package internal

import (
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

func NewZapCore(level zapcore.Level) *ZapCore {
	// 创建一个 ZapCore 对象，设置其日志级别为传入的参数 level
	entity := &ZapCore{level: level}

	// 获取写入同步器（WriteSyncer），这里假设 entity.WriteSyncer() 是一个方法或函数，用于返回一个写入同步器
	syncer := entity.WriteSyncer()

	// 创建一个日志级别启用器（LevelEnabler），使用 zap.LevelEnablerFunc 包装函数
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level // 只有当日志级别 l 等于传入的 level 时，启用日志记录
	})

	// 使用全局配置中的 Zap 编码器（Encoder）、上面获取的写入同步器和日志级别启用器创建一个新的核心（Core）
	entity.Core = zapcore.NewCore(global.GES_CONFIG.Zap.Encoder(), syncer, levelEnabler)

	// 返回创建好的 ZapCore 对象
	return entity
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	// 创建一个新的 Cutter 对象，用于日志文件切割
	cutter := NewCutter(
		global.GES_CONFIG.Zap.Director,     // 日志文件目录
		z.level.String(),                   // 当前日志级别的字符串表示
		global.GES_CONFIG.Zap.RetentionDay, // 日志文件保留天数
		CutterWithLayout(time.DateOnly),    // 设置日志文件名的日期布局为仅日期
		CutterWithFormats(formats...),      // 设置日志文件的格式
	)

	// 判断是否需要将日志同时输出到控制台
	if global.GES_CONFIG.Zap.LogInConsole {
		// 如果需要，在标准输出 os.Stdout 和 Cutter 同步器之间创建一个多写同步器
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		// 将多写同步器添加到 zapcore 的同步器中
		return zapcore.AddSync(multiSyncer)
	}

	// 否则，直接将 Cutter 同步器添加到 zapcore 的同步器中
	return zapcore.AddSync(cutter)
}
