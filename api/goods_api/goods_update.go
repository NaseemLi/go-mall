package goodsapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type GoodsUpdateRequest struct {
	ID              uint                 `json:"id" binding:"required"`           //商品ID
	Title           string               `json:"title" binding:"required,max=64"` //商品名称
	VideoPath       *string              `json:"videoPath"`
	Images          []string             `json:"images" binding:"required"`      //主图
	Price           int                  `json:"price" binding:"required,min=1"` //价格单位:分
	Inventory       *int                 `json:"inventory"`                      //库存
	Category        string               `json:"category"`                       //分类
	Abstract        string               `json:"abstract"`                       //商品简介
	GoodsConfigList []models.GoodsConfig `json:"goodsConfigList"`                //商品配置
}

func (GoodsApi) GoodsUpdateView(c *gin.Context) {
	cr := middleware.GetBind[GoodsUpdateRequest](c)
	//主图只能九张
	if len(cr.Images) > 9 {
		res.FailWithMsg("商品主图只能九张", c)
		return
	}
	// 更新商品信息
	var model models.GoodsModel
	err := global.DB.Take(&model, "id = ?", cr.ID).Error
	if err != nil {
		res.FailWithMsg("商品不存在", c)
		return
	}

	model.Title = cr.Title
	model.VideoPath = cr.VideoPath
	model.Images = cr.Images
	model.Price = cr.Price
	model.Inventory = cr.Inventory
	model.Category = cr.Category
	model.Abstract = cr.Abstract
	model.GoodsConfigList = cr.GoodsConfigList

	err = global.DB.Save(&model).Error
	if err != nil {
		res.FailWithMsg("商品更新失败", c)
		return
	}

	res.Ok(model.ID, "商品更新成功", c)
}
