package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/pwd"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type UpdatePwdRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,max=64" label:"老密码"`
	Password    string `json:"password" binding:"required,max=64" label:"密码"`
	RePassword  string `json:"rePassword" binding:"required,max=64" label:"确认密码"`
}

func (UserApi) UpdatePwdView(c *gin.Context) {
	cr := middleware.GetBind[UpdatePwdRequest](c)

	if cr.Password != cr.RePassword {
		res.FailWithMsg("两次密码不一致", c)
		return
	}

	claims := middleware.GetAuth(c)

	var user models.UserModel
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}

	if !pwd.CompareHashAndPassword(user.Password, cr.OldPassword) {
		res.FailWithMsg("密码错误", c)
		return
	}

	hashPwd := pwd.GenerateFromPassword(cr.Password)

	global.DB.Model(&user).Updates(models.UserModel{
		Password: hashPwd,
	})

	res.OkWithMsg("密码修改成功", c)
}
