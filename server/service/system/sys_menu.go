package system

import (
	"errors"

	"github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"gorm.io/gorm"
)

type MenuService struct{}

var MenuServiceApp = new(MenuService)

// UserAuthorityDefaultRouter 用户角色默认路由检查
//
//	Author [SliverHorn](https://github.com/SliverHorn)
func (menuService *MenuService) UserAuthorityDefaultRouter(user *system.SysUser) {
	var menuIds []string
	err := global.GES_DB.Model(&system.SysAuthorityMenu{}).Where("sys_authority_authority_id = ?", user.AuthorityId).Pluck("sys_base_menu_id", &menuIds).Error
	if err != nil {
		return
	}
	var am system.SysBaseMenu
	err = global.GES_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, menuIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}
}
