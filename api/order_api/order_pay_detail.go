package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderPayDetailRequest struct {
	No string `form:"no" binding:"required"` // 订单号
}

type OrderPayDetailResponse struct {
	No    string `json:"no"`    // 订单号
	Price int    `json:"price"` // 订单总价
}

func (OrderApi) OrderPayDetailView(c *gin.Context) {
	cr := middleware.GetBind[OrderPayDetailRequest](c)

	var order models.OrderModel
	err := global.DB.
		Take(&order, "no = ?", cr.No).Error
	if err != nil {
		res.FailWithMsg("订单不存在", c)
		return
	}

	if order.Status != 1 {
		res.FailWithMsg("订单状态异常,请勿支付", c)
		return
	}

	data := OrderPayDetailResponse{
		No:    cr.No,
		Price: order.Price,
	}

	res.OkWithData(data, c)

}
