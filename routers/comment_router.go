package routers

import (
	"fast_gin/api"
	commentapi "fast_gin/api/comment_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRouter(g *gin.RouterGroup) {
	app := api.App.CommentApi

	// 创建评价
	{
		g.POST("comment",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[commentapi.CommentCreateRequest],
			app.CommentCreateView)
	}
}
