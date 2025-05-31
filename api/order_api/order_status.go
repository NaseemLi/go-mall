package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderStatusRequest struct {
	No string `form:"no" binding:"required"`
}

type OrderStatusResponse struct {
	Status int8 `json:"status"`
}

func (OrderApi) OrderStatusView(c *gin.Context) {
	cr := middleware.GetBind[OrderStatusRequest](c)
	claims := middleware.GetAuth(c)

	var model models.OrderModel
	err := global.DB.Take(&model, "no = ? and user_id = ?", cr.No, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("订单不存在", c)
		return
	}

	data := OrderStatusResponse{
		Status: model.Status,
	}

	res.OkWithData(data, c)
}
