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
	g.POST("order/pay",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[orderapi.OrderPayRequest],
		app.OrderPayView)
	g.GET("order/status",
		middleware.AuthMiddleware,
		middleware.BindQueryMiddleware[orderapi.OrderStatusRequest],
		app.OrderStatusView)
	g.PUT("order/note",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[orderapi.OrderNoteUpdateRequest],
		app.OrderNoteUpdateView)
}
