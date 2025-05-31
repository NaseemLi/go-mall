package routers

import (
	"fast_gin/api"
	carapi "fast_gin/api/car_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func CarRouter(g *gin.RouterGroup) {
	app := api.App.CarApi

	// 添加、删除购物车商品
	{
		g.POST("car",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[carapi.CarCreateRequest],
			app.CarCreateView)
		g.DELETE("car",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.CarRemoveView)
	}

	// 购物车商品数量更新
	{
		g.PUT("car/num",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[carapi.CarNumUpdateRequest],
			app.CarNumUpdateView)
	}

	// 查询购物车列表
	{
		g.POST("car/list",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[models.PageInfo],
			app.CarListView)
	}

	// 将购物车商品转移至收藏
	{
		g.POST("car/collect",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.CarToCollectView)
	}
}
