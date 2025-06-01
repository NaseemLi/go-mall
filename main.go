package main

import (
	"fast_gin/core"
	"fast_gin/flags"
	"fast_gin/global"
	"fast_gin/routers"
	"fast_gin/service/cron_ser"
	redisdelay "fast_gin/service/redis_ser/redis_delay"
)

func main() {
	core.InitLogger()
	flags.Parse()
	global.Config = core.ReadConfig()
	global.DB = core.InitGorm()
	global.Redis = core.InitRedis()
	core.InitGorse()

	flags.Run()

	// 开启延时队列
	go redisdelay.PollOrderDelay()

	// 定时任务
	cron_ser.CronInit()

	routers.Run()
}
