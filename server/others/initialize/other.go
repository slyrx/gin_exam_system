package initialize

import (
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func OtherInit() {
	// 解析配置中的 JWT 过期时间，并将其转换为时间段（Duration）
	dr, err := utils.ParseDuration(global.GES_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err) // 如果解析失败，则抛出异常并终止程序
	}

	// 解析配置中的 JWT 缓冲时间，并将其转换为时间段（Duration）
	_, err = utils.ParseDuration(global.GES_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err) // 如果解析失败，则抛出异常并终止程序
	}

	// 初始化全局的 BlackCache（本地缓存），设置默认过期时间为解析后的过期时间
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr), // 设置缓存项的默认过期时间
	)
}
