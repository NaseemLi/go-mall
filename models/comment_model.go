package models

type CommentModel struct {
	Model
	UserID          uint            `json:"userID"`  //用户ID
	OrderID         uint            `json:"orderID"` // 订单ID
	GoodsID         uint            `json:"goodsID"` // 商品ID
	GoodsModel      GoodsModel      `gorm:"foreignKey:GoodsID" json:"-"`
	UserModel       UserModel       `gorm:"foreignKey:UserID" json:"-"`
	OrderGoodsID    uint            `json:"orderGoodsID"` //商品评价的是这个人下单买的商品
	OrderGoodsModel OrderGoodsModel `gorm:"foreignKey:OrderGoodsID" json:"-"`
	Level           int8            `json:"level"`                                        //满意度
	Content         string          `gorm:"type:longtext;serializer:json" json:"content"` //评价内
	Images          []string        `gorm:"type:longtext;serializer:json" json:"images"`  //评价图片
}
