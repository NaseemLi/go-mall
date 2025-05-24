package lookapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type LookApi struct {
}

type LookGoodsRequest struct {
	GoodsID uint `json:"goodsID"`
}

type LookGoodsResponse struct {
}

func (LookApi) LookGoodsView(c *gin.Context) {
	cr := middleware.GetBind[LookGoodsRequest](c)

	var goods models.GoodsModel
	err := global.DB.Take(&goods, cr.GoodsID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	claims := middleware.GetAuth(c)

	var model models.LookGoodsModel
	err = global.DB.Take(&model, "user_id = ? and goods_id = ? and date(created_at) = date(now())", claims.UserID, cr.GoodsID).Error
	if err == nil {
		res.OkWithMsg("OK", c)
		return
	}

	global.DB.Create(&models.LookGoodsModel{
		UserID:  claims.UserID,
		GoodsID: cr.GoodsID,
	})

	res.OkWithMsg("查看商品成功", c)
}
