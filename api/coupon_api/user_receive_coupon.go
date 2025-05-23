package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserReceiveCouponRequest struct {
	CouponID uint `json:"couponID" binding:"required"`
}

var mutex sync.Mutex

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
	if coupon.Type == ctype.CouponNewGoodsType {
		res.FailWithMsg("新商品优惠卷不可领取", c)
		return
	}

	//加锁解决并发问题
	mutex.Lock()
	defer mutex.Unlock()

	if coupon.Receive == coupon.Num {
		res.FailWithMsg("优惠卷已经被领取完了", c)
	}
	var userCoupon models.UserCouponModel
	err = global.DB.Take(&userCoupon, "user_id = ? and coupon_id = ?", user.ID, cr.CouponID).Error
	if err == nil {
		res.FailWithMsg("你已经领取过该优惠卷了", c)
		return
	}

	//原子性问题
	err = global.DB.Transaction(func(tx *gorm.DB) error {

		err = tx.Create(&models.UserCouponModel{
			UserID:   user.ID,
			CouponID: cr.CouponID,
			Status:   ctype.CouponStatusNotUsed,
			EndTime:  time.Now().Add(time.Duration(coupon.Validity) * time.Hour),
		}).Error
		if err != nil {
			return err
		}
		//增加领取数
		err = tx.Model(&coupon).Where("id = ?", cr.CouponID).Update("receive", gorm.Expr("receive + 1")).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		res.FailWithMsg("领取失败", c)
		return
	}

	res.OkWithMsg("领取成功", c)
}
