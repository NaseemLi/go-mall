package routers

import (
	"fast_gin/api"
	msgapi "fast_gin/api/msg_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func MsgRouter(g *gin.RouterGroup) {
	app := api.App.MsgApi

	// 消息列表
	{
		//用户消息列表
		g.GET("msg/user",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[msgapi.MsgUserListRequest],
			app.MsgUserListView)

		//管理员消息列表
		g.GET("msg/admin",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[msgapi.MsgAdminListRequest],
			app.MsgAdminListView)

		g.GET("msg/read/:id",
			middleware.AuthMiddleware,
			middleware.BindUriMiddleware[models.IDRequest],
			app.MsgReadView)
	}
}
