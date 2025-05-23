package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CouponRemoveRequest struct {
}

func (CouponApi) CouponRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var couponList []models.CouponModel
	global.DB.Find(&couponList, "id in ?", cr.IDList)
	if len(couponList) > 0 {
		global.DB.Delete(&couponList)
	}

	msg := fmt.Sprintf("优惠券删除成功,删除了 %d 个优惠券", len(couponList))

	res.OkWithMsg(msg, c)
}
