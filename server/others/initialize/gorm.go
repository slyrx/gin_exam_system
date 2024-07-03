package initialize

import (
	"github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
	"os"

	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GES_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.GES_DB
	err := db.AutoMigrate(

		system.SysApi{},
	// system.SysUser{},
	// system.SysBaseMenu{},
	// system.JwtBlacklist{},
	// system.SysAuthority{},
	// system.SysDictionary{},
	// system.SysOperationRecord{},
	// system.SysAutoCodeHistory{},
	// system.SysDictionaryDetail{},
	// system.SysBaseMenuParameter{},
	// system.SysBaseMenuBtn{},
	// system.SysAuthorityBtn{},
	// system.SysAutoCode{},
	// system.SysExportTemplate{},
	// system.Condition{},
	// system.JoinTemplate{},

	// example.ExaFile{},
	// example.ExaCustomer{},
	// example.ExaFileChunk{},
	// example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.GES_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GES_LOG.Info("register table success")
}
