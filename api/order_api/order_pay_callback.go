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
	err := global.DB.Preload("UserCouponList.UserCouponModel").Take(&order, "no = ?", cr.No).Error
	if err != nil {
		res.FailWithMsg("订单不存在", c)
		return
	}

	if order.Status != 1 {
		res.FailWithMsg("订单状态异常,请勿支付", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		//改订单状态
		tx.Model(&order).Where("no = ?", cr.No).Update("status", 2) // 2已付款/待发货
		if order.PzKey != "" {
			//订单支付完成,延长凭证
			pzUidKey := fmt.Sprintf("sec:pz_uid:%s", order.PzKey)
			val := global.Redis.Get(context.Background(), pzUidKey).Val()
			if val == "" {
				logrus.Warnf("[订单处理] 秒杀商品已经过期")
				return nil
			}

			var pzInfo redis_ser.PZinfo
			err = json.Unmarshal([]byte(val), &pzInfo)
			if err != nil {
				logrus.Warnf("[订单处理] 秒杀凭证信息解析失败: %s", order.PzKey)
				return nil
			}

			global.Redis.Expire(context.Background(), pzUidKey, 60*time.Second)
			global.Redis.Expire(context.Background(), pzInfo.PZKey, 60*time.Second)
		}

		//如果有购物车,清空购物车
		if len(order.CarIDList) > 0 {
			var carList []models.CarModel
			tx.Find(&carList, "id IN ?", order.CarIDList)
			if len(carList) > 0 {
				err := tx.Delete(&carList).Error
				if err != nil {
					return err
				}
			}
		}

		//如果使用了优惠卷,就修改优惠卷状态
		if len(order.UserCouponList) > 0 {
			var userCouponIDList []uint
			for _, v := range order.UserCouponList {
				userCouponIDList = append(userCouponIDList, v.UserCouponID)
			}
			var userCouponList []models.UserCouponModel
			tx.Find(&userCouponList, "id IN ?", userCouponIDList)
			if len(userCouponList) > 0 {
				for _, coupon := range userCouponList {
					err := tx.Delete(&coupon).Error
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		res.FailWithMsg("订单支付失败", c)
		return
	}

	res.OkWithMsg("订单支付成功", c)
}
