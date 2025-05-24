package routers

import (
	"fast_gin/api"
	collectapi "fast_gin/api/user_center_api/collect_api"
	lookapi "fast_gin/api/user_center_api/look_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func UserCenterRouter(g *gin.RouterGroup) {
	app := api.App.UserCenterApi

	// 浏览记录
	{
		g.POST("user_center/look",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[lookapi.LookGoodsRequest],
			app.LookApi.LookGoodsView)
		g.GET("user_center/look",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[models.PageInfo],
			app.LookApi.LookGoodsListView)
		g.DELETE("user_center/look",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.LookApi.LookRemoveView)
	}

	// 收藏商品
	{
		g.POST("user_center/collect",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[collectapi.CollectGoodsRequest],
			app.CollectApi.CollectGoodsView)
		g.GET("user_center/collect",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[models.PageInfo],
			app.CollectApi.CollectGoodsListView)
		g.DELETE("user_center/collect",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.CollectApi.CollectRemoveView)
	}
}
