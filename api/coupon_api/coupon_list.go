package couponapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

func (CouponApi) CouponListView(c *gin.Context) {
	var cr = middleware.GetBind[models.PageInfo](c)

	list, count, _ := common.QueryList(models.CouponModel{}, common.QueryOption{
		PageInfo: cr,
	})
	res.OkWithList(list, count, c)
}
