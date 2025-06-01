package commentapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CommentLevelListRequest struct {
	GoodsID uint `form:"goodsID" bind:"required"` // 商品ID
}

type CommentLevelListResponse struct {
	AllCount    int `json:"allCount"`    // 全部评论数量
	ImageCount  int `json:"imageCount"`  // 带图片评论数量
	GreatCount  int `json:"greatCount"`  // 好评数量
	MiddleCount int `json:"middleCount"` // 中评数量
	BadCount    int `json:"badCount"`    // 差评数量
}

func (CommentApi) CommentLevelListView(c *gin.Context) {
	cr := middleware.GetBind[CommentLevelListRequest](c)

	var list []models.CommentModel
	var allCount, imageCount, greatCount, middleCount, badCount int
	// 查询评论列表
	global.DB.Find(&list, "goods_id = ?", cr.GoodsID)
	for _, v := range list {
		switch v.Level {
		case 1, 2:
			badCount++
		case 3:
			middleCount++
		case 4, 5:
			greatCount++
		}
		if len(v.Images) > 0 {
			imageCount++
		}
	}

	allCount = len(list)
	data := CommentLevelListResponse{
		AllCount:    allCount,
		ImageCount:  imageCount,
		GreatCount:  greatCount,
		MiddleCount: middleCount,
		BadCount:    badCount,
	}
	res.OkWithData(data, c)
}
