package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserRemoveRequest struct {
}

func (UserApi) UserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var userList []models.UserModel
	global.DB.Find(&userList, "id in ?", cr.IDList)
	if len(userList) > 0 {
		global.DB.Delete(&userList)
	}

	msg := fmt.Sprintf("用户删除成功,删除了 %d 个用户", len(userList))

	res.OkWithMsg(msg, c)
}
