package models

import "time"

type CouponModel struct {
	Model
	Title        string     `json:"title"`        //优惠券名称
	Type         int8       `json:"type"`         //优惠券类型
	CouponPrice  int        `json:"couponPrice"`  //优惠券金额
	Threshold    int        `json:"threshold"`    //使用门槛
	StartTime    *time.Time `json:"startTime"`    //开始时间
	EndTime      *time.Time `json:"endTime"`      //结束时间
	Num          int        `json:"num"`          //优惠卷数量
	Push         int8       `json:"push"`         //推广方式
	GoodsID      *uint      `json:"goodsID"`      //关联的商品
	GoodCategory *string    `json:"goodCategory"` //关联的商品分类
}
