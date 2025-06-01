package msgapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (MsgApi) MsgUserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)
	claims := middleware.GetAuth(c)

	var msgList []models.MessageModel
	global.DB.Find(&msgList, "user_id = ? and id in ?", claims.UserID, cr.IDList)
	if len(msgList) > 0 {
		global.DB.Delete(&msgList)
	}

	msg := fmt.Sprintf("消息删除成功,删除了 %d 条消息", len(msgList))

	res.OkWithMsg(msg, c)
}
