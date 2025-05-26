package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/random"
	"fast_gin/utils/res"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
)

type OrderPayRequest struct {
	AddrID         uint             `json:"addrID" binding:"required"`              // 收货地址ID
	PayType        int8             `json:"payType" binding:"required,oneof=1 2 3"` // 支付方式 1:站内支付 2:微信支付 3:支付宝支付
	OrderGoodsList []OrderGoodsInfo `json:"orderGoodsList" binding:"required"`      // 订单商品列表
	CouponIDList   []uint           `json:"couponIDList"`                           // 优惠券ID列表
	CarIDList      []uint           `json:"carIDList"`                              // 购物车ID列表,如果有购物车,则删除购物车
}

type OrderPayResponse struct {
	No     string `json:"no"`     // 订单号
	PayUrl string `json:"payUrl"` // 支付链接
	Price  int    `json:"price"`  // 支付金额
}

func (OrderApi) OrderPayView(c *gin.Context) {
	cr := middleware.GetBind[OrderPayRequest](c)
	claims := middleware.GetAuth(c)

	//检验地址
	var addr models.AddrModel
	if err := global.DB.Take(&addr, "user_id = ? and id = ?", claims.ID, cr.AddrID).Error; err != nil {
		res.FailWithMsg("地址不存在", c)
		return
	}

	//优惠卷校验
	var myCouponList []models.UserCouponModel
	var goodsCouponMap = map[uint]models.UserCouponModel{}
	query := global.DB.Where("id in ? AND user_id = ? AND status = ? AND end_time > now()", cr.CouponIDList, claims.ID, ctype.CouponStatusNotUsed)
	// if cr.CouponIDList != nil {
	// 	query = query.Where("id in ?", cr.CouponIDList)
	// }
	global.DB.Where(query).Preload("CouponModel").Find(&myCouponList)
	if len(cr.CouponIDList) != len(myCouponList) {
		res.FailWithMsg("优惠卷校验失败", c)
		return
	}
	for _, v := range myCouponList {
		if v.CouponModel.Type == ctype.CouponGoodsType && v.CouponModel.GoodsID != nil {
			goodsCouponMap[*v.CouponModel.GoodsID] = v
		}
	}

	//判断来自购物车还是网页直接支付,有没有重复下单
	if len(cr.CarIDList) > 0 {
		var carList []models.CarModel
		if err := global.DB.Find(&carList, "user_id = ? and id in ? and status = 0", claims.ID, cr.CarIDList).Error; err != nil {
			res.FailWithMsg("获取购物车信息失败", c)
			return
		}
		if len(carList) != len(cr.CarIDList) {
			res.FailWithMsg("购物车商品重复下单", c)
			return
		}
	}

	//检验库存
	//找商品
	var goodsIDList []uint
	var orderGoodsMap = map[uint]OrderGoodsInfo{}
	for _, v := range cr.OrderGoodsList {
		goodsIDList = append(goodsIDList, v.GoodsID)
		orderGoodsMap[v.GoodsID] = v
	}
	var GoodsList []models.GoodsModel
	global.DB.Find(&GoodsList, "id in ? and status = ?", goodsIDList, ctype.GoodsStatusTop)
	if len(GoodsList) != len(cr.OrderGoodsList) {
		res.FailWithMsg("存在未下架商品请检查", c)
		return
	}

	var price int
	var couponPrice int
	for _, goodsModel := range GoodsList {
		info := orderGoodsMap[goodsModel.ID]
		if goodsModel.Inventory != nil {
			// 检查库存
			if info.Num > *goodsModel.Inventory {
				res.FailWithMsg(fmt.Sprintf("商品: %s 库存不足", goodsModel.Title), c)
				return
			}
		}

		goodsPrice := info.Num * goodsModel.Price
		coupon, ok := goodsCouponMap[goodsModel.ID]
		if ok && goodsPrice >= coupon.CouponModel.Threshold {
			//如果有优惠卷,则计算优惠
			goodsPrice -= coupon.CouponModel.CouponPrice
			couponPrice += coupon.CouponModel.CouponPrice
		}
		price += goodsPrice
	}

	//生成订单号
	no := random.GenerateOrderNumber()

	//算金额
	sort.Slice(myCouponList, func(i, j int) bool {
		return myCouponList[i].CouponModel.Type < myCouponList[j].CouponModel.Type
	})

	for _, info := range myCouponList {
		if info.CouponModel.Type != ctype.CouponGoodsType && price >= info.CouponModel.Threshold {
			price -= info.CouponModel.CouponPrice
		}
	}

	//创建订单
	var order = models.OrderModel{
		No:      no,
		UserID:  claims.UserID,
		AddrID:  cr.AddrID,
		Price:   price,
		Status:  1,
		PayType: cr.PayType,
		Coupon:  couponPrice,
	}

	if err := global.DB.Create(&order).Error; err != nil {
		res.FailWithMsg("创建订单失败", c)
		return
	}

	//锁后优惠卷
	//改购物车状态
}
