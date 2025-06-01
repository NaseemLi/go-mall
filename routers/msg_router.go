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

	// 用户消息相关
	{
		g.GET("msg/user",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[msgapi.MsgUserListRequest],
			app.MsgUserListView)
		g.DELETE("msg/user",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.MsgUserRemoveView)
	}

	// 管理员消息相关
	{
		g.GET("msg/admin",
			middleware.AdminMiddleware,
			middleware.BindQueryMiddleware[msgapi.MsgAdminListRequest],
			app.MsgAdminListView)
		g.DELETE("msg/admin",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.MsgAdminRemoveView)
	}

	// 用户标记已读
	{
		g.GET("msg/read/:id",
			middleware.AuthMiddleware,
			middleware.BindUriMiddleware[models.IDRequest],
			app.MsgReadView)
	}
}
