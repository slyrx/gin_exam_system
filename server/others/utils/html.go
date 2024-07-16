package utils

import (
	"regexp"
	"strings"
)

func Clear(htmlStr string) string {
	// 定义正则表达式
	regExScript := regexp.MustCompile(`(?i)<script[^>]*?>[\s\S]*?<\/script>`)
	regExStyle := regexp.MustCompile(`(?i)<style[^>]*?>[\s\S]*?<\/style>`)
	regExHTML := regexp.MustCompile(`(?i)<[^>]+>`)

	// 去除<script>标签
	htmlStr = regExScript.ReplaceAllString(htmlStr, "")

	// 去除<style>标签
	htmlStr = regExStyle.ReplaceAllString(htmlStr, "")

	// 去除HTML标签
	htmlStr = regExHTML.ReplaceAllString(htmlStr, "")

	// 去除多余的空白字符并返回
	return strings.TrimSpace(htmlStr)
}
