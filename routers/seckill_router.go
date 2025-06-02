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

	// 管理端：秒杀管理
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
	}

	// 前台展示：秒杀活动时间 & 秒杀商品列表
	{
		g.GET("sec_kill/date",
			app.IndexDateListView)

		g.GET("sec_kill/goods",
			middleware.BindQueryMiddleware[seckillapi.IndexSecKillGoodsListRequest],
			app.IndexSecKillGoodsListView)
	}

	// 用户端：参与秒杀 & 查询结果
	{
		g.POST("sec_kill/user",
			middleware.LimitMiddleware(100),
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[seckillapi.SecKillRequest],
			app.SecKillView)

		g.POST("sec_kill/detail",
			middleware.LimitMiddleware(100),
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[seckillapi.SecKillDetailRequest],
			app.SecKillDetailView)

		g.POST("sec_kill/order",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[seckillapi.SecKillOrderRequest],
			app.SecKillOrderView)
	}
}
