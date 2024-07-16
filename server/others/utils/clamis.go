package utils

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	systemReq "github.com/slyrx/gin_exam_system/server/model/system/request"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
)

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.Request.Header.Get("x-token")
	}
	return token
}

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GES_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

func GetCurrentUser(c *gin.Context) *systemMod.SysExamUser {

	// 获取名为 "example_cookie" 的 cookie
	cookie, err := c.Cookie("studentUserName")
	global.GES_LOG.Info("cookie1", zap.Any("cookie", err))
	if err != nil {
		if err == http.ErrNoCookie {
			// 如果 cookie 不存在，返回相应的响应
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Cookie not found",
			})
			return nil
		}
		// 其他错误情况
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred",
		})
		return nil
	}
	var user systemMod.SysExamUser
	user.UserName = cookie
	global.GES_LOG.Info("cookie2", zap.Any("cookie", cookie))
	global.GES_LOG.Info("cookie2", zap.Any("cookie", user))
	return &user
}
