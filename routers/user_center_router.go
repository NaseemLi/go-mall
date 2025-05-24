package routers

import (
	"fast_gin/api"
	lookapi "fast_gin/api/user_center_api/look_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func UserCenterRouter(g *gin.RouterGroup) {
	app := api.App.UserCenterApi
	{
		g.POST("user_center/look",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[lookapi.LookGoodsRequest],
			app.LookApi.LookGoodsView)
	}
}
