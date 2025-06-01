package flags

import (
	"fast_gin/global"
	"fast_gin/models"

	"github.com/sirupsen/logrus"
)

type Gorse struct {
}

func (Gorse) Sync() {
	if global.Gorse == nil {
		logrus.Fatalf("Gorse服务未配置,无法同步用户数据")
		return
	}
	//同步用户
	var myUserList []models.UserModel
	global.DB.Find(&myUserList)
	var userList []models.User
	global.DB.Find(&userList)

	if len(myUserList) != len(userList) {
		logrus.Infof("同步用户数据 %d 条", len(myUserList))
		for _, v := range myUserList {
			v.AfterCreate(nil)
		}
	}

	//同步商品
	var myGoodsList []models.GoodsModel
	global.DB.Find(&myGoodsList)
	var goodsList []models.Item
	global.DB.Find(&goodsList)

	if len(myGoodsList) != len(goodsList) {
		logrus.Infof("同步商品数据 %d 条", len(myGoodsList))
		for _, v := range myGoodsList {
			v.AfterCreate(nil)
		}
	}

	var lookList []models.LookGoodsModel
	global.DB.Find(&lookList)
	var collectList []models.CollectModel
	global.DB.Find(&collectList)

	var feedbackList []models.Feedback
	global.DB.Find(&feedbackList)
	if len(lookList)+len(collectList) > 0 && len(feedbackList) == 0 {
		for _, model := range lookList {
			model.AfterCreate(nil)
		}
		logrus.Infof("同步浏览记录 %d 条", len(lookList))
		for _, model := range collectList {
			model.AfterCreate(nil)
		}
		logrus.Infof("同步收藏记录 %d 条", len(collectList))
	}

}
