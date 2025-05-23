package models

import (
	"fast_gin/models/ctype"
	"time"
)

type UserCouponModel struct {
	Model
	UserID      uint               `json:"userID"` //用户ID
	UserModel   UserModel          `gorm:"foreignKey:UserID" json:"-"`
	CouponID    uint               `json:"couponID"` //优惠券ID
	CouponModel CouponModel        `gorm:"foreignKey:CouponID" json:"-"`
	Status      ctype.CouponStatus `json:"status"`  //状态
	EndTime     time.Time          `json:"endTime"` //过期时间
}
