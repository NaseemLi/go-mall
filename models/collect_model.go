package models

type CollectModel struct {
	Model
	UserID     uint       `json:"userID"` //用户ID
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`
	GoodsID    uint       `json:"goodsID"` //商品ID
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
}
