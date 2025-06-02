package redisdelay

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/redis_ser"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const waitTime = 900 // 延时处理时间，单位为秒 15min=900s
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
	logrus.Infof("[延时队列] 已添加订单: %s,延迟 %d 秒执行", no, waitTime)
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

var lock = &sync.Mutex{}

// 执行单个订单超时处理
func OrderDelay(no string) {
	logrus.WithField("order_no", no).Info("[订单处理] 开始处理订单")

	var model models.OrderModel
	err := global.DB.Preload("OrderGoodsList").Take(&model, "no = ?", no).Error
	if err != nil {
		logrus.WithField("order_no", no).Warn("[订单处理] 查询失败，订单不存在")
		return
	}

	if model.Status != 1 {
		logrus.WithFields(logrus.Fields{
			"order_no": no,
			"status":   model.Status,
		}).Info("[订单处理] 订单状态非待支付，跳过处理")
		return
	}

	lock.Lock()
	defer lock.Unlock()

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if model.PzKey != "" {
			uid := global.Redis.Get(context.Background(), model.PzKey).Val()
			if uid == "" {
				logrus.WithFields(logrus.Fields{
					"order_no": no,
					"pz_key":   model.PzKey,
				}).Warn("[订单处理] 凭证已过期，UID 不存在")
				return nil
			}

			pzUidKey := fmt.Sprintf("sec:pz_uid:%s", uid)
			val := global.Redis.Get(context.Background(), pzUidKey).Val()
			if val == "" {
				logrus.WithFields(logrus.Fields{
					"order_no":   no,
					"pz_uid_key": pzUidKey,
				}).Warn("[订单处理] 凭证信息已失效，跳过库存回退")
				return nil
			}

			var pzInfo redis_ser.PZinfo
			err := json.Unmarshal([]byte(val), &pzInfo)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"order_no":   no,
					"pz_uid_key": pzUidKey,
					"err":        err,
				}).Warn("[订单处理] 凭证信息解析失败")
				return nil
			}

			parts := strings.Split(pzInfo.PZKey, ":")
			if len(parts) < 4 {
				logrus.WithFields(logrus.Fields{
					"order_no": no,
					"pz_key":   pzInfo.PZKey,
				}).Warn("[订单处理] PZKey 格式非法")
				return nil
			}

			date := parts[2]
			field := parts[3]
			hashKey := fmt.Sprintf("sec:goods:%s", date)
			val = global.Redis.HGet(context.Background(), hashKey, field).Val()
			if val == "" {
				logrus.WithFields(logrus.Fields{
					"order_no": no,
					"hash_key": hashKey,
					"field":    field,
				}).Warn("[订单处理] Redis 秒杀商品信息缺失")
				return nil
			}

			var info models.SecKillInfo
			err = json.Unmarshal([]byte(val), &info)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"order_no": no,
					"err":      err,
				}).Warn("[订单处理] 秒杀商品信息解析失败")
				return nil
			}

			info.BuyNum--
			byteData, _ := json.Marshal(info)
			global.Redis.HSet(context.Background(), hashKey, field, string(byteData))
			global.Redis.Del(context.Background(), pzInfo.PZKey)
			global.Redis.Del(context.Background(), pzUidKey)

			logrus.WithFields(logrus.Fields{
				"order_no": no,
				"goods_id": info.GoodsID,
				"buy_num":  info.BuyNum,
				"hash_key": hashKey,
				"field":    field,
				"pz_uid":   uid,
				"pz_key":   pzInfo.PZKey,
			}).Info("[订单处理] 秒杀订单库存已回退并清除凭证")
			return nil
		}

		tx.Model(&model).Update("status", 7)
		logrus.WithField("order_no", no).Info("[订单处理] 非秒杀订单已标记为超时")

		//正产的订单 订单过期 将库存+1
		var goodsIDList []uint
		for _, v := range model.OrderGoodsList {
			goodsIDList = append(goodsIDList, v.GoodsID)
		}

		var goodsList []models.GoodsModel
		tx.Find(&goodsList, "id in ?", goodsIDList)
		for _, v := range goodsList {
			if v.Inventory == nil {
				continue
			}
			tx.Model(&v).Update("inventory", gorm.Expr("inventory + 1"))
		}

		var carList []models.CarModel
		if err := global.DB.Where("id in ?", model.CarIDList).Find(&carList).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"order_no": no,
				"err":      err,
			}).Error("[订单处理] 查询购物车失败")
			return err
		}
		if len(carList) > 0 {
			tx.Model(&carList).Update("status", 0)
			logrus.WithField("order_no", no).Info("[订单处理] 已释放购物车资源")
		}

		var userCoupon []models.UserCouponModel
		var userCouponIDList []uint
		for _, v := range model.UserCouponList {
			userCouponIDList = append(userCouponIDList, v.UserCouponID)
		}
		if err := global.DB.Where("id in ?", userCouponIDList).Find(&userCoupon).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"order_no": no,
				"err":      err,
			}).Error("[订单处理] 查询优惠券失败")
			return err
		}
		if len(userCoupon) > 0 {
			tx.Model(&userCoupon).Update("status", ctype.CouponStatusNotUsed)
			logrus.WithField("order_no", no).Info("[订单处理] 已归还用户优惠券")
		}

		return nil
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order_no": no,
			"err":      err,
		}).Error("[订单处理] 超时订单处理失败")
		return
	}

	logrus.WithField("order_no", no).Info("[订单处理] 超时订单处理完成")
}
