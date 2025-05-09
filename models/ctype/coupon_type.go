package ctype

type CouponType int8

const (
	CouponGeneralType   CouponType = iota + 1 //通用
	CouponConditionType                       //满减
	CouponCategoryType                        //门类
	CouponGoodsType                           //商品
)
