package global

import (
	"github.com/slyrx/gin_exam_system/server/others/config"
	"github.com/slyrx/gin_exam_system/server/others/utils/timer"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/songzhibin97/gkit/cache/local_cache"
)

var (
	GES_CONFIG config.Server
	GES_LOG    *zap.Logger
	GES_VP     *viper.Viper
	GES_DB     *gorm.DB
	GES_Timer  timer.Timer = timer.NewTimerTask()

	BlackCache local_cache.Cache
)
