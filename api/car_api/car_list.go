package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsInfo struct {
	CarID     uint                  `json:"carId"`     //购物车ID
	GoodsID   uint                  `json:"goodsId"`   //商品ID
	Cover     string                `json:"cover"`     //商品封面
	Title     string                `json:"title"`     //商品标题
	Price     int                   `json:"price"`     //商品价格
	Inventory *int                  `json:"inventory"` //商品库存
	Num       int                   `json:"num"`       //购买数量
	Status    ctype.GoodsStatusType `json:"status"`    //商品状态
	Used      bool                  `json:"used"`      //是否选购
}

type CouponInfo struct {
	ID          uint             `json:"id"`          //优惠券 ID
	Type        ctype.CouponType `json:"type"`        //优惠券类型
	Title       string           `json:"title"`       //优惠券标题
	CouponPrice int              `json:"couponPrice"` //优惠券金额
	Used        bool             `json:"used"`        //是否使用
	Threshold   int              `json:"threshold"`   //使用门槛
	SubPrice    int              `json:"subPrice"`    //差多少,为 0 表示可选
}

type CarListRequest struct {
	models.PageInfo
	CarIDList    []uint `json:"carIdList"`    //购物车ID列表
	CouponIDList []uint `json:"couponIdList"` //优惠券ID列表
}

type CarListResponse struct {
	GoodsList  []GoodsInfo  `json:"goodsList"`
	CouponList []CouponInfo `json:"couponList"`
	Count      int64        `json:"count"`      //购物车商品总数
	TotalPrice int          `json:"totalPrice"` //优惠前的金额
	Price      int          `json:"price"`      //优惠后的金额
}

func (CarApi) CarListView(c *gin.Context) {
	page := middleware.GetBind[models.PageInfo](c)

	var cr CarListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	claims := middleware.GetAuth(c)

	_list, count, _ := common.QueryList(models.CarModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: page,
		Likes:    []string{"goods_title"},
		Preloads: []string{"GoodsModel"},
	})

	//我能用的商品优惠卷
	var couponList []models.UserCouponModel
	global.DB.
		Preload("CouponModel").
		Find(&couponList, "user_id = ? AND status = ? AND end_time > now()", claims.UserID, ctype.CouponStatusNotUsed)

	var useCouponMap = map[uint]bool{}
	for _, v := range cr.CouponIDList {
		useCouponMap[v] = true
	}

	var couponInfoList = make([]CouponInfo, 0)
	for _, v := range couponList {
		couponInfoList = append(couponInfoList, CouponInfo{
			ID:          v.ID,
			Type:        v.CouponModel.Type,
			CouponPrice: v.CouponModel.CouponPrice,
			Used:        useCouponMap[v.ID], //默认未使用
			//SubPrice:    v.CouponModel.CouponPrice, //差多少
			Threshold: v.CouponModel.Threshold,
		})
	}

	var useCarMap = map[uint]bool{}
	for _, v := range cr.CarIDList {
		useCarMap[v] = true
	}

	var totalPrice int
	var goodsList []GoodsInfo
	for _, v := range _list {
		goodsList = append(goodsList, GoodsInfo{
			CarID:     v.ID,
			GoodsID:   v.GoodsModel.ID,
			Cover:     v.GoodsModel.Images[0],
			Title:     v.GoodsModel.Title,
			Price:     v.GoodsModel.Price,
			Inventory: v.GoodsModel.Inventory,
			Num:       v.Num,
			Status:    v.GoodsModel.Status,
			Used:      useCarMap[v.ID],
		})

		if useCarMap[v.ID] {
			totalPrice += v.GoodsModel.Price
		}
	}

	price := totalPrice
	//算优惠

	data := CarListResponse{
		Count:      count,
		GoodsList:  goodsList,
		CouponList: couponInfoList,
		TotalPrice: price,
	}

	res.OkWithData(data, c)
}
