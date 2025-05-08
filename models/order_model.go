package models

type OrderModel struct {
	Model
	UserID         uint              `json:"userID"`
	UserModel      UserModel         `gorm:"foreignKey:UserID" json:"-"`
	AddrID         uint              `json:"addrID"`                      //地址ID
	OrderGoodsList []OrderGoodsModel `gorm:"foreignKey:OrderID" json:"-"` //订单商品列表
	PayType        int8              `json:"payType"`                     //支付方式
	Price          int               `json:"price"`                       //价格单位:分
	Coupon         int               `json:"coupon"`                      //优惠券
	No             string            `json:"no"`                          //订单号
	Status         int8              `json:"status"`                      //订单状态
}

type OrderGoodsModel struct {
	Model
	OrderID    uint       `json:"orderID"`
	OrderModel OrderModel `gorm:"foreignKey:OrderID" json:"-"`
	GoodsID    uint       `json:"goodsID"`
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
	Price      uint       `json:"price"` //价格单位:分
	Num        int        `json:"num"`
}
