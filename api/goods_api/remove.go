package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserRemoveRequest struct {
}

func (GoodsApi) GoodsRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var goodsList []models.GoodsModel
	global.DB.Find(&goodsList, "id in ?", cr.IDList)
	if len(goodsList) > 0 {
		global.DB.Delete(&goodsList)
	}

	msg := fmt.Sprintf("商品删除成功,删除了 %d 个商品", len(goodsList))

	res.OkWithMsg(msg, c)
}
