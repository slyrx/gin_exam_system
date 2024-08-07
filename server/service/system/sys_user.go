package system

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
)

type UserService struct{}

var UserServiceApp = new(UserService)

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

func (userService *UserService) CreateUser(user *system.SysExamUser) error {
	user.UserUUID = uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = global.GES_DB.Create(user).Error
	if err != nil {
		return err
	}
	return UserServiceApp.AddAdminStudentRelation(uint(user.ID), user.RealName)
}

func (userService *UserService) FindUserByRealName(realName string) (*system.SysExamUser, error) {
	var user system.SysExamUser
	result := global.GES_DB.Debug().Where("real_name = ?", realName).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 没有找到记录，返回 nil
		}
		return nil, result.Error // 其他错误
	}
	return &user, nil
}

func (userService *UserService) AddAdminStudentRelation(studentID uint, RealName string) error {
	admin, err := UserServiceApp.FindUserByRealName(RealName)
	if err != nil {
		return err
	}
	relation := system.AdminStudentRelation{
		AdminID:   uint(admin.ID),
		StudentID: studentID,
	}

	result := global.GES_DB.Create(&relation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userService *UserService) AddOrUpdateAdminStudentRelation(studentID uint, realName string) error {
	admin, err := userService.FindUserByRealName(realName)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	var relation system.AdminStudentRelation
	err = global.GES_DB.Debug().Where("admin_id = ? AND student_id = ?", admin.ID, studentID).First(&relation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果记录不存在，则创建
		relation = system.AdminStudentRelation{
			AdminID:   uint(admin.ID),
			StudentID: studentID,
		}
		return global.GES_DB.Debug().Create(&relation).Error
	} else if err != nil {
		return err
	}

	// 如果记录存在，则更新
	relation.UpdatedAt = time.Now()
	return global.GES_DB.Save(&relation).Error
}
