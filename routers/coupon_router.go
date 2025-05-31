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

	// 优惠券管理（管理员）
	{
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

	// 优惠券展示（面向所有用户，无需鉴权）
	{
		g.GET("coupon/acceptable",
			middleware.BindQueryMiddleware[models.PageInfo],
			app.CouponUserAcceptableListView)
	}

	// 用户领券 & 用户券列表
	{
		g.POST("coupon/receive",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[couponapi.UserReceiveCouponRequest],
			app.UserReceiveCouponView)
		g.GET("coupon/user",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[couponapi.UserCouponListRequest],
			app.UserCouponListView)
	}
}
