package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsDetailCoupon struct {
	CouponPrice int  `json:"couponPrice"` //优惠券金额
	Threshold   int  `json:"threshold"`   //使用门槛
	IsReceive   bool `json:"isReceive"`   //是否领取
}

type GoodsDetailResponse struct {
	models.GoodsModel
	IsCollect         bool               `json:"isCollect"`         //是否收藏
	IsGoodsCoupon     bool               `json:"isGoodsCoupon"`     //是否有优惠券
	GoodsDetailCoupon *GoodsDetailCoupon `json:"goodsDetailCoupon"` //优惠券信息
}

func (GoodsApi) GoodsDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	// 查询商品信息
	var model models.GoodsModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	data := GoodsDetailResponse{
		GoodsModel: model,
	}

	//判断用户是否登录
	token := c.GetHeader("token")
	claims, err := jwts.CheckToken(token)
	//查这个商品是否有优惠券
	if err == nil && claims != nil {
		// 用户登录
		var userCollect models.CollectModel
		err = global.DB.Take(&userCollect, "user_id = ? and goods_id = ?", claims.UserID, model.ID).Error
		if err == nil {
			data.IsCollect = true
		}
	}

	res.OkWithData(data, c)
}
