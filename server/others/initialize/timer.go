package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/robfig/cron/v3"
	"github.com/slyrx/gin_exam_system/server/others/global"
)

func Timer() {
	go func() { // 启动一个新的 goroutine 以异步运行定时任务
		var option []cron.Option
		option = append(option, cron.WithSeconds()) // 设置 cron 选项，支持秒级别的任务调度

		// 清理数据库的定时任务
		_, err := global.GES_Timer.AddTaskByFunc("ClearDB", "@daily", func() { // 添加一个每天执行一次的任务
			err := task.ClearTable(global.GES_DB) // 定时任务具体执行的方法位于 task 包中
			if err != nil {
				fmt.Println("timer error:", err) // 如果执行任务时出现错误，打印错误信息
			}
		}, "定时清理数据库【日志，黑名单】内容", option...) // 添加任务的描述信息和选项
		if err != nil {
			fmt.Println("add timer error:", err) // 如果添加任务时出现错误，打印错误信息
		}

		// 其他定时任务可以在这里添加，参考上方的使用方法
		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "cron 表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err) // 如果添加任务时出现错误，打印错误信息
		//}
	}()
}
