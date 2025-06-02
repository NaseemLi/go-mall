package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func DataRouter(g *gin.RouterGroup) {
	app := api.App.DataApi

	{
		g.GET("data/user",
			middleware.AuthMiddleware,
			app.StatisticsUserView)
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
