package models

type MessageModel struct {
	Model
	UserID       uint       `json:"userID"` //用户ID
	UserModel    UserModel  `gorm:"foreignKey:UserID" json:"-"`
	OrderID      uint       `json:"orderID"` //订单ID
	OrderModel   OrderModel `gorm:"foreignKey:OrderID" json:"-"`
	MsgList      []string   `gorm:"type:longtext;serializer:json" json:"msgList"`
	IsRead       bool       `json:"isRead"`  //是否已读
	GoodsID      uint       `json:"goodsID"` //商品ID
	GoodsModel   GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
	OrderGoodsID uint       `json:"orderGoodsID"` //订单商品ID
}
