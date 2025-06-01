package seckillapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (SecKillApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	var seckillList []models.SecKillModel
	global.DB.Find(&seckillList, "id in ?", cr.IDList)
	if len(seckillList) > 0 {
		global.DB.Delete(&seckillList)
	}

	msg := fmt.Sprintf("秒杀商品删除成功,删除了 %d 个商品", len(seckillList))

	res.OkWithMsg(msg, c)
}
