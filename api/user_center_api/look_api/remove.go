package lookapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CouponRemoveRequest struct {
}

func (LookApi) LookRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var lookList []models.LookGoodsModel
	global.DB.Find(&lookList, "id in ?", cr.IDList)
	if len(lookList) > 0 {
		global.DB.Delete(&lookList)
	}

	msg := fmt.Sprintf("浏览记录删除成功,删除了 %d 个浏览记录", len(lookList))

	res.OkWithMsg(msg, c)
}
