package routers

import (
	"fast_gin/api"

	"github.com/gin-gonic/gin"
)

func CaptchaRouter(g *gin.RouterGroup) {
	app := api.App.CaptchaApi

	// 生成图形验证码（无需鉴权）
	{
		g.GET("captcha/generate", app.GenerateView)
	}
}
