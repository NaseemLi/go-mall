package routers

import (
	"fast_gin/api"
	"fast_gin/api/user_api"
	"fast_gin/middleware"
	"fast_gin/models"

	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	app := api.App.UserApi

	// 用户认证相关
	{
		g.POST("users/login",
			middleware.BindJsonMiddleware[user_api.LoginRequest],
			app.LoginView)
		g.POST("users/register",
			middleware.BindJsonMiddleware[user_api.RegisterRequest],
			app.RegisterView)
		g.POST("users/logout",
			middleware.AuthMiddleware,
			app.LogoutView)
	}

	// 用户自身信息操作
	{
		g.PUT("users/pwd",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[user_api.UpdatePwdRequest],
			app.UpdatePwdView)
		g.PUT("users/info",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[user_api.UpdateInfoRequest],
			app.UpdateInfoView)
		g.GET("users/detail",
			middleware.AuthMiddleware,
			app.UserDetailView)
	}

	// 管理员操作用户信息
	{
		g.PUT("users/info/admin",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[user_api.AdminUpdateInfoRequest],
			app.AdminUpdateInfoView)
		g.GET("users",
			middleware.LimitMiddleware(10),
			middleware.AdminMiddleware,
			middleware.BindQueryMiddleware[models.PageInfo],
			app.UserListView)
		g.DELETE("users",
			middleware.AdminMiddleware,
			middleware.BindJsonMiddleware[models.IDListRequest],
			app.UserRemoveView)
	}
}
