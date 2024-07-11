package system

import (
	"errors"
	"fmt"

	"github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
)

type UserService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GES_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.GES_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}
