package dataapi

import (
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginTrendResponse struct {
	DateList  []string `json:"date"`
	LoginList []int    `json:"login"`
	SignList  []int    `json:"sign"`
}

func (DataApi) UserLoginTrendView(c *gin.Context) {
	now := time.Now()
	endTime := now.Format("2006-01-02") + " 23:59:59"
	startTime := now.AddDate(0, 0, -7).Format("2006-01-02") + " 00:00:00"
	_startTime := now.AddDate(0, 0, -7)

	var data UserLoginTrendResponse
	var userLoginList []models.UserLoginModel
	global.DB.Find(&userLoginList, "created_at >= ? AND created_at <= ?", startTime, endTime)

	var loginMap = make(map[string]int)
	for _, v := range userLoginList {
		date := v.CreatedAt.Format("2006-01-02")
		count, ok := loginMap[date]
		if ok {
			loginMap[date] = count + 1
		} else {
			loginMap[date] = 1
		}
	}

	var userSignList []models.UserModel
	global.DB.Find(&userSignList, "created_at >= ? AND created_at <= ?", startTime, endTime)

	var signMap = make(map[string]int)
	for _, v := range userSignList {
		date := v.CreatedAt.Format("2006-01-02")
		count, ok := signMap[date]
		if ok {
			signMap[date] = count + 1
		} else {
			signMap[date] = 1
		}
	}

	for i := 7; i > 0; i-- {
		date := _startTime.AddDate(0, 0, i).Format("2006-01-02")
		data.DateList = append(data.DateList, date)
		data.LoginList = append(data.LoginList, loginMap[date])
		data.SignList = append(data.SignList, signMap[date])
	}

	res.OkWithData(data, c)
}
