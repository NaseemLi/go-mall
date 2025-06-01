package commentapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type GoodsCommentListRequest struct {
	models.PageInfo
	GoodsID     uint  `form:"goodsID" binding:"required"` // 商品ID
	CommentType int8  `form:"commentType"`                // 评论类型 0:全部 1:好评 2:中评 3:差评
	IsImages    *bool `form:"isImages"`                   // 是否有图片
}

type GoodsCommentListResponse struct {
	Content      string    `json:"content"`      // 评论内容
	Images       []string  `json:"images"`       // 评论图片
	Level        int8      `json:"level"`        // 评论等级
	UserNickname string    `json:"userNickname"` // 用户昵称
	UserAvatar   string    `json:"userAvatar"`   // 用户头像
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
}

func (CommentApi) GoodsCommentListView(c *gin.Context) {
	cr := middleware.GetBind[GoodsCommentListRequest](c)
	query := global.DB.Where("")
	if cr.IsImages != nil {
		if *cr.IsImages {
			query.Where("images != '[]'")
		} else {
			query.Where("images = '[]'")
		}
	}
	switch cr.CommentType {
	case 1:
		query.Where("level in ?", []int8{4, 5}) // 好评
	case 2:
		query.Where("level in ?", []int8{3}) // 中评
	case 3:
		query.Where("level in ?", []int8{1, 2}) // 差评
	}

	_list, count, _ := common.QueryList(models.CommentModel{
		GoodsID: cr.GoodsID,
	}, common.QueryOption{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: []string{"UserModel"},
	})

	var list = make([]GoodsCommentListResponse, 0)
	for _, item := range _list {
		list = append(list, GoodsCommentListResponse{
			Content:      item.Content,
			Images:       item.Images,
			Level:        item.Level,
			UserNickname: item.UserModel.Nickname,
			UserAvatar:   item.UserModel.Avatar,
			CreatedAt:    item.CreatedAt,
		})
	}
	res.OkWithList(list, count, c)
}
