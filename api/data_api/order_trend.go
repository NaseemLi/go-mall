package dataapi

import (
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderTrendResponse struct {
	DateList  []string `json:"date"`
	CountList []int    `json:"count"`
}

func (DataApi) OrderTrendView(c *gin.Context) {
	now := time.Now()
	endTime := now.Format("2006-01-02") + " 23:59:59"
	startTime := now.AddDate(0, 0, -7).Format("2006-01-02") + " 00:00:00"
	_startTime := now.AddDate(0, 0, -7)

	var data OrderTrendResponse
	var orderList []models.OrderModel
	global.DB.Find(&orderList, "created_at >= ? AND created_at <= ? and status not in ?",
		startTime, endTime, []int8{1, 6, 7})

	var countMap = make(map[string]int)
	for _, v := range orderList {
		date := v.CreatedAt.Format("2006-01-02")
		count, ok := countMap[date]
		if ok {
			countMap[date] = count + 1
		} else {
			countMap[date] = 1
		}
	}

	for i := 7; i > 0; i-- {
		date := _startTime.AddDate(0, 0, i).Format("2006-01-02")
		data.DateList = append(data.DateList, date)
		data.CountList = append(data.CountList, countMap[date])
	}

	res.OkWithData(data, c)
}
