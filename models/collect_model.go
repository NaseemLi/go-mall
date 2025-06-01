package models

import (
	"context"
	"fast_gin/global"
	"fmt"
	"time"

	"github.com/zhenghaoz/gorse/client"
	"gorm.io/gorm"
)

type CollectModel struct {
	Model
	UserID     uint       `json:"userID"` //用户ID
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`
	GoodsID    uint       `json:"goodsID"`    //商品ID
	GoodsTitle string     `json:"goodsTitle"` //商品标题
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
}

func (c *CollectModel) AfterCreate(tx *gorm.DB) error {
	if global.Gorse == nil {
		return nil
	}

	global.Gorse.InsertFeedback(context.Background(), []client.Feedback{
		{
			FeedbackType: "like",
			UserId:       fmt.Sprintf("%d", c.UserID),
			ItemId:       fmt.Sprintf("%d", c.GoodsID),
			Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		},
	})

	return nil
}

func (c *CollectModel) BeforeDelete(tx *gorm.DB) error {
	if global.Gorse == nil {
		return nil
	}

	global.Gorse.InsertFeedback(context.Background(), []client.Feedback{
		{
			FeedbackType: "read",
			UserId:       fmt.Sprintf("%d", c.UserID),
			ItemId:       fmt.Sprintf("%d", c.GoodsID),
			Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		},
	})

	return nil
}
