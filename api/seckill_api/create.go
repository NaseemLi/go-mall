package seckillapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateResquest struct {
	GoodsID       uint   `json:"goodsID" binding:"required"`       //商品ID
	KillPrice     int    `json:"killPrice" binding:"required"`     //秒杀价格
	KillInventory int    `json:"killInventory" binding:"required"` //秒杀库存
	StartTime     string `json:"startTime" binding:"required"`     //秒杀开始时间 2006-01-02 15:04:05Z07:00
}

func (SecKillApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateResquest](c)

	var goods models.GoodsModel
	err := global.DB.Take(&goods, "id = ?", cr.GoodsID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05Z07:00", cr.StartTime)
	if err != nil {
		res.FailWithMsg("开始时间格式错误", c)
		return
	}
	endTime := startTime.Add(time.Hour) // 默认结束时间为开始时间后1小时
	//1.同一个时间节点下，秒杀商品不能重复
	var model models.SecKillModel
	err = global.DB.Take(&model, "goods_id = ? AND start_time = ?", cr.GoodsID, startTime).Error
	if err == nil {
		res.FailWithMsg("同一时间节点下，秒杀商品不能重复", c)
		return
	}
	//2.开始时间是整小时，并且开始时间可以大于当前时间
	//3.库存数量不能大于商品本身的库存
	if goods.Inventory != nil {
		if cr.KillInventory > *goods.Inventory {
			res.FailWithMsg("秒杀库存不能大于商品本身的库存", c)
			return
		}
	}
	//4.价格不能大于等于商品本身的价格
	if cr.KillPrice >= goods.Price {
		res.FailWithMsg("秒杀价格不能大于等于商品本身的价格", c)
		return
	}

	//创建秒杀
	err = global.DB.Create(&models.SecKillModel{
		GoodsID:       cr.GoodsID,
		KillPrice:     cr.KillPrice,
		KillInventory: cr.KillInventory,
		StartTime:     startTime,
		EndTime:       endTime,
	}).Error
	if err != nil {
		res.FailWithMsg("创建秒杀失败", c)
		return
	}

	res.OkWithMsg("创建秒杀成功", c)
}
