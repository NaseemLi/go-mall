package models

type UserCouponModel struct {
	Model
	UserID      uint        `json:"userID"` //用户ID
	UserModel   UserModel   `gorm:"foreignKey:UserID" json:"-"`
	CouponID    uint        `json:"couponID"` //优惠券ID
	CouponModel CouponModel `gorm:"foreignKey:CouponID" json:"-"`
	Status      int8        `json:"status"` //状态
}
