package routers

import (
	"fast_gin/api"
	orderapi "fast_gin/api/order_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(g *gin.RouterGroup) {
	app := api.App.OrderApi

	g.POST("order/confirm",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[orderapi.OrderConfirmRequest],
		app.OrderConfirmView)
}
