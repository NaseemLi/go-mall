package models

type MessageModel struct {
	Model
	UserID          uint            `json:"userID"` //用户ID
	UserModel       UserModel       `gorm:"foreignKey:UserID" json:"-"`
	OrderGoodsID    uint            `json:"orderGoodsID"` //订单商品ID
	OrderGoodsModel OrderGoodsModel `gorm:"foreignKey:OrderGoodsID" json:"-"`
	MsgList         []string        `gorm:"type:longtext;serializer:json" json:"msgList"`
	IsRead          bool            `json:"isRead"` //是否已读
}
