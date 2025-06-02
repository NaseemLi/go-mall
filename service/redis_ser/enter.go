package redis_ser

import "fast_gin/models"

type PZinfo struct {
	PZKey     string             `json:"PZKey"`     // 购买凭证
	GoodsInfo models.SecKillInfo `json:"GoodsInfo"` // 商品信息
}
