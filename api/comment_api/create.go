package commentapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type CommentOrderGoodsRequest struct {
	OrderGoodsID uint     `json:"orderGoodsID"`
	Comment      string   `json:"comment"`
	Level        int8     `json:"level"` //1-5
	Images       []string `json:"images"`
}

type CommentCreateRequest struct {
	List []CommentOrderGoodsRequest `json:"list" binding:"required"`
}

func (CommentApi) CommentCreateView(c *gin.Context) {
	cr := middleware.GetBind[CommentCreateRequest](c)
	claims := middleware.GetAuth(c)

	var orderGoodsIDList []uint
	for _, info := range cr.List {
		orderGoodsIDList = append(orderGoodsIDList, info.OrderGoodsID)
	}

	var orderGoodsList []models.OrderGoodsModel
	global.DB.Preload("OrderModel").Find(&orderGoodsList, "id IN ? AND user_id = ?",
		orderGoodsIDList, claims.UserID)
	if len(orderGoodsList) == 0 {
		res.FailWithMsg("未找到相关订单商品", c)
		return
	}

	//判断是否归属一个订单
	firstOrder := orderGoodsList[0].OrderModel
	for _, v := range orderGoodsList {
		if v.OrderID != firstOrder.ID {
			res.FailWithMsg("商品不属于同一订单", c)
			return
		}
	}

	if firstOrder.Status != 4 {
		res.FailWithMsg("订单状态必须为已收货/待评价", c)
		return
	}

	//重复判断
	var list []models.CommentModel
	global.DB.Find(&list, "order_goods_id in ?", orderGoodsIDList)
	if len(list) > 0 {
		res.FailWithMsg("已评价过的商品不能重复评价", c)
		return
	}

	for _, v := range cr.List {
		list = append(list, models.CommentModel{
			UserID:       claims.UserID,
			OrderGoodsID: v.OrderGoodsID,
			Content:      v.Comment,
			Level:        v.Level,
			Images:       v.Images,
		})
	}

	if err := global.DB.Create(&list).Error; err != nil {
		res.FailWithMsg("创建评价失败", c)
		return
	}

	res.OkWithMsg("创建评价成功", c)
}
