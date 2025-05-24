package models

type LookGoodsModel struct {
	Model
	UserID     uint       `json:"userID"` //用户ID
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`
	GoodsID    uint       `json:"goodsID"` //商品ID
	GoodsTitle string     `json:"goodsTitle"`
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
}
