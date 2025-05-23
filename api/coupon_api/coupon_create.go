package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CouponCreateRequest struct {
	Title       string           `json:"title"`                          //优惠券名称
	Type        ctype.CouponType `json:"type" binding:"required"`        //优惠券类型
	CouponPrice int              `json:"couponPrice" binding:"required"` //优惠券金额
	Threshold   int              `json:"threshold"`                      //使用门槛
	Validity    int              `json:"validity" binding:"required"`    //有效期 单位小时
	Num         int              `json:"num" binding:"required"`         //优惠卷数量
	GoodsID     *uint            `json:"goodsID"`                        //关联的商品
	Festival    *string          `json:"festival"`                       //节日活动
}

func (CouponApi) CouponCreateView(c *gin.Context) {
	cr := middleware.GetBind[CouponCreateRequest](c)

	switch cr.Type {
	case ctype.CouponFestivalType:
		if cr.Festival == nil || *cr.Festival == "" {
			res.FailWithMsg("节日优惠卷必须输入节日名称", c)
			return
		}
	case ctype.CouponGoodsType:
		if cr.GoodsID == nil || *cr.GoodsID == 0 {
			res.FailWithMsg("商品优惠卷必须选择商品", c)
			return
		}
	}

	err := global.DB.Create(&models.CouponModel{
		Title:       cr.Title,
		Type:        cr.Type,
		CouponPrice: cr.CouponPrice,
		Threshold:   cr.Threshold,
		Validity:    cr.Validity,
		Num:         cr.Num,
		GoodsID:     cr.GoodsID,
		Festival:    cr.Festival,
	}).Error
	if err != nil {
		res.FailWithMsg("创建优惠卷失败", c)
		return
	}

	res.OkWithMsg("创建优惠卷成功", c)
}
