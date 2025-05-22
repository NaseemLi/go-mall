package ctype

type CouponType int8

const (
	CouponFestivalType CouponType = iota + 1 //节日
	CouponNewUserType                        //新用户
	CouponNewGoodsType                       //新商品
	CouponGoodsType                          //商品
	CouponGeneralType                        //通用
)
