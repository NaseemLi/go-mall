package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"sort"

	"github.com/gin-gonic/gin"
)

type GoodsInfo struct {
	CarID      uint                  `json:"carId"`      //购物车ID
	GoodsID    uint                  `json:"goodsId"`    //商品ID
	Cover      string                `json:"cover"`      //商品封面
	Title      string                `json:"title"`      //商品标题
	Price      int                   `json:"price"`      //实际价格
	TotalPrice int                   `json:"totalPrice"` //商品总价
	PayPrice   int                   `json:"payPrice"`   //实际支付价格
	Inventory  *int                  `json:"inventory"`  //商品库存
	Num        int                   `json:"num"`        //购买数量
	Status     ctype.GoodsStatusType `json:"status"`     //商品状态
	Used       bool                  `json:"used"`       //是否选购
	CouponInfo *CouponInfo           `json:"couponInfo"` //优惠券信息
}

type CouponInfo struct {
	ID          uint             `json:"id"`          //优惠券 ID
	Type        ctype.CouponType `json:"type"`        //优惠券类型
	Title       string           `json:"title"`       //优惠券标题
	CouponPrice int              `json:"couponPrice"` //优惠券金额
	Used        bool             `json:"used"`        //是否使用
	Threshold   int              `json:"threshold"`   //使用门槛
	SubPrice    int              `json:"subPrice"`    //差多少,为 0 表示可选
	GoodsID     *uint            `json:"goodsId"`     //关联的商品ID,如果是商品优惠卷
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
		Where:    global.DB.Where("status = ?", ctype.CarStatusPending),
		Likes:    []string{"goods_title"},
		Preloads: []string{"GoodsModel"},
	})

	//我能用的商品优惠卷
	var couponList []models.UserCouponModel
	global.DB.
		Preload("CouponModel").
		Find(&couponList, "user_id = ? AND status = ? AND end_time > now()",
			claims.UserID, ctype.CouponStatusNotUsed)

	var useCouponMap = map[uint]bool{}
	for _, v := range cr.CouponIDList {
		useCouponMap[v] = true
	}

	var couponInfoList = make([]CouponInfo, 0)
	var goodsCouponList = make([]CouponInfo, 0)
	var couponGoodsMap = map[uint]*CouponInfo{} //商品优惠卷的商品ID

	for _, v := range couponList {
		info := CouponInfo{
			ID:          v.ID,
			Type:        v.CouponModel.Type,
			CouponPrice: v.CouponModel.CouponPrice,
			Used:        useCouponMap[v.ID],    //默认未使用
			GoodsID:     v.CouponModel.GoodsID, //如果是商品优惠卷,则有商品ID
			Threshold:   v.CouponModel.Threshold,
			//SubPrice:    v.CouponModel.CouponPrice, //差多少
		}
		if v.CouponModel.GoodsID != nil {
			couponGoodsMap[*v.CouponModel.GoodsID] = &info //商品优惠卷的商品ID
		}
		if v.CouponModel.Type == ctype.CouponGoodsType {
			goodsCouponList = append(goodsCouponList, info)
		}
		couponInfoList = append(couponInfoList, info)
	}

	sort.Slice(couponInfoList, func(i, j int) bool {
		if couponInfoList[i].Used && !couponInfoList[j].Used {
			return true // i 用了，j 没用 → i 前面
		}
		if !couponInfoList[i].Used && couponInfoList[j].Used {
			return false // i 没用，j 用了 → j 前面
		}
		// 如果两个都用 or 都没用，就按类型从小到大排
		return uint(couponInfoList[i].Type) < uint(couponInfoList[j].Type)
	})

	var useCarMap = map[uint]bool{}
	for _, v := range cr.CarIDList {
		useCarMap[v] = true
	}

	//如果有商品优惠卷,得判断有没有选择这样优惠卷中指定的商品

	var totalPrice int
	var couponPrice int //商品优惠卷的优惠金额
	var goodsList []GoodsInfo
	for _, v := range _list {
		goods := GoodsInfo{
			CarID:      v.ID,
			GoodsID:    v.GoodsModel.ID,
			Cover:      v.GoodsModel.Images[0],
			Title:      v.GoodsModel.Title,
			TotalPrice: v.GoodsModel.Price * v.Num,
			Price:      v.GoodsModel.Price,
			PayPrice:   v.GoodsModel.Price * v.Num,
			Inventory:  v.GoodsModel.Inventory,
			Num:        v.Num,
			Status:     v.GoodsModel.Status,
			Used:       useCarMap[v.ID],
			CouponInfo: couponGoodsMap[v.GoodsID],
		}

		if goods.Used {
			//判断这个商品有没有可用的商品优惠卷
			for _, info := range goodsCouponList {
				if info.Used && info.Type == ctype.CouponGoodsType &&
					*info.GoodsID == v.GoodsID && info.Threshold <= goods.TotalPrice {
					//我选择的商品优惠卷,我也选择了对应的商品
					if goods.CouponInfo != nil {
						goods.CouponInfo.Used = true
						goods.PayPrice -= goods.CouponInfo.CouponPrice
						couponPrice += goods.CouponInfo.CouponPrice
					}
				}
			}
		}
		if goods.Used {
			totalPrice += goods.TotalPrice
		}
		goodsList = append(goodsList, goods)
	}

	price := totalPrice - couponPrice
	//算其他优惠卷的优惠金额
	for _, info := range couponInfoList {
		if info.Used && info.Type != ctype.CouponGoodsType && info.Threshold <= price {
			price -= info.CouponPrice
			info.SubPrice = 0
		} else {
			info.SubPrice = info.Threshold - price
		}
	}

	if price < 0 {
		price = 0 //不能为负数
	}

	data := CarListResponse{
		Count:      count,
		GoodsList:  goodsList,
		CouponList: couponInfoList,
		TotalPrice: totalPrice,
		Price:      price,
	}

	res.OkWithData(data, c)
}
