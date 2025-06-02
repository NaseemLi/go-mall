package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func DataRouter(g *gin.RouterGroup) {
	app := api.App.DataApi

	// 用户数据相关
	{
		g.GET("data/user",
			middleware.AuthMiddleware,
			app.StatisticsUserView)
	}

	// 管理后台数据分析
	{
		g.GET("data/system",
			middleware.AdminMiddleware,
			app.StatisticsSystemView)

		g.GET("data/user_trend",
			middleware.AdminMiddleware,
			app.UserLoginTrendView)

		g.GET("data/order_trend",
			middleware.AdminMiddleware,
			app.OrderTrendView)

		g.GET("data/computer",
			middleware.AdminMiddleware,
			app.ComputerView)
	}
}
