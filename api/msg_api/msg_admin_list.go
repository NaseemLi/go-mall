package msgapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type MsgAdminListRequest struct {
	models.PageInfo
	UserID uint `form:"userID"` // 用户ID
}

type MsgAdminListResponse struct {
	models.MessageModel
	GoodsTitle   string `json:"goodsTitle"`   // 商品标题
	GoodsID      uint   `json:"goodsID"`      // 商品ID
	UserNickname string `json:"userNickname"` // 用户昵称
}

func (MsgApi) MsgAdminListView(c *gin.Context) {
	cr := middleware.GetBind[MsgAdminListRequest](c)

	_list, count, _ := common.QueryList(models.MessageModel{
		UserID: cr.UserID,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Unscoped: true,
		Preloads: []string{"GoodsModel", "UserModel"},
	})

	list := make([]MsgAdminListResponse, 0)
	for _, v := range _list {
		list = append(list, MsgAdminListResponse{
			MessageModel: v,
			GoodsTitle:   v.GoodsModel.Title,
			GoodsID:      v.GoodsID,
			UserNickname: v.UserModel.Nickname,
		})
	}

	res.OkWithList(list, count, c)
}
