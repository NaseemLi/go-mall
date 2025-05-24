package collectapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CollectApi struct {
}

type CollectGoodsRequest struct {
	GoodsID uint `json:"goodsID"`
}

type CollectGoodsListResponse struct {
	models.CollectModel
	Cover    string `json:"cover"`
	Title    string `json:"title"`
	Price    int    `json:"price"`
	SalesNum int    `json:"salesNum"` // 销量

}

func (CollectApi) CollectGoodsView(c *gin.Context) {
	cr := middleware.GetBind[CollectGoodsRequest](c)

	var goods models.GoodsModel
	err := global.DB.Take(&goods, cr.GoodsID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	claims := middleware.GetAuth(c)

	var model models.CollectModel
	err = global.DB.Take(&model, "user_id = ? and goods_id = ?", claims.UserID, cr.GoodsID).Error
	if err == nil {
		// 取消收藏
		res.OkWithMsg("取消收藏成功", c)
		global.DB.Delete(&model)
		return
	}

	global.DB.Create(&models.CollectModel{
		GoodsTitle: goods.Title,
		UserID:     claims.UserID,
		GoodsID:    cr.GoodsID,
	})

	res.OkWithMsg("商品收藏成功", c)
}

func (CollectApi) CollectGoodsListView(c *gin.Context) {
	var cr = middleware.GetBind[models.PageInfo](c)

	claims := middleware.GetAuth(c)

	_list, count, _ := common.QueryList(models.CollectModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: cr,
		Preloads: []string{"GoodsModel"},
		Likes:    []string{"goods_title"},
	})
	var list = make([]CollectGoodsListResponse, 0)
	for _, model := range _list {
		list = append(list, CollectGoodsListResponse{
			CollectModel: model,
			Cover:        model.GoodsModel.Images[0],
			Title:        model.GoodsModel.Title,
			Price:        model.GoodsModel.Price,
			SalesNum:     model.GoodsModel.SalesNum,
		})
	}
	res.OkWithList(list, count, c)
}

func (CollectApi) CollectRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)

	claims := middleware.GetAuth(c)
	var list []models.CollectModel
	global.DB.Find(&list, "user_id = ? and id in ?", claims.UserID, cr.IDList)
	if len(list) > 0 {
		global.DB.Delete(&list)
	}
	msg := fmt.Sprintf("收藏商品删除成功，共删除%d个", len(list))

	res.OkWithMsg(msg, c)
}
