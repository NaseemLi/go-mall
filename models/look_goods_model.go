package models

import (
	"context"
	"fast_gin/global"
	"fmt"
	"time"

	"github.com/zhenghaoz/gorse/client"
	"gorm.io/gorm"
)

type LookGoodsModel struct {
	Model
	UserID     uint       `json:"userID"` //用户ID
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`
	GoodsID    uint       `json:"goodsID"` //商品ID
	GoodsTitle string     `json:"goodsTitle"`
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
}

func (l *LookGoodsModel) AfterCreate(tx *gorm.DB) error {
	if global.Gorse == nil {
		return nil
	}

	global.Gorse.InsertFeedback(context.Background(), []client.Feedback{
		{
			FeedbackType: "read",
			UserId:       fmt.Sprintf("%d", l.UserID),
			ItemId:       fmt.Sprintf("%d", l.GoodsID),
			Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		},
	})

	return nil
}
