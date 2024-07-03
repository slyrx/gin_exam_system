package internal

import (
	"github.com/slyrx/gin_exam_system/server/others/config"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB
	switch global.GES_CONFIG.System.DbType {
	case "mysql":
		general = global.GES_CONFIG.Mysql.GeneralDB
	default:
		general = global.GES_CONFIG.Mysql.GeneralDB
	}
	return &gorm.Config{
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
