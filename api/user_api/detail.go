package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

func (UserApi) UserDetailView(c *gin.Context) {
	claims := middleware.GetAuth(c)

	var user models.UserModel
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}

	res.OkWithData(user, c)
}
