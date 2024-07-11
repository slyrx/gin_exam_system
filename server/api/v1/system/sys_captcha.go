package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemRes "github.com/slyrx/gin_exam_system/server/model/system/response"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
)

type BaseApi struct{}

var store = base64Captcha.DefaultMemStore

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (b *BaseApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := global.GES_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GES_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.GES_CONFIG.Captcha.ImgHeight, global.GES_CONFIG.Captcha.ImgWidth, global.GES_CONFIG.Captcha.KeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.GES_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GES_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证码获取成功", c)
}
