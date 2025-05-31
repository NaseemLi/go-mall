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
	OrderID       uint   `json:"orderID" binding:"required"` // 订单ID
	WaybillNumber string `json:"waybillNumber"`              // 运单号
	Message       string `json:"message"`                    // 备注信息
}

func (OrderApi) OrderSendOutGoodsView(c *gin.Context) {
	cr := middleware.GetBind[OrderSendOutGoodsRequest](c)

	// 检查订单是否存在
	var model models.OrderModel
	err := global.DB.Where("id = ?", cr.OrderID).Take(&model).Error
	if err != nil {
		res.FailWithMsg("订单不存在或无权限操作", c)
		return
	}

	// 更新订单状态为已发货
	if model.Status != 2 {
		res.FailWithMsg("订单状态不允许发货", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model).Updates(map[string]any{
			"status":         3,
			"waybill_number": cr.WaybillNumber,
		}).Error
		if err != nil {
			res.FailWithMsg("发货失败", c)
			return err
		}

		if cr.Message != "" {
			// 添加备注信息
			err := tx.Create(&models.MessageModel{
				UserID:  model.UserID,
				OrderID: model.ID,
				MsgList: []string{cr.Message},
				IsRead:  false,
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
