package routers

import (
	"fast_gin/api"
	couponapi "fast_gin/api/coupon_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func CouponRouter(g *gin.RouterGroup) {
	app := api.App.CouponApi
	g.POST("coupon",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[couponapi.CouponCreateRequest],
		app.CouponCreateView)
}
