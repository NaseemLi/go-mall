package couponapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserCouponListResponse struct {
	models.Model
	Title       string             `json:"title"`       //优惠券名称
	Type        ctype.CouponType   `json:"type"`        //优惠券类型
	CouponPrice int                `json:"couponPrice"` //优惠券金额
	Threshold   int                `json:"threshold"`   //使用门槛
	Validity    int                `json:"validity"`    //有效期 单位小时
	Status      ctype.CouponStatus `json:"status"`      //状态
	CouponID    uint               `json:"couponID"`    //优惠券id
}

type UserCouponListRequest struct {
	models.PageInfo
	Status ctype.CouponStatus `form:"status"` //状态
}

func (CouponApi) UserCouponListView(c *gin.Context) {
	var cr = middleware.GetBind[UserCouponListRequest](c)
	claims := middleware.GetAuth(c)

	//把用户领取过的优惠卷过滤掉
	_list, count, _ := common.QueryList(models.UserCouponModel{
		Status: cr.Status,
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Preloads: []string{"CouponModel"},
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

	var list = make([]UserCouponListResponse, 0)
	var updateList []models.UserCouponModel
	for _, model := range _list {
		list = append(list, UserCouponListResponse{
			Model:       model.Model,
			Title:       model.CouponModel.Title,
			Type:        model.CouponModel.Type,
			CouponPrice: model.CouponModel.CouponPrice,
			Threshold:   model.CouponModel.Threshold,
			Validity:    model.CouponModel.Validity,
			Status:      cr.Status,
			CouponID:    model.CouponID,
		})
		if time.Until(model.EndTime) < 0 && model.Status != ctype.CouponStatusNotUsed {
			//过期,并且状态是未使用
			updateList = append(updateList, model)
			logrus.Infof("优惠券过期: %v", model.ID)
		}
	}
	if len(updateList) > 0 {
		global.DB.Model(&updateList).Update("status", ctype.CouponStatusExpired)
	}

	res.OkWithList(list, count, c)
}
