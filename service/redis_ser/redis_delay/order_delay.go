package redisdelay

import (
	"context"
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const waitTime = 5 // 延时处理时间，单位为秒
const queue = "delay_order_queue"

// 添加订单到延时队列
func AddOrderDelay(no string) {
	if global.Redis == nil {
		logrus.Warnf("[延时队列] Redis 未初始化，无法添加订单: %s", no)
		return
	}

	global.Redis.ZAdd(context.Background(), queue, redis.Z{
		Member: no,
		Score:  float64(time.Now().Add(waitTime * time.Second).Unix()),
	})
	logrus.Infof("[延时队列] 已添加订单: %s，延迟 %d 秒执行", no, waitTime)
}

// 启动延时订单轮询处理
func PollOrderDelay() {
	if global.Redis == nil {
		logrus.Warnf("[轮询任务] Redis 未初始化，无法启动延时处理")
		return
	}

	ctx := context.Background()
	for {
		val, err := global.Redis.ZRangeByScore(ctx, queue, &redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprintf("%d", time.Now().Unix()),
		}).Result()

		if err != nil {
			logrus.Errorf("[轮询任务] 查询 Redis 队列失败: %v", err)
			return
		}

		for _, no := range val {
			logrus.Infof("[轮询任务] 检测到超时订单: %s", no)
			OrderDelay(no)
			global.Redis.ZRem(ctx, queue, no)
			logrus.Infof("[轮询任务] 已移除处理完成订单: %s", no)
		}

		time.Sleep(1 * time.Second)
	}
}

// 执行单个订单超时处理
func OrderDelay(no string) {
	logrus.Infof("[订单处理] 开始处理订单: %s", no)

	var model models.OrderModel
	err := global.DB.Take(&model, "no = ?", no).Error
	if err != nil {
		logrus.Warnf("[订单处理] 查询失败，订单不存在: %s", no)
		return
	}

	// 判断订单状态是否为待支付
	if model.Status != 1 {
		logrus.Infof("[订单处理] 订单状态非待支付，跳过处理: %s", no)
		return
	}

	// 执行订单超时事务处理
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		tx.Model(&model).Update("status", 7)
		logrus.Infof("[订单处理] 订单标记为超时: %s", no)

		// 释放购物车
		var carList []models.CarModel
		err := global.DB.Where("id in ?", model.CarIDList).Find(&carList).Error
		if err != nil {
			logrus.Errorf("[订单处理] 查询购物车失败: %s", no)
			return err
		}
		if len(carList) > 0 {
			tx.Model(&carList).Update("status", 0)
			logrus.Infof("[订单处理] 已释放购物车: %s", no)
		}

		// 归还优惠券
		var userCoupon []models.UserCouponModel
		var userCouponIDList []uint
		for _, v := range model.UserCouponList {
			userCouponIDList = append(userCouponIDList, v.UserCouponID)
		}
		err = global.DB.Where("id in ?", userCouponIDList).Find(&userCoupon).Error
		if err != nil {
			logrus.Errorf("[订单处理] 查询优惠券失败: %s", no)
			return err
		}
		if len(userCoupon) > 0 {
			tx.Model(&userCoupon).Update("status", ctype.CouponStatusNotUsed)
			logrus.Infof("[订单处理] 已归还优惠券: %s", no)
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("[订单处理] 订单处理失败: %s, err: %v", no, err)
		return
	}

	logrus.Infof("[订单处理] 超时订单处理完成: %s", no)
}
