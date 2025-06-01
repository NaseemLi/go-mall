package commentapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CommentUserListRequest struct {
	models.PageInfo
}

type CommentUserListResponse struct {
	models.CommentModel
	GoodsID    uint   `json:"goodsID"`    // 商品ID
	GoodsTitle string `json:"goodsTitle"` // 商品标题
}

func (CommentApi) CommentUserListView(c *gin.Context) {
	cr := middleware.GetBind[CommentUserListRequest](c)
	claims := middleware.GetAuth(c)

	_list, count, _ := common.QueryList(models.CommentModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Preloads: []string{"OrderGoodsModel.GoodsModel"},
	})

	var list = make([]CommentUserListResponse, 0)
	for _, item := range _list {
		list = append(list, CommentUserListResponse{
			CommentModel: item,
			GoodsID:      item.OrderGoodsModel.GoodsID,
			GoodsTitle:   item.OrderGoodsModel.GoodsModel.Title,
		})
	}
	res.OkWithList(list, count, c)
}
