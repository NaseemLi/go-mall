package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsInfo struct {
	GoodsID    uint                  `json:"goodsId"`   //商品ID
	Cover      string                `json:"cover"`     //商品封面
	Title      string                `json:"title"`     //商品标题
	Price      int                   `json:"price"`     //实际价格
	TotalPrice int                   `json:"-"`         //商品总价
	Inventory  *int                  `json:"inventory"` //商品库存
	Num        int                   `json:"num"`       //购买数量
	Status     ctype.GoodsStatusType `json:"status"`    //商品状态
}

type CouponInfo struct {
	ID          uint             `json:"id"`          //优惠券 ID
	Type        ctype.CouponType `json:"type"`        //优惠券类型
	Title       string           `json:"title"`       //优惠券标题
	CouponPrice int              `json:"couponPrice"` //优惠券金额
	Threshold   int              `json:"threshold"`   //使用门槛
	GoodsID     *uint            `json:"goodsId"`     //关联的商品ID,如果是商品优惠卷
}

type OrderGoodsInfo struct {
	GoodsID uint   `json:"goodsID"` // 商品ID
	Num     int    `json:"num"`     // 商品数量
	Note    string `json:"note"`    // 商品备注
}

type OrderConfirmRequest struct {
	OrderGoodsList []OrderGoodsInfo `json:"orderGoodsList"` // 订单商品列表
	CouponIDList   *[]uint          `json:"couponIDList"`   // 优惠券ID列表
}

type OrderConfirmResponse struct {
	GoodsList  []GoodsInfo  `json:"goodsList"`
	CouponList []CouponInfo `json:"couponList"`
	Price      int          `json:"price"` //优惠后的金额
}

func (o *OrderApi) OrderConfirmView(c *gin.Context) {
	cr := middleware.GetBind[OrderConfirmRequest](c)
	claims := middleware.GetAuth(c)

	var totalPrice int
	var discountPrice int
	var finalPrice int

	// step 1: 获取商品信息（仅查找已上架的）
	var goodsIDList []uint
	for _, item := range cr.OrderGoodsList {
		goodsIDList = append(goodsIDList, item.GoodsID)
	}

	var goodsList []models.GoodsModel
	global.DB.Find(&goodsList, "id IN ? AND status = ?", goodsIDList, ctype.GoodsStatusTop)

	// 校验是否存在未上架商品
	if len(goodsList) != len(cr.OrderGoodsList) {
		res.FailWithMsg("存在未上架商品，请检查", c)
		return
	}

	goodsMap := make(map[uint]models.GoodsModel)
	for _, g := range goodsList {
		goodsMap[g.ID] = g
	}

	var goodsInfoList []GoodsInfo
	for _, item := range cr.OrderGoodsList {
		g := goodsMap[item.GoodsID]
		subTotal := g.Price * item.Num
		totalPrice += subTotal

		goodsInfoList = append(goodsInfoList, GoodsInfo{
			GoodsID:    g.ID,
			Cover:      g.GetCover(),
			Title:      g.Title,
			Price:      g.Price,
			TotalPrice: subTotal,
			Inventory:  g.Inventory,
			Num:        item.Num,
			Status:     g.Status,
		})
	}

	// step 2: 获取用户选中的可用优惠券（未使用 + 未过期 + 可用）
	var userCoupons []models.UserCouponModel
	query := global.DB.Where("user_id = ? AND status = ? AND end_time > now()",
		claims.UserID, ctype.CouponStatusNotUsed)
	if cr.CouponIDList != nil {
		query = query.Where("id IN ?", *cr.CouponIDList)
	}
	query.Preload("CouponModel").Find(&userCoupons)

	var couponInfoList []CouponInfo
	var usedCouponMap = make(map[uint]bool)
	var bestGeneralCoupon *CouponInfo
	bestDiscount := 0

	for _, uc := range userCoupons {
		c := uc.CouponModel

		// 商品券
		if c.Type == ctype.CouponGoodsType && c.GoodsID != nil {
			if _, ok := goodsMap[*c.GoodsID]; ok {
				discountPrice += c.CouponPrice
				usedCouponMap[uc.ID] = true
			}
			continue
		}

		// 普通券，挑选最大优惠且满足门槛的券
		if c.Threshold <= totalPrice {
			if c.CouponPrice > bestDiscount {
				bestDiscount = c.CouponPrice
				bestGeneralCoupon = &CouponInfo{
					ID:          uc.ID,
					Type:        c.Type,
					Title:       c.Title,
					CouponPrice: c.CouponPrice,
					Threshold:   c.Threshold,
					GoodsID:     c.GoodsID,
				}
			}
		}
	}

	// 只选一个最划算的普通券（如果有）
	if bestGeneralCoupon != nil {
		discountPrice += bestGeneralCoupon.CouponPrice
		couponInfoList = append(couponInfoList, *bestGeneralCoupon)
	}

	// step 3: 计算最终价格
	finalPrice = totalPrice - discountPrice
	if finalPrice < 0 {
		finalPrice = 0
	}

	// step 4: 返回结果
	res.OkWithData(OrderConfirmResponse{
		GoodsList:  goodsInfoList,
		CouponList: couponInfoList,
		Price:      finalPrice,
	}, c)
}
