package config

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
	RetentionDay  int    `mapstructure:"retention-day" json:"retention-day" yaml:"retention-day"`    // 日志保留天数
}

func (c *Zap) Encoder() zapcore.Encoder {
	// 定义一个 EncoderConfig 结构体，配置日志的各种属性
	config := zapcore.EncoderConfig{
		TimeKey:       "time",          // 时间键名
		NameKey:       "name",          // 日志记录器名称键名
		LevelKey:      "level",         // 日志级别键名
		CallerKey:     "caller",        // 调用者信息键名
		MessageKey:    "message",       // 日志消息键名
		StacktraceKey: c.StacktraceKey, // 堆栈跟踪信息键名，使用 Zap 结构体中的 StacktraceKey 字段

		LineEnding: zapcore.DefaultLineEnding, // 换行符，默认为系统默认换行符
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			// 定义时间编码器，将时间按指定格式编码为字符串，加上前缀 Prefix
			encoder.AppendString(c.Prefix + t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeLevel:    c.LevelEncoder(),               // 使用 Zap 结构体中的 LevelEncoder 方法来编码日志级别
		EncodeCaller:   zapcore.FullCallerEncoder,      // 调用者信息编码器，完整路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder, // 时间段编码器，使用秒作为单位的时间编码器
	}

	// 根据 Zap 结构体中的 Format 字段判断日志格式，选择相应的编码器
	if c.Format == "json" {
		return zapcore.NewJSONEncoder(config) // 返回一个新的 JSON 编码器实例，使用上面定义的 config 配置
	}

	// 默认情况下，返回一个新的控制台（Console）编码器实例，使用上面定义的 config 配置
	return zapcore.NewConsoleEncoder(config)
}

// LevelEncoder 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author [SliverHorn](https://github.com/SliverHorn)
func (c *Zap) LevelEncoder() zapcore.LevelEncoder {
	switch {
	case c.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case c.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case c.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case c.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// Levels 根据字符串转化为 zapcore.Levels
func (c *Zap) Levels() []zapcore.Level {
	// 创建一个长度为 0，容量为 7 的 zapcore.Level 切片，用于存储不同的日志级别。
	levels := make([]zapcore.Level, 0, 7)
	// 将 c.Level 字符串解析为 zapcore.Level 类型的 level 变量，如果解析失败，默认设置为 zapcore.DebugLevel。
	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		// 如果解析过程中发生错误，则将 level 设置为 zapcore.DebugLevel。
		level = zapcore.DebugLevel
	}
	// 从当前的日志级别开始，依次增加，直到 zapcore.FatalLevel，并将每个级别添加到 levels 切片中。
	for ; level <= zapcore.FatalLevel; level++ {
		// 将当前的 level 添加到 levels 切片中。
		levels = append(levels, level)
	}
	// 返回包含所有级别的切片。
	return levels
}
