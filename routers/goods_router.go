package routers

import (
	"fast_gin/api"
	goodsapi "fast_gin/api/goods_api"
	"fast_gin/middleware"

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
}
