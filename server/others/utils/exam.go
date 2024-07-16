package utils

import "fmt"

func ContentToString(contentArray []string) string {
	// Implement the logic to convert contentArray to a string.
	return ""
}

func ScoreToVM(score int) string {
	if score%10 == 0 {
		return fmt.Sprintf("%d", score/10)
	} else {
		return fmt.Sprintf("%.1f", float64(score)/10.0)
	}
}

func SecondToVM(second int) string {
	// 计算天、小时、分钟和秒
	days := second / (60 * 60 * 24)
	hours := (second % (60 * 60 * 24)) / (60 * 60)
	minutes := (second % (60 * 60)) / 60
	seconds := second % 60

	var dateTimes string

	// 根据时间长度选择合适的格式
	if days > 0 {
		dateTimes = fmt.Sprintf("%d天 %d时 %d分 %d秒", days, hours, minutes, seconds)
	} else if hours > 0 {
		dateTimes = fmt.Sprintf("%d时 %d分 %d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		dateTimes = fmt.Sprintf("%d分 %d秒", minutes, seconds)
	} else {
		dateTimes = fmt.Sprintf("%d秒", seconds)
	}

	return dateTimes
}
