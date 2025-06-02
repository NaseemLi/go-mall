package seckillapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type SecKillDetailRequest struct {
	Key string `json:"key" binding:"required"` // 购买凭证
}

type SecKillDetailResponse struct {
	models.SecKillInfo
}

func (SecKillApi) SecKillDetailView(c *gin.Context) {
	cr := middleware.GetBind[SecKillDetailRequest](c)

	val, _ := global.Redis.Get(context.Background(), "sec:pz_uid:"+cr.Key).Result()
	if val == "" {
		res.FailWithMsg("购买凭证无效", c)
		return
	}

	var info PZinfo
	err := json.Unmarshal([]byte(val), &info)
	if err != nil {
		res.FailWithMsg("秒杀商品信息Json解析失败", c)
		return
	}

	data := SecKillDetailResponse{
		SecKillInfo: info.GoodsInfo,
	}

	res.OkWithData(data, c)
}
