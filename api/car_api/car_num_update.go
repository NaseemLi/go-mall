package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CarNumUpdateRequest struct {
	CarID uint `json:"carID" binding:"required"`
	Num   int  `json:"num" binding:"required,min=1,max=10000"`
}

func (CarApi) CarNumUpdateView(c *gin.Context) {
	cr := middleware.GetBind[CarNumUpdateRequest](c)

	user, err := middleware.GetUser(c)
	if err != nil {
		res.FailWithMsg("用户信息不存在", c)
		return
	}

	var car models.CarModel
	err = global.DB.Preload("GoodsModel").Take(&car, "user_id = ? and id = ?", user.ID, cr.CarID).Error
	if err != nil {
		res.FailWithMsg("购物车记录不存在", c)
		return
	}

	// 判断能不能加数量
	if car.GoodsModel.Inventory != nil {
		if cr.Num > *car.GoodsModel.Inventory {
			res.FailWithMsg("选择数量大于商品库存", c)
			return
		}
	}
	global.DB.Model(&car).Update("num", cr.Num)
	res.OkWithMsg("购物车商品数量修改成功", c)
}
