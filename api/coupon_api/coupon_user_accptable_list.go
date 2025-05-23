package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CouponUserAcceptableListResponse struct {
	models.Model
	Title       string           `json:"title"`       //优惠券名称
	Type        ctype.CouponType `json:"type"`        //优惠券类型
	CouponPrice int              `json:"couponPrice"` //优惠券金额
	Threshold   int              `json:"threshold"`   //使用门槛
	Validity    int              `json:"validity"`    //有效期 单位小时
	IsReceive   bool             `json:"isReceive"`   //是否领取
}

// 用户可以领取的优惠券列表
func (CouponApi) CouponUserAcceptableListView(c *gin.Context) {
	var cr = middleware.GetBind[models.PageInfo](c)

	query := global.DB.Where("`type` in ?", []ctype.CouponType{
		ctype.CouponFestivalType,
		ctype.CouponGoodsType,
		ctype.CouponGeneralType,
	})

	//把用户领取过的优惠卷过滤掉
	_list, count, _ := common.QueryList(models.CouponModel{}, common.QueryOption{
		PageInfo: cr,
		Where:    query,
	})
	claims, err := jwts.CheckToken(c.GetHeader("token"))
	var reviceMap = make(map[uint]bool)
	if err == nil && claims != nil {
		//用户登录了
		//查这个用户领取过的优惠卷id列表
		var userCouponList []models.UserCouponModel
		global.DB.Find(&userCouponList, "user_id = ?", claims.UserID)
		for _, v := range userCouponList {
			reviceMap[v.CouponID] = true
		}
	}

	var list = make([]CouponUserAcceptableListResponse, 0)
	for _, model := range _list {
		list = append(list, CouponUserAcceptableListResponse{
			Model:       model.Model,
			Title:       model.Title,
			Type:        model.Type,
			CouponPrice: model.CouponPrice,
			Threshold:   model.Threshold,
			Validity:    model.Validity,
			IsReceive:   reviceMap[model.ID],
		})
	}

	res.OkWithList(list, count, c)
}
