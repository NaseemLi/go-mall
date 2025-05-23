package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type UserReceiveCouponRequest struct {
	CouponID uint `json:"couponID"`
}

func (CouponApi) UserReceiveCouponView(c *gin.Context) {
	cr := middleware.GetBind[UserReceiveCouponRequest](c)

	user, err := middleware.GetUser(c)
	if err != nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}

	var coupon models.CouponModel
	err = global.DB.Take(&coupon, cr.CouponID).Error
	if err != nil {
		res.FailWithMsg("优惠卷不存在", c)
		return
	}

	if coupon.Receive == coupon.Num {
		res.FailWithMsg("优惠卷已经被领取完了", c)
	}
	var userCoupon models.UserCouponModel
	err = global.DB.Take(&userCoupon, "user_id = ? and coupon_id = ?", user.ID, cr.CouponID).Error
	if err == nil {
		res.FailWithMsg("你已经领取过该优惠卷了", c)
		return
	}

	global.DB.Create(&models.UserCouponModel{
		UserID:   user.ID,
		CouponID: cr.CouponID,
		EndTime:  time.Now().Add(time.Duration(coupon.Validity) * time.Hour),
	})
	//增加领取数
	global.DB.Model(&coupon).Where("id = ?", cr.CouponID).Update("receive", coupon.Receive+1)

	res.OkWithMsg("领取成功", c)
}
