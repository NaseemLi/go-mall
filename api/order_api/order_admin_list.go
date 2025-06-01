package orderapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderAdminGoodsInfo struct {
	GoodsID      uint   `json:"goodsID"`    // 商品ID
	GoodsCover   string `json:"goodsCover"` // 商品封面
	GoodsTitle   string `json:"goodsTitle"` // 商品标题
	OrderGoodsID uint   `form:"orderGoodsID"`
	Status       int8   `form:"status"` // 商品状态: 0 1 已发货
}

// 订单状态: 1代付款 2已付款/待发货 3已发货/待收货 4已收货/待评价 5已评论/已完成 6已取消 7已超时
type OrderAdminListRequest struct {
	models.PageInfo
	UserID uint   `form:"userID"`
	No     string `form:"no"`
	Status int8   `form:"status"`
}

type OrderAdminListResponse struct {
	models.OrderModel
	UserNickname   string                `json:"userNickname"`
	OrderGoodsList []OrderAdminGoodsInfo `json:"orderGoodsList"` // 订单商品信息
}

func (OrderApi) OrderAdminListView(c *gin.Context) {
	cr := middleware.GetBind[OrderAdminListRequest](c)

	_list, count, _ := common.QueryList(models.OrderModel{
		UserID: cr.UserID,
		No:     cr.No,
		Status: cr.Status,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Unscoped: true,
		Preloads: []string{"UserModel", "OrderGoodsList.GoodsModel"},
	})

	var list = make([]OrderAdminListResponse, 0)
	for _, item := range _list {
		var goodsList = make([]OrderAdminGoodsInfo, 0)
		for _, v := range item.OrderGoodsList {
			goodsList = append(goodsList, OrderAdminGoodsInfo{
				GoodsID:      v.GoodsID,
				GoodsCover:   v.GoodsModel.GetCover(),
				GoodsTitle:   v.GoodsModel.Title,
				OrderGoodsID: v.ID,
				Status:       v.Status,
			})
		}
		list = append(list, OrderAdminListResponse{
			OrderModel:     item,
			UserNickname:   item.UserModel.Nickname,
			OrderGoodsList: goodsList,
		})
	}
	res.OkWithList(list, count, c)
}
