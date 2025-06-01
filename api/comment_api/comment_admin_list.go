package commentapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CommentAdminListRequest struct {
	models.PageInfo
	UserID uint `json:"userID"` // 用户ID
}

type CommentAdminListResponse struct {
	models.CommentModel
	GoodsID         uint   `json:"goodsID"`         // 商品ID
	GoodsTitle      string `json:"goodsTitle"`      // 商品标题
	UserNickname    string `json:"userNickname"`    // 用户昵称
	OrderGoodsPrice int    `json:"orderGoodsPrice"` // 订单商品价格
}

func (CommentApi) CommentAdminListView(c *gin.Context) {
	cr := middleware.GetBind[CommentAdminListRequest](c)

	_list, count, _ := common.QueryList(models.CommentModel{
		UserID: cr.UserID,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Preloads: []string{"OrderGoodsModel.GoodsModel", "UserModel"},
	})

	var list = make([]CommentAdminListResponse, 0)
	for _, item := range _list {
		list = append(list, CommentAdminListResponse{
			CommentModel:    item,
			GoodsID:         item.OrderGoodsModel.GoodsID,
			GoodsTitle:      item.OrderGoodsModel.GoodsModel.Title,
			UserNickname:    item.UserModel.Nickname,
			OrderGoodsPrice: item.OrderGoodsModel.Price,
		})
	}
	res.OkWithList(list, count, c)
}
