package routers

import (
	"fast_gin/api"
	commentapi "fast_gin/api/comment_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRouter(g *gin.RouterGroup) {
	app := api.App.CommentApi

	// 用户评价操作（创建 & 查询自己的）
	{
		g.POST("comment",
			middleware.AuthMiddleware,
			middleware.BindJsonMiddleware[commentapi.CommentCreateRequest],
			app.CommentCreateView)
		g.GET("comment/user",
			middleware.AuthMiddleware,
			middleware.BindQueryMiddleware[commentapi.CommentUserListRequest],
			app.CommentUserListView)
	}

	// 管理员查看所有评价
	{
		g.GET("comment/admin",
			middleware.AdminMiddleware,
			middleware.BindQueryMiddleware[commentapi.CommentAdminListRequest],
			app.CommentAdminListView)
	}

	// 商品评价展示（公开）
	{
		g.GET("comment/level",
			middleware.BindQueryMiddleware[commentapi.CommentLevelListRequest],
			app.CommentLevelListView)
		g.GET("comment/goods",
			middleware.BindQueryMiddleware[commentapi.GoodsCommentListRequest],
			app.GoodsCommentListView)
	}
}
