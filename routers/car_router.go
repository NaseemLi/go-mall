package routers

import (
	"fast_gin/api"
	carapi "fast_gin/api/car_api"
	"fast_gin/middleware"

	"github.com/gin-gonic/gin"
)

func CarRouter(g *gin.RouterGroup) {
	app := api.App.CarApi
	g.POST("car",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[carapi.CarCreateRequest],
		app.CarCreateView)
}
