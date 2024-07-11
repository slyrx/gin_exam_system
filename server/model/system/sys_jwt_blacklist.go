package system

import (
	"github.com/slyrx/gin_exam_system/server/others/global"
)

type JwtBlacklist struct {
	global.GES_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
