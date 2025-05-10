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
	g.POST("goods",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[goodsapi.GoodsAddRequest],
		app.GoodsAddView,
	)

	g.GET("goods/admin",
		middleware.AdminMiddleware,
		middleware.BindQueryMiddleware[goodsapi.GoodsListRequest],
		app.GoodsListView,
	)

	g.DELETE("goods/admin",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[models.IDListRequest],
		app.GoodsRemoveView,
	)

	g.PUT("goods",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[goodsapi.GoodsUpdateRequest],
		app.GoodsUpdateView,
	)

	g.GET("goods/:id",
		middleware.BindUriMiddleware[models.IDRequest],
		app.GoodsDetailView,
	)

	g.PUT("goods/status",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[goodsapi.GoodsStatusUpdateRequest],
		app.GoodsStatusUpdateView,
	)

	g.GET("goods/category",
		app.GoodsCategoryListView,
	)
	g.GET("goods/options/admin",
		middleware.AdminMiddleware,
		app.GoodsOptionsListView,
	)
}
