package routers

import (
	"fast_gin/api"
	seckillapi "fast_gin/api/seckill_api"
	"fast_gin/middleware"

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
	}
}
