package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderSendOutGoodsRequest struct {
	OrderGoodsID  uint   `json:"orderGoodsID" binding:"required"` // 订单商品ID
	WaybillNumber string `json:"waybillNumber"`                   // 运单号
	Message       string `json:"message"`                         // 备注信息
}

func (OrderApi) OrderSendOutGoodsView(c *gin.Context) {
	cr := middleware.GetBind[OrderSendOutGoodsRequest](c)

	// 检查订单商品是否存在
	var orderGoods models.OrderGoodsModel
	err := global.DB.Preload("OrderModel").
		Where("id = ?", cr.OrderGoodsID).
		Take(&orderGoods).Error
	if err != nil {
		res.FailWithMsg("订单不存在或无权限操作", c)
		return
	}

	// 更新订单状态为已发货
	if orderGoods.OrderModel.Status != 2 {
		res.FailWithMsg("订单状态不允许发货", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		//改订单商品
		err := tx.Model(&orderGoods).Updates(map[string]any{
			"status":         1,
			"waybill_number": cr.WaybillNumber,
		}).Error
		if err != nil {
			res.FailWithMsg("发货失败", c)
			return err
		}
		//查订单是不是已经改完
		var orderGoodsList []models.OrderGoodsModel
		tx.Find(&orderGoodsList, "order_id = ? AND status = ?", orderGoods.OrderID, 0)
		if len(orderGoodsList) == 0 {
			tx.Model(&orderGoods.OrderModel).Updates(map[string]any{
				"status": 3, // 更新订单状态为已发货
			})
		}

		if cr.Message != "" {
			// 添加备注信息
			err := tx.Create(&models.MessageModel{
				UserID:       orderGoods.UserID,
				OrderID:      orderGoods.OrderID,
				GoodsID:      orderGoods.GoodsID,
				OrderGoodsID: orderGoods.ID,
				MsgList:      []string{cr.Message},
				IsRead:       false,
			}).Error
			if err != nil {
				res.FailWithMsg("添加备注信息失败", c)
				return err
			}
		}
		return nil
	})
	if err != nil {
		res.FailWithMsg("发货失败", c)
		return
	}

	res.OkWithMsg("发货成功", c)
}
