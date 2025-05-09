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
	g.POST("users/login",
		middleware.BindJsonMiddleware[user_api.LoginRequest],
		app.LoginView)
	g.POST("users/register",
		middleware.BindJsonMiddleware[user_api.RegisterRequest],
		app.RegisterView)
	g.GET("users/detail",
		middleware.AuthMiddleware,
		app.UserDetailView)
	g.GET("users",
		middleware.LimitMiddleware(10),
		middleware.AdminMiddleware,
		middleware.BindQueryMiddleware[models.PageInfo],
		app.UserListView)
	g.POST("users/logout",
		middleware.AuthMiddleware,
		app.LogoutView)
}
