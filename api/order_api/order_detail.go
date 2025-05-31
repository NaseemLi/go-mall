package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type OrderDetailCoupon struct {
	CouponPrice int              `json:"couponPrice"`
	Type        ctype.CouponType `json:"type"`
}

type OrderDetailGoods struct {
	GoodsID      uint   `json:"goodsID"`
	OrderGoodsID uint   `json:"orderGoodsID"`
	Cover        string `json:"cover"`
	Title        string `json:"title"`
	Price        int    `json:"price"`
	Num          int    `json:"num"`
	Note         string `json:"note"`
}
type OrderDetailResponse struct {
	models.OrderModel
	GoodsList  []OrderDetailGoods  `json:"goodsList"`
	CouponList []OrderDetailCoupon `json:"couponList"`
	AddrInfo   models.AddrModel    `json:"addrInfo"`
}

func (OrderApi) OrderDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	claims := middleware.GetAuth(c)

	var order models.OrderModel
	err := global.DB.
		Preload("AddrModel").
		Preload("OrderGoodsList.GoodsModel").
		Preload("UserCouponList.UserCouponModel.CouponModel").
		Take(&order, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("订单不存在", c)
		return
	}

	var goodsList = make([]OrderDetailGoods, 0)
	for _, model := range order.OrderGoodsList {
		goodsList = append(goodsList, OrderDetailGoods{
			GoodsID:      model.GoodsID,
			OrderGoodsID: model.ID,
			Cover:        model.GoodsModel.GetCover(),
			Title:        model.GoodsModel.Title,
			Price:        model.GoodsModel.Price,
			Num:          model.Num,
			Note:         model.Note,
		})
	}

	var couponList = make([]OrderDetailCoupon, 0)
	for _, model := range order.UserCouponList {
		couponList = append(couponList, OrderDetailCoupon{
			CouponPrice: model.UserCouponModel.CouponModel.CouponPrice,
			Type:        model.UserCouponModel.CouponModel.Type,
		})
	}

	data := OrderDetailResponse{
		OrderModel: order,
		GoodsList:  goodsList,
		CouponList: couponList,
		AddrInfo:   order.AddrModel,
	}

	res.OkWithData(data, c)

}
