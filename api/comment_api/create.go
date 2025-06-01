package commentapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CommentOrderGoodsRequest struct {
	OrderGoodsID uint     `json:"orderGoodsID" binding:"required"` //订单商品ID
	Comment      string   `json:"comment"`
	Level        int8     `json:"level" binding:"required,min=1,max=5"` //1-5
	Images       []string `json:"images"`
}

type CommentCreateRequest struct {
	List []CommentOrderGoodsRequest `json:"list" binding:"required,dive"`
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
	var orderGoodsMap = make(map[uint]models.OrderGoodsModel)

	for _, v := range orderGoodsList {
		if v.OrderID != firstOrder.ID {
			res.FailWithMsg("商品不属于同一订单", c)
			return
		}
		orderGoodsMap[v.ID] = v
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
			OrderID:      firstOrder.ID,
			GoodsID:      orderGoodsMap[v.OrderGoodsID].GoodsID,
		})
	}

	if err := global.DB.Create(&list).Error; err != nil {
		res.FailWithMsg("创建评价失败", c)
		return
	}

	//如果这个订单下的商品都评论过了,则改变订单状态
	var commentList []models.CommentModel
	global.DB.Find(&commentList, "order_id = ?", firstOrder.ID)

	if err := global.DB.Preload("OrderGoodsList").
		First(&firstOrder, firstOrder.ID).Error; err != nil {
		res.FailWithMsg("查询订单失败", c)
		return
	}

	//如果评论数量等于订单商品数量,则修改订单状态为已完成
	if len(commentList) == len(firstOrder.OrderGoodsList) {
		global.DB.Model(&firstOrder).Update("status", 5)
		logrus.Infof("订单 %d 全部商品已评价, 改变订单状态为已完成", firstOrder.ID)
	}

	res.OkWithMsg("创建评价成功", c)
}
