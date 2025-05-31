package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func ImageRouter(g *gin.RouterGroup) {
	app := api.App.ImageApi

	// 图片上传（需登录）
	{
		g.POST("images/upload",
			middleware.AuthMiddleware,
			app.UploadView)
	}
}
