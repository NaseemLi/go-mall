package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderRevRequest struct {
	OrderID uint `json:"orderID" binding:"required"` // 订单ID
}

func (OrderApi) OrderRevGoodsView(c *gin.Context) {
	cr := middleware.GetBind[OrderRevRequest](c)
	claims := middleware.GetAuth(c)

	// 检查订单是否存在
	var model models.OrderModel
	err := global.DB.Where("id = ? AND user_id = ?", cr.OrderID, claims.UserID).Take(&model).Error
	if err != nil {
		res.FailWithMsg("订单不存在或无权限操作", c)
		return
	}

	err = global.DB.Model(&model).Update("status", 4).Error
	if err != nil {
		res.FailWithMsg("收货失败", c)
		return
	}

	res.OkWithMsg("收货成功", c)
}
