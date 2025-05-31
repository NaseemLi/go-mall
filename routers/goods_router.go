package routers

import (
	"fast_gin/api"
	goodsapi "fast_gin/api/goods_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func GoodsRouter(g *gin.RouterGroup) {
	app := api.App.GoodsApi

	// 商品新增 / 修改 / 删除（管理员权限）
	{
		g.POST("goods",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[goodsapi.GoodsAddRequest],
			app.GoodsAddView)
		g.PUT("goods",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[goodsapi.GoodsUpdateRequest],
			app.GoodsUpdateView)
		g.DELETE("goods/admin",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.GoodsRemoveView)
		g.PUT("goods/status",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[goodsapi.GoodsStatusUpdateRequest],
			app.GoodsStatusUpdateView)
	}

	// 商品查询（管理员视角）
	{
		g.GET("goods/admin",
			middleware.AdminMiddleware,
			middleware.BindQueryMiddleware[goodsapi.GoodsListRequest],
			app.GoodsListView)
		g.GET("goods/options/admin",
			middleware.AdminMiddleware,
			app.GoodsOptionsListView)
	}

	// 商品查询（用户视角）
	{
		g.GET("goods/:id",
			middleware.BindUriMiddleware[models.IDRequest],
			app.GoodsDetailView)
		g.GET("goods/category",
			app.GoodsCategoryListView)
		g.GET("goods/index",
			middleware.BindQueryMiddleware[goodsapi.GoodsIndexListRequest],
			app.GoodsIndexListView)
	}
}
