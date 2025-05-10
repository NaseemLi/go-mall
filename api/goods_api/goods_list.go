package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsListRequest struct {
	models.PageInfo
	Category string `form:"category"` //分类
}

type GoodsListResponse struct {
	models.GoodsModel
}

func (GoodsApi) GoodsListView(c *gin.Context) {
	var cr = middleware.GetBind[GoodsListRequest](c)

	query := global.DB.Where("")
	_list, count, _ := common.QueryList(models.GoodsModel{
		Category: cr.Category,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Where:    query,
	})

	var list = make([]GoodsListResponse, 0)

	for _, item := range _list {
		list = append(list, GoodsListResponse{
			GoodsModel: item,
		})
	}

	res.OkWithList(list, count, c)
}
