package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CarCreateRequest struct {
	GoodsID uint `json:"goodsId" binding:"required"`
	Num     int  `json:"num" binding:"required,min=1,max=999"` //数量
}

var mutex sync.Mutex

func (CarApi) CarCreateView(c *gin.Context) {
	cr := middleware.GetBind[CarCreateRequest](c)
	user, err := middleware.GetUser(c)
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var goods models.GoodsModel
	err = global.DB.Take(&goods, cr.GoodsID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	//商品本身就在购物车中,加数量
	var car models.CarModel
	err = global.DB.Take(&car, "user_id = ? AND goods_id = ?", user.ID, cr.GoodsID).Error
	if err == nil {
		//加数量
		global.DB.Model(&car).Update("num", gorm.Expr("num + ?", cr.Num))
		res.OkWithMsg("商品加入购物车成功", c)
		return
	}

	err = global.DB.Create(&models.CarModel{
		UserID:  user.ID,
		GoodsID: cr.GoodsID,
		Price:   goods.Price,
		Num:     cr.Num,
	}).Error
	if err != nil {
		res.FailWithMsg("加入购物车失败", c)
		return
	}

	res.OkWithMsg("商品加入购物车成功", c)
}
