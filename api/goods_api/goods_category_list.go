package goodsapi

import (
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

func (GoodsApi) GoodsCategoryListView(c *gin.Context) {
	var list []models.LabelResponse[string]
	global.DB.Model(&models.GoodsModel{}).
		Where("status = ? and category != ''", ctype.GoodsStatusTop).
		Group("category").
		Select("category as label", "category as value").
		Scan(&list)

	res.OkWithData(list, c)
}
