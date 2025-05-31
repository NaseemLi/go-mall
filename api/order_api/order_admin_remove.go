package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type OrderAdminRemoveRequest struct {
}

func (OrderApi) OrderAdminRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var orderList []models.OrderModel
	global.DB.Find(&orderList, "id in ?", cr.IDList)
	if len(orderList) > 0 {
		global.DB.Delete(&orderList)
	}

	msg := fmt.Sprintf("订单删除成功,删除了 %d 个订单", len(orderList))

	res.OkWithMsg(msg, c)
}
