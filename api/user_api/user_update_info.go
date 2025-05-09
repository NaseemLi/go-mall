package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type UpdateInfoRequest struct {
	Avatar   string `json:"avatar" binding:"required,max=256"`
	Nickname string `json:"nickname" binding:"required,max=32"`
}

func (UserApi) UpdateInfoView(c *gin.Context) {
	cr := middleware.GetBind[UpdateInfoRequest](c)
	claims := middleware.GetAuth(c)
	var user models.UserModel
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}

	global.DB.Model(&user).Updates(map[string]any{
		"avatar":   cr.Avatar,
		"nickname": cr.Nickname,
	})

	res.OkWithMsg("用户信息修改成功", c)
}
