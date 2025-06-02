package orderapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/redis_ser"
	"fast_gin/utils/res"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderPayCallbackRequest struct {
	No string `json:"no" binding:"required"` // 订单号
}

func (OrderApi) OrderPayCallbackView(c *gin.Context) {
	cr := middleware.GetBind[OrderPayCallbackRequest](c)

	var order models.OrderModel
	err := global.DB.
		Select("*").
		Preload("UserCouponList.UserCouponModel").
		Take(&order, "no = ?", cr.No).Error
	if err != nil {
		res.FailWithMsg("订单不存在", c)
		return
	}

	logrus.WithFields(logrus.Fields{
		"order_no": order.No,
		"pz_key":   order.PzKey,
	}).Info("收到订单支付回调")

	if order.Status != 1 {
		res.FailWithMsg("订单状态异常,请勿支付", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态为已付款
		tx.Model(&order).Where("no = ?", cr.No).Update("status", 2)

		if order.PzKey != "" {
			uid := global.Redis.Get(context.Background(), order.PzKey).Val()
			if uid == "" {
				logrus.WithFields(logrus.Fields{
					"order_no": order.No,
					"pz_key":   order.PzKey,
				}).Warn("从 Redis 获取 UID 失败，凭证已过期")
				return nil
			}

			pzUidKey := fmt.Sprintf("sec:pz_uid:%s", uid)
			val := global.Redis.Get(context.Background(), pzUidKey).Val()
			if val == "" {
				logrus.WithFields(logrus.Fields{
					"uid":       uid,
					"pz_uidKey": pzUidKey,
				}).Warn("凭证 JSON 未找到，可能已过期")
				return nil
			}

			var pzInfo redis_ser.PZinfo
			if err := json.Unmarshal([]byte(val), &pzInfo); err != nil {
				logrus.WithFields(logrus.Fields{
					"uid":    uid,
					"rawval": val,
				}).Warn("凭证信息解析失败")
				return nil
			}

			ok1, err1 := global.Redis.Expire(context.Background(), pzUidKey, 60*time.Minute).Result()
			ok2, err2 := global.Redis.Expire(context.Background(), pzInfo.PZKey, 60*time.Minute).Result()

			logrus.WithFields(logrus.Fields{
				"pz_uid_key": pzUidKey,
				"success":    ok1,
				"err":        err1,
			}).Info("重置 TTL 成功（UID凭证）")

			logrus.WithFields(logrus.Fields{
				"pz_key":  pzInfo.PZKey,
				"success": ok2,
				"err":     err2,
			}).Info("重置 TTL 成功（用户凭证）")
		}

		// 清空购物车
		if len(order.CarIDList) > 0 {
			var carList []models.CarModel
			tx.Find(&carList, "id IN ?", order.CarIDList)
			if len(carList) > 0 {
				if err := tx.Delete(&carList).Error; err != nil {
					return err
				}
			}
		}

		// 删除已使用的优惠券
		if len(order.UserCouponList) > 0 {
			var userCouponIDList []uint
			for _, v := range order.UserCouponList {
				userCouponIDList = append(userCouponIDList, v.UserCouponID)
			}
			var userCouponList []models.UserCouponModel
			tx.Find(&userCouponList, "id IN ?", userCouponIDList)
			if len(userCouponList) > 0 {
				for _, coupon := range userCouponList {
					if err := tx.Delete(&coupon).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order_no": cr.No,
			"error":    err,
		}).Error("订单支付事务失败")
		res.FailWithMsg("订单支付失败", c)
		return
	}

	logrus.WithField("order_no", cr.No).Info("订单支付成功")
	res.OkWithMsg("订单支付成功", c)
}
