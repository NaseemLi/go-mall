package seckillapi

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/redis_ser"
	"fast_gin/utils/res"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SecKillRequest struct {
	Date    string `json:"date" binding:"required"`    // 秒杀日期
	GoodsID uint   `json:"goodsID" binding:"required"` // 商品ID
}

type SecKillResponse struct {
	Key string `json:"key"` // 购买凭证
}

var lock = sync.Mutex{}

func (SecKillApi) SecKillView(c *gin.Context) {
	cr := middleware.GetBind[SecKillRequest](c)
	claims := middleware.GetAuth(c)
	date, err := time.Parse("2006-01-02-15", cr.Date) // 校验日期格式
	if err != nil {
		res.FailWithMsg("日期格式错误", c)
		return
	}
	date = date.Local().Add(-8 * time.Hour)

	subHours := time.Since(date).Hours()
	if subHours < 0 {
		res.FailWithMsg("秒杀未开始", c)
		return
	}
	if subHours >= 1 {
		res.FailWithMsg("秒杀已结束", c)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	//如果需要分布式部署,这个地方需要改为分布式锁

	dateStr := date.Format("2006-01-02-15") // 示例输出: 2025-06-02-12
	key := fmt.Sprintf("sec:goods:%s", dateStr)
	field := fmt.Sprintf("%d", cr.GoodsID)

	result, err := global.Redis.HGet(context.Background(), key, field).Result()
	if err != nil {
		res.FailWithMsg("秒杀商品不存在", c)
		return
	}

	var info models.SecKillInfo
	err = json.Unmarshal([]byte(result), &info)
	if err != nil {
		res.FailWithMsg("秒杀商品信息解析失败", c)
		return
	}

	if info.BuyNum >= info.KillInventory {
		res.FailWithMsg("秒杀商品已售罄", c)
		return
	}

	pzKey := fmt.Sprintf("sec:pz:%s:%d:%d", dateStr, info.GoodsID, claims.UserID)
	_uid := global.Redis.Get(context.Background(), pzKey).Val()
	if _uid != "" {
		res.FailWithMsg("您已购买过该商品", c)
		return
	}

	//TODO:BuyNum不加问题
	info.BuyNum++
	byteData, _ := json.Marshal(info)
	global.Redis.HSet(context.Background(), key, field, byteData)

	uuid, _ := uuid.NewUUID()
	uid := uuid.String()

	global.Redis.Set(context.Background(), pzKey, uid, 15*time.Minute)
	pzInfoByteData, _ := json.Marshal(redis_ser.PZinfo{
		PZKey:     pzKey,
		GoodsInfo: info,
	})
	global.Redis.Set(context.Background(), "sec:pz_uid:"+uid, string(pzInfoByteData), 15*time.Minute)
	go func(pzKey string, key string, field string) {
		time.Sleep(15 * time.Minute)
		_uid := global.Redis.Get(context.Background(), pzKey).Val()
		if _uid != "" {
			//说明凭证已经延期了
			return
		}
		//已经过期了
		result, err := global.Redis.HGet(context.Background(), key, field).Result()
		if err != nil {
			return
		}
		var info models.SecKillInfo
		err = json.Unmarshal([]byte(result), &info)
		if err != nil {
			res.FailWithMsg("秒杀商品信息解析失败", c)
			return
		}

		info.BuyNum--
		byteData, _ := json.Marshal(info)
		global.Redis.HSet(context.Background(), key, field, string(byteData))
		logrus.Warnf("秒杀商品 %d:%d 购买凭证已过期", info.GoodsID, claims.UserID)

	}(pzKey, key, field)

	data := SecKillResponse{
		Key: uid,
	}

	res.OkWithData(data, c)
}
