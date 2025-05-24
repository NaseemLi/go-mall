package goodsapi

import (
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsIndexListRequest struct {
	models.PageInfo
	Category string `form:"category"`
}

type GoodsIndexListResponse struct {
	ID       uint   `json:"id"`
	Cover    string `json:"cover"`
	Title    string `json:"title"`
	Price    int    `json:"price"`
	SalesNum int    `json:"salesNum"`
}

func (GoodsApi) GoodsIndexListView(c *gin.Context) {
	var cr = middleware.GetBind[GoodsIndexListRequest](c)

	sortMap := map[string]bool{
		"":               true,
		"price desc":     true,
		"price asc":      true,
		"sales_num desc": true,
		"sales_num asc":  true,
	}

	_, ok := sortMap[cr.Order]
	if !ok {
		res.FailWithMsg("排序方式错误", c)
		return
	}

	_list, count, _ := common.QueryList(models.GoodsModel{
		Category: cr.Category,
		Status:   ctype.GoodsStatusTop,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
	})

	var list = make([]GoodsIndexListResponse, 0)

	for _, item := range _list {
		list = append(list, GoodsIndexListResponse{
			ID:       item.ID,
			Cover:    item.Images[0],
			Title:    item.Title,
			Price:    item.Price,
			SalesNum: item.SalesNum,
		})
	}

	res.OkWithList(list, count, c)
}
