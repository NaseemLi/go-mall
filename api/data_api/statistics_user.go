package dataapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type StatisticsUserResponse struct {
	CarNum          int64 `json:"carNum"`          // 购物车商品数量
	MsgNum          int64 `json:"msgNum"`          // 消息数量
	Obligation      int64 `json:"obligation"`      // 待付款订单数量
	PendingShipment int64 `json:"pendingShipment"` // 待发货订单数量
	PendingPut      int64 `json:"pendingPut"`      // 待收货订单数量
	PendingComment  int64 `json:"pendingComment"`  // 待评价订单数量
}

func (DataApi) StatisticsUserView(c *gin.Context) {
	claims := middleware.GetAuth(c)
	var data StatisticsUserResponse

	global.DB.Model(models.CarModel{}).Where("user_id = ?", claims.UserID).Count(&data.CarNum)
	global.DB.Model(models.MessageModel{}).Where("user_id = ? and is_read = ?", claims.UserID, false).Count(&data.MsgNum)

	var orderList []models.OrderModel
	global.DB.Where("user_id = ?", claims.UserID).Find(&orderList)
	for _, v := range orderList {
		switch v.Status {
		case 1:
			data.Obligation++
		case 2:
			data.PendingShipment++
		case 3:
			data.PendingPut++
		case 4:
			data.PendingComment++
		}
	}

	res.OkWithData(data, c)
}
