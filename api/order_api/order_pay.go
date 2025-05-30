package orderapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	payser "fast_gin/service/pay_ser"
	redisdelay "fast_gin/service/redis_ser/redis_delay"
	"fast_gin/utils/random"
	"fast_gin/utils/res"
	"fmt"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderPayRequest struct {
	AddrID         uint             `json:"addrID" binding:"required"`
	PayType        int8             `json:"payType" binding:"required,oneof=1 2 3"`
	OrderGoodsList []OrderGoodsInfo `json:"orderGoodsList" binding:"required"`
	CouponIDList   []uint           `json:"couponIDList"`
	CarIDList      []uint           `json:"carIDList"`
}

type OrderPayResponse struct {
	No     string `json:"no"`
	PayUrl string `json:"payUrl"`
	Price  int    `json:"price"`
}

func (OrderApi) OrderPayView(c *gin.Context) {
	cr := middleware.GetBind[OrderPayRequest](c)
	claims := middleware.GetAuth(c)

	var addr models.AddrModel
	if err := global.DB.Take(&addr, "user_id = ? and id = ?",
		claims.UserID, cr.AddrID).Error; err != nil {
		res.FailWithMsg("地址不存在", c)
		return
	}

	var myCouponList []models.UserCouponModel
	var goodsCouponMap = map[uint]models.UserCouponModel{}
	if len(cr.CouponIDList) > 0 {
		global.DB.Where("id in ? AND user_id = ? AND status = ? AND end_time > now()",
			cr.CouponIDList, claims.UserID, ctype.CouponStatusNotUsed).
			Preload("CouponModel").Find(&myCouponList)

		if len(myCouponList) == 0 {
			res.FailWithMsg("请选择有效的优惠券", c)
			return
		}
		if len(myCouponList) != len(cr.CouponIDList) {
			res.FailWithMsg("部分优惠券无效，请重新选择", c)
			return
		}
		for _, v := range myCouponList {
			if v.CouponModel.Type == ctype.CouponGoodsType && v.CouponModel.GoodsID != nil {
				goodsCouponMap[*v.CouponModel.GoodsID] = v
			}
		}
	}

	if len(cr.CarIDList) > 0 {
		var carList []models.CarModel
		err := global.DB.Find(&carList, "user_id = ? and id in ? and status = ?",
			claims.UserID, cr.CarIDList, ctype.CarStatusPending).Error
		if err != nil || len(carList) != len(cr.CarIDList) {
			res.FailWithMsg("购物车商品重复下单", c)
			return
		}
	}

	var goodsIDList []uint
	orderGoodsMap := map[uint]OrderGoodsInfo{}
	for _, v := range cr.OrderGoodsList {
		if v.Num <= 0 || v.Num > 100 {
			res.FailWithMsg("非法的商品数量", c)
			return
		}
		goodsIDList = append(goodsIDList, v.GoodsID)
		orderGoodsMap[v.GoodsID] = v
	}

	var GoodsList []models.GoodsModel
	global.DB.Find(&GoodsList, "id in ? and status = ?", goodsIDList, ctype.GoodsStatusTop)
	if len(GoodsList) != len(cr.OrderGoodsList) {
		res.FailWithMsg("存在未上架商品，请检查", c)
		return
	}

	type GoodsInfo struct {
		models.GoodsModel
		GoodsPrice int
	}
	price := 0
	couponPrice := 0
	goodsMap := map[uint]GoodsInfo{}
	for _, goodsModel := range GoodsList {
		info := orderGoodsMap[goodsModel.ID]
		if goodsModel.Inventory != nil && info.Num > *goodsModel.Inventory {
			res.FailWithMsg(fmt.Sprintf("商品: %s 库存不足", goodsModel.Title), c)
			return
		}

		goodsPrice := info.Num * goodsModel.Price
		if coupon, ok := goodsCouponMap[goodsModel.ID]; ok && goodsPrice >= coupon.CouponModel.Threshold {
			goodsPrice -= coupon.CouponModel.CouponPrice
			couponPrice += coupon.CouponModel.CouponPrice
		}
		price += goodsPrice
		goodsMap[goodsModel.ID] = GoodsInfo{
			GoodsModel: goodsModel,
			GoodsPrice: goodsPrice,
		}
	}
	//重复下单时,若 15min 中内有下单,则挨个查询订单里商品和传入商品是否重复
	var myOrderList []models.OrderModel
	global.DB.Order("created_at desc").
		Preload("OrderGoodsList").
		Find(&myOrderList, "user_id = ? and status = ? and created_at > ?",
			claims.UserID, 1, time.Now().Add(-15*time.Minute))

	if len(myOrderList) > 0 {
		for _, v := range myOrderList {
			if len(v.OrderGoodsList) != len(cr.OrderGoodsList) {
				continue
			}
			var repeat bool = true
			sort.Slice(v.OrderGoodsList, func(i, j int) bool {
				return v.OrderGoodsList[i].GoodsID < v.OrderGoodsList[j].GoodsID
			})
			sort.Slice(cr.OrderGoodsList, func(i, j int) bool {
				return cr.OrderGoodsList[i].GoodsID < cr.OrderGoodsList[j].GoodsID
			})
			for i := 0; i < len(cr.OrderGoodsList); i++ {
				if v.OrderGoodsList[i].GoodsID != cr.OrderGoodsList[i].GoodsID ||
					v.OrderGoodsList[i].Num != cr.OrderGoodsList[i].Num {
					repeat = false
					break
				}
			}
			if repeat {
				logrus.Warnf("用户重复下单,订单号: %v", v.ID)
				res.FailWithMsg("存在未支付的重复订单,请勿重复下单", c)
				return
			}
		}
	}

	// 生成订单号
	no := random.GenerateOrderNumber()

	sort.Slice(myCouponList, func(i, j int) bool {
		return myCouponList[i].CouponModel.Type < myCouponList[j].CouponModel.Type
	})

	for _, info := range myCouponList {
		if info.CouponModel.Type != ctype.CouponGoodsType && price >= info.CouponModel.Threshold {
			price -= info.CouponModel.CouponPrice
			couponPrice += info.CouponModel.CouponPrice
		}
	}

	if price < 0 {
		price = 0
	}

	payUrl, err := payser.Pay(cr.PayType, no, price)
	if err != nil || payUrl == "" {
		logrus.Errorf("支付服务调用失败: %v", err)
		res.FailWithMsg("支付服务调用失败", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		order := models.OrderModel{
			No:        no,
			UserID:    claims.UserID,
			AddrID:    cr.AddrID,
			Price:     price,
			Status:    1,
			PayType:   cr.PayType,
			Coupon:    couponPrice,
			PayTime:   time.Now(),
			PayUrl:    payUrl,
			CarIDList: cr.CarIDList,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var goodsOrderList []models.OrderGoodsModel
		for _, v := range cr.OrderGoodsList {
			goodsOrderList = append(goodsOrderList, models.OrderGoodsModel{
				OrderID: order.ID,
				GoodsID: v.GoodsID,
				UserID:  claims.UserID,
				Price:   goodsMap[v.GoodsID].GoodsPrice,
				Num:     v.Num,
				Note:    v.Note,
			})
		}
		if err := tx.Create(&goodsOrderList).Error; err != nil {
			return err
		}

		if len(myCouponList) > 0 {
			var orderCouponList []models.OrderCouponModel
			for _, v := range myCouponList {
				orderCouponList = append(orderCouponList, models.OrderCouponModel{
					OrderID:      order.ID,
					UserID:       claims.UserID,
					UserCouponID: v.ID,
				})
			}
			if err := tx.Create(&orderCouponList).Error; err != nil {
				return err
			}
			err := tx.Model(&models.UserCouponModel{}).
				Where("id in ? AND user_id = ? AND status = ?",
					cr.CouponIDList, claims.UserID, ctype.CouponStatusNotUsed).
				Update("status", ctype.CouponStatusLocked).Error
			if err != nil {
				return err
			}
		}

		if len(cr.CarIDList) > 0 {
			var carList []models.CarModel
			tx.Find(&carList, "user_id = ? and id in ? and status = 0", claims.UserID, cr.CarIDList)
			err := tx.Model(&carList).
				Where("id in ?", cr.CarIDList).
				Update("status", ctype.GoodsStatusBottom).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("创建订单失败: %v", err)
		res.FailWithMsg("创建订单失败", c)
		return
	}

	data := OrderPayResponse{
		No:     no,
		PayUrl: payUrl,
		Price:  price,
	}

	redisdelay.AddOrderDelay(data.No)

	res.OkWithData(data, c)
}
