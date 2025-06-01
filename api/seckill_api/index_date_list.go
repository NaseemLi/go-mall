package seckillapi

import (
	"context"
	"fast_gin/global"
	"fast_gin/utils/res"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type IndexDateListResponse struct {
	Date string `json:"date"`
}

func (SecKillApi) IndexDateListView(c *gin.Context) {
	keys := global.Redis.Keys(context.Background(), "sec:goods:*").Val()
	var date []IndexDateListResponse
	var dateList []time.Time
	for _, item := range keys {
		_list := strings.Split(item, ":")
		date := _list[len(_list)-1]
		dateObj, err := time.Parse("2006-01-02 15", date)
		if err != nil {
			logrus.Warnf("时间格式错误: %v", err)
			continue
		}
		dateList = append(dateList, dateObj)
	}

	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i].Before(dateList[j])
	})
	for _, v := range dateList {
		date = append(date, IndexDateListResponse{
			Date: v.Format("2006-01-02 15:00:00"),
		})
	}

	res.OkWithData(date, c)
}
