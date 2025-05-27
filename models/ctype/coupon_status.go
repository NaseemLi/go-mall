package ctype

type CouponStatus int8

const (
	CouponStatusNotUsed CouponStatus = iota + 1 //未使用
	CouponStatusUsed                            //已使用
	CouponStatusExpired                         //已过期
	CouponStatusLocked                          //已锁定
)
