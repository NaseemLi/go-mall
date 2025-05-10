package goodsapi

import (
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsOptionsListResponse struct {
	ID    uint   `json:"id"`
	Cover string `json:"cover"`
	Title string `json:"title"`
}

func (GoodsApi) GoodsOptionsListView(c *gin.Context) {
	var list = make([]GoodsOptionsListResponse, 0)
	var goods []models.GoodsModel
	global.DB.Find(&goods, "status = ? ", ctype.GoodsStatusTop)

	for _, v := range goods {
		list = append(list, GoodsOptionsListResponse{
			ID:    v.ID,
			Cover: v.Images[0],
			Title: v.Title,
		})
	}

	res.OkWithData(list, c)
}
