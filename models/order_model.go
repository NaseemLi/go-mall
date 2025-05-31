package models

import "time"

// 订单状态: 1代付款 2已付款/待发货 3已发货/待收货 4已收货/待评价 5已完成 6已取消 7已超时

type OrderModel struct {
	Model
	UserID         uint               `json:"userID"`
	UserModel      UserModel          `gorm:"foreignKey:UserID" json:"-"`
	AddrID         uint               `json:"addrID"`                                         //地址ID
	OrderGoodsList []OrderGoodsModel  `gorm:"foreignKey:OrderID" json:"-"`                    //订单商品列表
	UserCouponList []OrderCouponModel `gorm:"foreignKey:OrderID" json:"-"`                    //订单优惠券列表
	PayType        int8               `json:"payType"`                                        //支付方式
	Price          int                `json:"price"`                                          //价格单位:分
	Coupon         int                `json:"coupon"`                                         //优惠券
	No             string             `json:"no"`                                             //订单号
	Status         int8               `json:"status"`                                         //订单状态
	PayTime        time.Time          `json:"payTime"`                                        //支付时间
	PayUrl         string             `json:"payUrl"`                                         //支付链接
	CarIDList      []uint             `gorm:"type:longtext;serializer:json" json:"carIDList"` //购物车ID列表
	WaybillNumber  string             `gorm:"size:32" json:"waybillNumber"`                   //运单号
}

type OrderGoodsModel struct {
	Model
	UserID     uint       `json:"userID"`
	UserModel  UserModel  `gorm:"foreignKey:UserID" json:"-"`
	OrderID    uint       `json:"orderID"`
	OrderModel OrderModel `gorm:"foreignKey:OrderID" json:"-"`
	GoodsID    uint       `json:"goodsID"`
	GoodsModel GoodsModel `gorm:"foreignKey:GoodsID" json:"-"`
	Price      int        `json:"price"` //价格单位:分
	Num        int        `json:"num"`
	Note       string     `json:"note"` //商品备注
}

type OrderCouponModel struct {
	Model
	OrderID         uint            `json:"orderID"`
	OrderModel      OrderModel      `gorm:"foreignKey:OrderID" json:"-"`
	UserID          uint            `json:"userID"`
	UserModel       UserModel       `gorm:"foreignKey:UserID" json:"-"`
	UserCouponID    uint            `json:"couponID"`
	UserCouponModel UserCouponModel `gorm:"foreignKey:UserCouponID" json:"-"`
}
