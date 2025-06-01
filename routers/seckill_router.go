package routers

import (
	"fast_gin/api"
	seckillapi "fast_gin/api/seckill_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func SecKillRouter(g *gin.RouterGroup) {
	app := api.App.SecKillApi

	// 创建秒杀
	{
		g.POST("sec_kill",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[seckillapi.CreateResquest],
			app.CreateView)
		g.GET("sec_kill",
			middleware.AdminMiddleware,
			middleware.BindQueryMiddleware[seckillapi.ListRequest],
			app.ListView)
		g.DELETE("sec_kill",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.RemoveView)
		g.GET("sec_kill/date",
			app.IndexDateListView)
		g.GET("sec_kill/goods",
			middleware.BindQueryMiddleware[seckillapi.IndexSecKillGoodsListRequest],
			app.IndexSecKillGoodsListView)
		g.POST("sec_kill/user",
			middleware.LimitMiddleware(100),
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[seckillapi.SecKillRequest],
			app.SecKillView)
	}
}
