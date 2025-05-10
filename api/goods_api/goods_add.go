package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

// type SecKillRequest struct {
// 	Price     int              `json:"price"`     //秒杀价格
// 	StartTime *ctype.LocalTime `json:"startTime"` //秒杀开始时间
// 	EndTime   *ctype.LocalTime `json:"endTime"`   //秒杀结束时间
// }

// type CouponRequest struct {
// 	Num         int              `json:"num"`         //优惠券数量
// 	Threshold   int              `json:"threshold"`   //满多少元可用
// 	CouponPrice int              `json:"couponPrice"` //优惠券价格
// 	StartTime   *ctype.LocalTime `json:"startTime"`   //秒杀开始时间
// 	EndTime     *ctype.LocalTime `json:"endTime"`     //秒杀结束时间
// }

type GoodsAddRequest struct {
	Title           string               `json:"title" binding:"required,max=64"` //商品名称
	VideoPath       *string              `json:"videoPath"`
	Images          []string             `json:"images" binding:"required"`      //主图
	Price           int                  `json:"price" binding:"required,min=1"` //价格单位:分
	Inventory       *int                 `json:"inventory"`                      //库存
	Category        string               `json:"category"`                       //分类
	Abstract        string               `json:"abstract"`                       //商品简介
	GoodsConfigList []models.GoodsConfig `json:"goodsConfigList"`                //商品配置
	// Seckill         bool                 `json:"seckill"`                        //是否参与秒杀
	// SecKillInfo     *SecKillRequest      `json:"secKillInfo"`                    //秒杀信息
	// Coupon          bool                 `json:"coupon"`                         //是否开启优惠券
	// CouponInfo      *CouponRequest       `json:"couponInfo"`                     //优惠券信息
}

func (GoodsApi) GoodsAddView(c *gin.Context) {
	cr := middleware.GetBind[GoodsAddRequest](c)

	//主图只能九张
	if len(cr.Images) > 9 {
		res.FailWithMsg("商品主图只能九张", c)
		return
	}
	//名称不可以重复
	var model models.GoodsModel
	err := global.DB.Take(&model, "title = ?", cr.Title).Error
	if err == nil {
		res.FailWithMsg("商品名称已存在", c)
		return
	}

	// //判断是否需要将商品加入到秒杀
	// if cr.Seckill {
	// 	//如果需要秒杀,有些参数必填
	// 	if cr.SecKillInfo == nil {
	// 		res.FailWithMsg("秒杀信息不能为空", c)
	// 		return
	// 	}
	// 	if cr.SecKillInfo.Price <= 0 {
	// 		res.FailWithMsg("秒杀价格必须大于0", c)
	// 		return
	// 	}
	// 	if cr.SecKillInfo.StartTime == nil {
	// 		res.FailWithMsg("秒杀开始时间不能为空", c)
	// 		return
	// 	}
	// 	if cr.SecKillInfo.EndTime == nil {
	// 		res.FailWithMsg("秒杀结束时间不能为空", c)
	// 		return
	// 	}

	// 	now := time.Now()
	// 	sub := time.Time(*cr.SecKillInfo.EndTime).Sub(now)
	// 	if sub.Seconds() <= 0 {
	// 		res.FailWithMsg("秒杀结束时间应该大于当前时间", c)
	// 		return
	// 	}
	// 	sub = time.Time(*cr.SecKillInfo.EndTime).Sub(time.Time(*cr.SecKillInfo.StartTime))
	// 	if sub.Seconds() <= 0 {
	// 		res.FailWithMsg("秒杀结束时间应该大于开始时间", c)
	// 		return
	// 	}
	// }

	// if cr.Coupon {
	// 	if cr.CouponInfo == nil {
	// 		res.FailWithMsg("优惠券信息不能为空", c)
	// 		return
	// 	}
	// 	if cr.CouponInfo.Num <= 0 {
	// 		res.FailWithMsg("优惠券数量必须大于0", c)
	// 		return
	// 	}
	// 	if cr.CouponInfo.CouponPrice <= 0 {
	// 		res.FailWithMsg("优惠券金额必须大于0", c)
	// 		return
	// 	}
	// }

	model = models.GoodsModel{
		Title:           cr.Title,
		VideoPath:       cr.VideoPath,
		Images:          cr.Images,
		Price:           cr.Price,
		Inventory:       cr.Inventory,
		Category:        cr.Category,
		Abstract:        cr.Abstract,
		GoodsConfigList: cr.GoodsConfigList,
		Status:          ctype.GoodsStatusTop,
		// SecKill:         cr.Seckill,
		// Coupon:          cr.Coupon,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("商品添加失败", c)
		return
	}
	// //判断是否需要将商品加入到秒杀
	// if cr.Seckill {
	// 	global.DB.Create(&models.SecKillModel{
	// 		GoodsID:   model.ID,
	// 		KillPrice: cr.SecKillInfo.Price,
	// 		StartTime: time.Time(*cr.SecKillInfo.StartTime),
	// 		EndTime:   time.Time(*cr.SecKillInfo.EndTime),
	// 	})
	// }
	// //是否需要创建优惠卷
	// if cr.Coupon {
	// 	co := models.CouponModel{
	// 		Title:       cr.Title,
	// 		Type:        ctype.CouponGoodsType,
	// 		CouponPrice: cr.CouponInfo.CouponPrice,
	// 		Threshold:   cr.CouponInfo.Threshold,
	// 		Num:         cr.CouponInfo.Num,
	// 		GoodsID:     &model.ID,
	// 	}
	// 	if cr.CouponInfo.StartTime != nil {
	// 		startTime := time.Time(*cr.CouponInfo.StartTime)
	// 		co.StartTime = &startTime
	// 	}
	// 	if cr.CouponInfo.EndTime != nil {
	// 		endTime := time.Time(*cr.CouponInfo.EndTime)
	// 		co.EndTime = &endTime
	// 	}
	// 	global.DB.Create(&co)
	// }

	res.Ok(model.ID, "商品添加成功", c)
}
