package msgapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type MsgUserListRequest struct {
	models.PageInfo
	IsRead *bool `form:"isRead"` // 是否已读
}

type MsgUserListResponse struct {
	models.MessageModel
	GoodsTitle string `json:"goodsTitle"` // 商品标题
	GoodsID    uint   `json:"goodsID"`    // 商品ID
}

func (MsgApi) MsgUserListView(c *gin.Context) {
	cr := middleware.GetBind[MsgUserListRequest](c)
	claims := middleware.GetAuth(c)
	query := global.DB.Where("")

	if cr.IsRead != nil {
		query = query.Where("is_read = ?", *cr.IsRead)
	}

	_list, count, _ := common.QueryList(models.MessageModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		Where:    query,
		PageInfo: cr.PageInfo,
		Preloads: []string{"GoodsModel"},
	})

	list := make([]MsgUserListResponse, 0)
	for _, v := range _list {
		list = append(list, MsgUserListResponse{
			MessageModel: v,
			GoodsTitle:   v.GoodsModel.Title,
			GoodsID:      v.GoodsID,
		})
	}

	res.OkWithList(list, count, c)
}
