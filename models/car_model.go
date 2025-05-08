package models

type CarModel struct {
	Model
	UserID     uint       `json:"userID"`                      //用户ID
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`  //用户信息
	GoodsID    uint       `json:"goodsID"`                     //商品ID
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"` //商品信息
	Price      int        `json:"price"`                       //价格单位:分
	Num        int        `json:"num"`                         //数量
}
