package internal

import (
	"fmt"
	"github.com/slyrx/gin_exam_system/server/others/config"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type Writer struct {
	config config.GeneralDB
	writer logger.Writer
}

func NewWriter(config config.GeneralDB, writer logger.Writer) *Writer {
	return &Writer{config: config, writer: writer}
}

func (c *Writer) Printf(message string, data ...any) {
	if c.config.LogZap {
		switch c.config.LogLevel() {
		case logger.Silent:
			zap.L().Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			zap.L().Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			zap.L().Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			zap.L().Info(fmt.Sprintf(message, data...))
		default:
			zap.L().Info(fmt.Sprintf(message, data...))
		}
		return
	}
	c.writer.Printf(message, data...)
}
