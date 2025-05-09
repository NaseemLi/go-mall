package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type AdminUpdateInfoRequest struct {
	UserID   uint       `json:"userID" binding:"required"`
	Avatar   string     `json:"avatar" binding:"required,max=256"`
	Nickname string     `json:"nickname" binding:"required,max=32"`
	RoleID   ctype.Role `json:"roleID"`
}

func (UserApi) AdminUpdateInfoView(c *gin.Context) {
	cr := middleware.GetBind[AdminUpdateInfoRequest](c)
	var user models.UserModel
	err := global.DB.Take(&user, cr.UserID).Error
	if err != nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}

	global.DB.Model(&user).Updates(models.UserModel{
		Avatar:   cr.Avatar,
		Nickname: cr.Nickname,
		RoleID:   cr.RoleID,
	})

	res.OkWithMsg("用户信息修改成功", c)
}
