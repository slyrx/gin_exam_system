package utils

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)         // 去掉字符串首尾的空格
	dr, err := time.ParseDuration(d) // 尝试将字符串解析为 time.Duration 类型
	if err == nil {
		return dr, nil // 如果解析成功，返回解析结果
	}
	if strings.Contains(d, "d") { // 如果字符串包含 'd'，说明可能是天数
		index := strings.Index(d, "d") // 找到 'd' 的位置

		hour, _ := strconv.Atoi(d[:index])          // 将 'd' 之前的部分转换为整数，表示天数
		dr = time.Hour * 24 * time.Duration(hour)   // 将天数转换为小时数
		ndr, err := time.ParseDuration(d[index+1:]) // 解析 'd' 之后的部分
		if err != nil {
			return dr, nil // 如果解析失败，只返回天数部分的时间段
		}
		return dr + ndr, nil // 返回天数部分和 'd' 之后部分的总时间段
	}

	dv, err := strconv.ParseInt(d, 10, 64) // 尝试将整个字符串解析为整数
	return time.Duration(dv), err          // 返回解析结果，单位为纳秒
}
