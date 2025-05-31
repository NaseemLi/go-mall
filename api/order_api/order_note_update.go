package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderNoteUpdateRequest struct {
	OrderGoodsID uint   `json:"orderGoodsID" binding:"required"`
	Note         string `json:"note" binding:"required"`
}

func (OrderApi) OrderNoteUpdateView(c *gin.Context) {
	cr := middleware.GetBind[OrderNoteUpdateRequest](c)
	claims := middleware.GetAuth(c)

	var order models.OrderGoodsModel
	if err := global.DB.Preload("OrderModel").Take(&order, "user_id = ? and id = ?",
		claims.UserID, cr.OrderGoodsID).Error; err != nil {
		res.FailWithMsg("订单商品不存在", c)
		return
	}

	if !(order.OrderModel.Status == 1 || order.OrderModel.Status == 2) {
		res.FailWithMsg("订单备注不可修改", c)
		return
	}

	if err := global.DB.Model(&order).Update("note", cr.Note).Error; err != nil {
		res.FailWithMsg("更新订单备注失败", c)
		return
	}

	res.OkWithMsg("更新商品订单备注成功", c)
}
