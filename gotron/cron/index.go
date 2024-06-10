package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func Run() {
	Task := cron.New(cron.WithSeconds())

	// 每3秒执行一次 usdt转账日志
	_, err := Task.AddFunc("*/3 * * * * *", func() {
		GetUsdtlog()
	})

	if err != nil {
		fmt.Println("添加定时任务失败：", err)
		return
	}
	Task.Start()
}
