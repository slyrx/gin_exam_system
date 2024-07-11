package response

import (
	"github.com/slyrx/gin_exam_system/server/model/system"
)

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}
