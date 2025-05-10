package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsStatusUpdateRequest struct {
	ID     uint                  `json:"id"`
	Status ctype.GoodsStatusType `json:"status"`
}

func (GoodsApi) GoodsStatusUpdateView(c *gin.Context) {
	cr := middleware.GetBind[GoodsStatusUpdateRequest](c)

	var model models.GoodsModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}
	global.DB.Model(&model).Update("status", cr.Status)
	//如果商品在购物车中 此时商品下架 直接修改购物车 对应状态
	res.OkWithMsg("商品状态更新成功", c)
}
