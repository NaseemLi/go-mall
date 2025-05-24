package carapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"
	"fast_gin/utils/set"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (CarApi) CarToCollectView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	claims := middleware.GetAuth(c)
	var list []models.CarModel
	var goodsIDList []uint
	global.DB.Find(&list, "user_id = ? and id in ?", claims.UserID, cr.IDList).Select("goods_id").Scan(&goodsIDList)
	if len(list) > 0 {
		//算出在收藏里面的 id 列表
		var collectIDList []uint
		global.DB.Model(&models.CollectModel{}).Where("user_id = ?", claims.UserID).Pluck("goods_id", &collectIDList)

		// 计算需要添加到收藏的商品 ID 列表
		addCollectIDList := set.DiffArray(goodsIDList, collectIDList)
		var goodsList []models.GoodsModel
		global.DB.Find(&goodsList, "id in ?", addCollectIDList)

		var collectModels []models.CollectModel
		for _, v := range goodsList {
			collectModels = append(collectModels, models.CollectModel{
				UserID:     claims.UserID,
				GoodsID:    v.ID,
				GoodsTitle: v.Title,
			})
		}
		if len(collectModels) > 0 {
			global.DB.Create(&collectModels)
		}
		global.DB.Delete(&list)
	}
	msg := fmt.Sprintf("移入收藏成功,共%d个", len(list))

	res.OkWithMsg(msg, c)
}
