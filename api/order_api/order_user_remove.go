package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type OrderUserRemoveRequest struct {
}

func (OrderApi) OrderUserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)
	claims := middleware.GetAuth(c)

	var orderList []models.OrderModel
	global.DB.Find(&orderList, "user_id = ? and id in ?", claims.UserID, cr.IDList)
	if len(orderList) > 0 {
		global.DB.Delete(&orderList)
	}

	msg := fmt.Sprintf("订单删除成功,删除了 %d 个订单", len(orderList))

	res.OkWithMsg(msg, c)
}
