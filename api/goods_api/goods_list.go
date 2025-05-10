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
	SecKill  *bool  `form:"secKill"`  //是否参与秒杀
	Category string `form:"category"` //分类
}

type GoodsListResponse struct {
	models.GoodsModel
	BuyUserNum int `json:"buyUserNum"` //购买人数

}

func (GoodsApi) GoodsListView(c *gin.Context) {
	var cr = middleware.GetBind[GoodsListRequest](c)

	query := global.DB.Where("")
	if cr.SecKill != nil {
		query = query.Where("sec_kill = ?", *cr.SecKill)
	}
	_list, count, _ := common.QueryList(models.GoodsModel{
		Category: cr.Category,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: []string{"UserBuyGoodsList"},
		Where:    query,
	})

	var list = make([]GoodsListResponse, 0)

	for _, item := range _list {
		list = append(list, GoodsListResponse{
			GoodsModel: item,
			BuyUserNum: len(item.UserBuyGoodsList),
		})
	}

	res.OkWithList(list, count, c)
}
