package seckillapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	payser "fast_gin/service/pay_ser"
	"fast_gin/service/redis_ser"
	redisdelay "fast_gin/service/redis_ser/redis_delay"
	"fast_gin/utils/random"
	"fast_gin/utils/res"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SecKillOrderRequest struct {
	Key     string `json:"key" binding:"required"`
	AddrID  uint   `json:"addrID" binding:"required"`
	PayType int8   `json:"payType" binding:"required"`
	Note    string `json:"note"` // 订单备注
}

type SecKillOrderResponse struct {
	No     string `json:"no"`     // 订单号
	PayUrl string `json:"payUrl"` // 支付链接
	Price  int    `json:"price"`  // 订单价格
}

var orderLock = sync.Mutex{}

func (SecKillApi) SecKillOrderView(c *gin.Context) {
	cr := middleware.GetBind[SecKillOrderRequest](c)
	claims := middleware.GetAuth(c)
	var order models.OrderModel
	err := global.DB.Take(&order, "user_id = ? and pz_key = ?", claims.UserID, cr.Key).Error
	if err == nil {
		res.FailWithMsg("请勿重复下单", c)
		return
	}

	val, _ := global.Redis.Get(context.Background(), "sec:pz_uid:"+cr.Key).Result()
	if val == "" {
		res.FailWithMsg("购买凭证无效", c)
		return
	}

	var info redis_ser.PZinfo
	err = json.Unmarshal([]byte(val), &info)
	if err != nil {
		res.FailWithMsg("秒杀商品信息Json解析失败", c)
		return
	}

	orderLock.Lock()
	defer orderLock.Unlock()

	no := random.GenerateOrderNumber()
	price := info.GoodsInfo.KillPrice
	payUrl, err := payser.Pay(cr.PayType, no, price)
	if err != nil || payUrl == "" {
		logrus.Errorf("支付服务调用失败: %v", err)
		res.FailWithMsg("支付服务调用失败", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		order := models.OrderModel{
			No:      no,
			UserID:  claims.UserID,
			AddrID:  cr.AddrID,
			Price:   price,
			Status:  1,
			PayType: cr.PayType,
			Coupon:  info.GoodsInfo.Price - info.GoodsInfo.KillPrice,
			PayTime: time.Now(),
			PayUrl:  payUrl,
			PzKey:   info.PZKey,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var goodsOrderList = models.OrderGoodsModel{
			OrderID: order.ID,
			GoodsID: info.GoodsInfo.GoodsID,
			UserID:  claims.UserID,
			Price:   price,
			Num:     1,
			Note:    cr.Note,
		}

		if err := tx.Create(&goodsOrderList).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Errorf("创建订单失败: %v", err)
		res.FailWithMsg("创建订单失败", c)
		return
	}

	data := SecKillOrderResponse{
		No:     no,
		PayUrl: payUrl,
		Price:  price,
	}

	//延时队列
	redisdelay.AddOrderDelay(data.No)
	//延长凭证时间
	//TODO:延时问题
	global.Redis.Expire(context.Background(), info.PZKey, 20*time.Minute)
	global.Redis.Expire(context.Background(), "sec:pz_uid:"+cr.Key, 20*time.Minute)

	res.OkWithData(data, c)
}
