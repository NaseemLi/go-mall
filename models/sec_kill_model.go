package models

import "time"

type SecKillModel struct {
	Model
	GoodsID       uint       `json:"goodsID"`                     //商品ID
	GoodsModel    GoodsModel `gorm:"foreignKey:GoodsID" json:"-"` //商品模型
	KillPrice     int        `json:"killPrice"`                   //秒杀价格
	KillInventory int        `json:"killInventory"`               //秒杀库存
	BuyNum        int        `json:"buyNum"`                      //购买数量
	StartTime     time.Time  `json:"startTime"`                   //秒杀开始时间
	EndTime       time.Time  `json:"endTime"`                     //秒杀结束时间
}
