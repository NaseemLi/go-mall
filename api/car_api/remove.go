package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (CarApi) CarRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	claims := middleware.GetAuth(c)
	var list []models.CarModel
	global.DB.Find(&list, "user_id = ? and id in ?", claims.UserID, cr.IDList)
	if len(list) > 0 {
		global.DB.Delete(&list)
	}
	msg := fmt.Sprintf("购物车删除成功，共删除%d个", len(list))

	res.OkWithMsg(msg, c)
}
