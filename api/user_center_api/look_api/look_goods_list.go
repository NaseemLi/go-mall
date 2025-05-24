package lookapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type LookGoodsListRequest struct {
	models.PageInfo
}

type LookGoodsListResponse struct {
	models.LookGoodsModel
	Cover    string `json:"cover"`
	Title    string `json:"title"`
	Price    int    `json:"price"`
	SalesNum int    `json:"salesNum"` // 销量

}

func (LookApi) LookGoodsListView(c *gin.Context) {
	var cr = middleware.GetBind[models.PageInfo](c)

	claims := middleware.GetAuth(c)

	_list, count, _ := common.QueryList(models.LookGoodsModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: cr,
		Preloads: []string{"GoodsModel"},
		Likes:    []string{"goods_title"},
	})
	var list = make([]LookGoodsListResponse, 0)
	for _, model := range _list {
		list = append(list, LookGoodsListResponse{
			LookGoodsModel: model,
			Cover:          model.GoodsModel.Images[0],
			Title:          model.GoodsModel.Title,
			Price:          model.GoodsModel.Price,
			SalesNum:       model.GoodsModel.SalesNum,
		})
	}
	res.OkWithList(list, count, c)
}
