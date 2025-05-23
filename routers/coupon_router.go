package routers

import (
	"fast_gin/api"
	couponapi "fast_gin/api/coupon_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func CouponRouter(g *gin.RouterGroup) {
	app := api.App.CouponApi
	g.POST("coupon",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[couponapi.CouponCreateRequest],
		app.CouponCreateView)
	g.GET("coupon",
		middleware.AdminMiddleware,
		middleware.BindQueryMiddleware[models.PageInfo],
		app.CouponListView)
	g.DELETE("coupon",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[models.IDListRequest],
		app.CouponRemoveView)
}
