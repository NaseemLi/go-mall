package models

import (
	"fast_gin/models/ctype"
)

type CouponModel struct {
	Model
	Title       string           `json:"title"`       //优惠券名称
	Type        ctype.CouponType `json:"type"`        //优惠券类型
	CouponPrice int              `json:"couponPrice"` //优惠券金额
	Threshold   int              `json:"threshold"`   //使用门槛
	Validity    int              `json:"validity"`    //有效期 单位小时
	Num         int              `json:"num"`         //优惠卷数量
	Receive     int              `json:"receive"`     //领取数量
	GoodsID     *uint            `json:"goodsID"`     //关联的商品
	Festival    *string          `json:"festival"`    //节日活动
}
