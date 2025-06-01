package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderUserGoodsInfo struct {
	GoodsID      uint   `json:"goodsID"`      // 商品ID
	GoodsCover   string `json:"goodsCover"`   // 商品封面
	GoodsTitle   string `json:"goodsTitle"`   // 商品标题
	GoodsPrice   int    `json:"goodsPrice"`   // 商品价格
	Num          int    `json:"num"`          // 商品数量
	Note         string `json:"note"`         // 商品备注
	OrderGoodsID uint   `json:"orderGoodsID"` // 订单商品ID
}

type OrderUserListRequest struct {
	models.PageInfo
	GoodsTitle string `form:"goodsTitle"` // 商品标题
	No         string `form:"no"`         // 订单号
}

type OrderUserListResponse struct {
	ID             uint                 `json:"id"`
	CreatedAt      time.Time            `json:"createdAt"`      // 创建时间
	No             string               `json:"no"`             // 订单号
	Status         int8                 `json:"status"`         // 订单状态
	Price          int                  `json:"price"`          // 订单总价
	CouponPrice    int                  `json:"couponPrice"`    // 优惠券抵扣金额
	OrderGoodsList []OrderUserGoodsInfo `json:"orderGoodsList"` // 订单商品信息
}

func (OrderApi) OrderUserListView(c *gin.Context) {
	cr := middleware.GetBind[OrderUserListRequest](c)
	claims := middleware.GetAuth(c)

	//TODO:模糊匹配测试
	query := global.DB.Where("")
	if cr.GoodsTitle != "" {
		var goodsIDList []uint
		global.DB.Model(&models.GoodsModel{}).
			Where("title LIKE ?", "%"+cr.GoodsTitle+"%").
			Pluck("id", &goodsIDList)

		var orderIDList []uint
		if len(goodsIDList) > 0 {
			global.DB.Model(&models.OrderGoodsModel{}).
				Where("user_id = ? and goods_id IN ?", claims.UserID, goodsIDList).
				Pluck("order_id", &orderIDList)
		}

		if len(orderIDList) > 0 {
			query = query.Where("id IN ?", orderIDList)
		}
	}

	_list, count, _ := common.QueryList(models.OrderModel{
		No: cr.No,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: []string{"OrderGoodsList.GoodsModel"},
	})

	var list = make([]OrderUserListResponse, 0)
	for _, item := range _list {
		var goodsList = make([]OrderUserGoodsInfo, 0)
		for _, v := range item.OrderGoodsList {
			goodsList = append(goodsList, OrderUserGoodsInfo{
				GoodsID:      v.GoodsID,
				GoodsCover:   v.GoodsModel.GetCover(),
				GoodsTitle:   v.GoodsModel.Title,
				GoodsPrice:   v.GoodsModel.Price,
				Note:         v.Note,
				Num:          v.Num,
				OrderGoodsID: v.ID,
			})
		}
		list = append(list, OrderUserListResponse{
			ID:             item.ID,
			CreatedAt:      item.CreatedAt,
			No:             item.No,
			Status:         item.Status,
			Price:          item.Price,
			CouponPrice:    item.Coupon,
			OrderGoodsList: goodsList,
		})
	}
	res.OkWithList(list, count, c)
}
