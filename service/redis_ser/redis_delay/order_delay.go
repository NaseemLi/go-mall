package redisdelay

import (
	"context"
	"fast_gin/global"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const queue = "delay_order_queue"

func AddOrderDelay(no string) {
	// 添加订单到延时队列
	global.Redis.ZAdd(context.Background(), queue, redis.Z{
		Member: no,
		Score:  float64(time.Now().Add(2 * time.Second).Unix()),
	})
}

func PollOrderDelay() {
	ctx := context.Background()
	for {
		// 获取当前时间之前的所有任务
		val, err := global.Redis.ZRangeByScore(ctx, queue, &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%d", time.Now().Unix()),
		}).Result()

		if err != nil {
			logrus.Errorf("查询任务失败: %v", err)
			return
		}

		for _, no := range val {
			logrus.Infof("处理订单: %s", no)
			global.Redis.ZRem(ctx, queue, no) // 从队列中移除已处理的订单
		}

		time.Sleep(1 * time.Second)
	}
}
