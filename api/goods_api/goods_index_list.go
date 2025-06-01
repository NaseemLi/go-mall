package goodsapi

import (
	"context"
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/common"
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"
	"fast_gin/utils/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	if cr.Order == "" {
		//猜你喜欢
		token := c.GetHeader("token")
		claims, err := jwts.CheckToken(token)
		if err == nil && claims != nil && global.Gorse != nil {
			//有 gorse 的同时, 用户也进行登录了
			if cr.Page <= 0 {
				cr.Page = 1
			}
			if cr.Limit <= 0 {
				cr.Limit = -1
			}

			offset := (cr.Page - 1) * cr.Limit
			itemList, err := global.Gorse.GetRecommendOffSet(context.Background(),
				fmt.Sprintf("%d", claims.UserID), cr.Category, cr.Limit, offset)
			if err != nil {
				logrus.Errorf("获取推荐商品失败: %v", err)
				res.FailWithMsg("推荐系统异常", c)
				return
			}
			logrus.Infof("获取推荐商品成功: %v", itemList)
			cr.Order = sql.OrderRevert(itemList, "id")
			logrus.Infof("获取推荐商品排序: %s", cr.Order)
		}
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
