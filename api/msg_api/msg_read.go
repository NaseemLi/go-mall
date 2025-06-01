package msgapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type MsgReadViewResponse struct {
	MsgList []string `json:"msgList"`
}

func (MsgApi) MsgReadView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := middleware.GetAuth(c)

	var msg models.MessageModel
	err := global.DB.Take(&msg, "user_id = ? AND id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}

	if !msg.IsRead {
		err := global.DB.Model(&msg).Where("id = ?", cr.ID).Update("is_read", true).Error
		if err != nil {
			res.FailWithMsg("消息已读失败", c)
			return
		}
	}

	data := MsgReadViewResponse{MsgList: msg.MsgList}

	res.OkWithData(data, c)
}
