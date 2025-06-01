package seckillapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type IndexSecKillGoodsListRequest struct {
	Date string `form:"date" binding:"required"` // 秒杀日期，格式为 "2006-01-02-15"
}

type IndexSecKillGoodsListResponse struct {
	models.SecKillInfo
}

func (SecKillApi) IndexSecKillGoodsListView(c *gin.Context) {
	cr := middleware.GetBind[IndexSecKillGoodsListRequest](c)

	key := fmt.Sprintf("sec:goods:%s", cr.Date)
	var list = make([]IndexSecKillGoodsListResponse, 0)
	infoMap := global.Redis.HGetAll(context.Background(), key).Val()
	for _, v := range infoMap {
		var info models.SecKillInfo
		err := json.Unmarshal([]byte(v), &info)
		if err != nil {
			logrus.Warnf("秒杀商品信息Json解析失败: %v", err)
			continue
		}
		list = append(list, IndexSecKillGoodsListResponse{
			SecKillInfo: info,
		})
	}
	res.OkWithData(list, c)
}
