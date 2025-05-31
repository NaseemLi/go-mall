package routers

import (
	"fast_gin/api"
	orderapi "fast_gin/api/order_api"
	"fast_gin/middleware"
	"fast_gin/models"

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
	g.PUT("order/admin",
		middleware.AdminMiddleware,
		middleware.BindQueryMiddleware[orderapi.OrderAdminListRequest],
		app.OrderAdminListView)
	g.PUT("order/user",
		middleware.AuthMiddleware,
		middleware.BindQueryMiddleware[orderapi.OrderUserListRequest],
		app.OrderUserListView)
	g.POST("order/send_out_goods",
		middleware.AdminMiddleware,
		middleware.BindJsonMiddleware[orderapi.OrderSendOutGoodsRequest],
		app.OrderSendOutGoodsView)
	g.POST("order/rev_goods",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[orderapi.OrderRevRequest],
		app.OrderRevGoodsView)
	g.GET("order/detail/:id",
		middleware.AuthMiddleware,
		middleware.BindUriMiddleware[models.IDRequest],
		app.OrderDetailView)
	g.PUT("order/callback",
		middleware.BindJsonMiddleware[orderapi.OrderPayCallbackRequest],
		app.OrderPayCallbackView)
	g.GET("order/pay/page",
		middleware.BindQueryMiddleware[orderapi.OrderPayDetailRequest],
		app.OrderPayDetailView)
	g.DELETE("order/admin/remove",
		middleware.BindJsonMiddleware[models.IDListRequest],
		app.OrderAdminRemoveView)
}
