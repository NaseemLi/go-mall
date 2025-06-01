package models

import (
	"context"
	"encoding/json"
	"fast_gin/global"
	"fmt"
	"time"

	"gorm.io/gorm"
)

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

type SecKillInfo struct {
	GoodsID       uint      `json:"goodsID"`       //商品ID
	Title         string    `json:"title"`         //商品标题
	Cover         string    `json:"cover"`         //商品封面图
	Price         int       `json:"price"`         //商品价格
	KillPrice     int       `json:"killPrice"`     //秒杀价格
	KillInventory int       `json:"killInventory"` //秒杀库存
	BuyNum        int       `json:"buyNum"`        //购买数量
	StartTime     time.Time `json:"startTime"`     //秒杀开始时间
	EndTime       time.Time `json:"endTime"`       //秒杀结束时间
}

func (s *SecKillModel) GetSecKillInfo() SecKillInfo {
	return SecKillInfo{
		GoodsID:       s.GoodsID,
		Title:         s.GoodsModel.Title,
		Cover:         s.GoodsModel.GetCover(),
		Price:         s.GoodsModel.Price,
		KillPrice:     s.KillPrice,
		KillInventory: s.KillInventory,
		BuyNum:        s.BuyNum,
		StartTime:     s.StartTime,
		EndTime:       s.EndTime,
	}
}

func (s *SecKillModel) Key() string {
	return fmt.Sprintf("sec:goods:%s", s.StartTime.Format("2006-01-02-15"))
}

func (s *SecKillModel) AfterCreate(tx *gorm.DB) error {
	key := s.Key()
	field := fmt.Sprintf("%d", s.GoodsID)
	byteData, _ := json.Marshal(s.GetSecKillInfo())

	global.Redis.HSetNX(context.Background(), key, field, string(byteData))
	global.Redis.ExpireAt(context.Background(), key, s.StartTime.Add(time.Hour))

	return nil
}

func (s *SecKillModel) AfterDelete(tx *gorm.DB) error {
	key := s.Key()
	field := fmt.Sprintf("%d", s.GoodsID)

	global.Redis.HDel(context.Background(), key, field)

	return nil
}
