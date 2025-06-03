package dataapi

import (
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsSystemResponse struct {
	UserNum         int64 `json:"userNum"`         // 用户数量
	GoodsNum        int64 `json:"goodsNum"`        // 商品数量
	SecKillNum      int64 `json:"secKillNum"`      // 秒杀活动数量
	SussessOrderNum int64 `json:"sussessOrderNum"` // 成功订单数量
	NewLoginCount   int64 `json:"newLoginCount"`   // 新登录用户数量
	Obligation      int64 `json:"obligation"`      // 待付款订单数量
	PendingShipment int64 `json:"pendingShipment"` // 待发货订单数量
	PendingPut      int64 `json:"pendingPut"`      // 待收货订单数量
	PendingComment  int64 `json:"pendingComment"`  // 待评价订单数量
}

func (DataApi) StatisticsSystemView(c *gin.Context) {
	var data StatisticsSystemResponse
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, loc)
	tomorrow := today.Add(24 * time.Hour)

	global.DB.Model(models.UserModel{}).Count(&data.UserNum)
	global.DB.Model(models.GoodsModel{}).Where("status = ?", ctype.GoodsStatusTop).Count(&data.GoodsNum)
	global.DB.Model(models.SecKillModel{}).Count(&data.SecKillNum)
	global.DB.Model(models.OrderModel{}).Where("status not in ?", []int8{1, 6, 7}).Count(&data.SussessOrderNum)
	global.DB.Debug().Model(models.UserLoginModel{}).Where("created_at >= ? AND created_at < ?", today, tomorrow).Count(&data.NewLoginCount)

	var orderList []models.OrderModel
	global.DB.Find(&orderList)
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
