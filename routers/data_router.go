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
	}
}
