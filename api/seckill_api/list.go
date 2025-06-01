package seckillapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type ListRequest struct {
	models.PageInfo
}

type ListResponse struct {
	models.SecKillInfo
	CreateAt time.Time `json:"createAt"` // 创建时间
	ID       uint      `json:"id"`
}

func (SecKillApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)

	_list, count, _ := common.QueryList(models.SecKillModel{}, common.QueryOption{
		PageInfo: cr.PageInfo,
	})

	var list = make([]ListResponse, 0)

	for _, v := range _list {

		item := ListResponse{
			CreateAt: v.CreatedAt,
			ID:       v.ID,
		}

		val, err := global.Redis.HGet(context.Background(), v.Key(), fmt.Sprintf("%d", v.GoodsID)).Result()
		if err != nil {
			if err == redis.Nil {
				logrus.Warnf("Redis中未找到该商品字段: key=%s field=%d", v.Key(), v.GoodsID)
				continue
			}
			logrus.Errorf("Redis 访问出错: %v", err)
			continue
		}

		var info models.SecKillInfo
		err = json.Unmarshal([]byte(val), &info)
		if err != nil {
			logrus.Warnf("从 Redis 解析秒杀商品购买数量失败: %v", err)
			continue
		}

		item.SecKillInfo = info
		list = append(list, item)
	}

	res.OkWithList(list, count, c)
}
