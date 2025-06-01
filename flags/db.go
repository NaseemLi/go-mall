package flags

import (
	"fast_gin/global"
	"fast_gin/models"

	"github.com/sirupsen/logrus"
)

func MigrateDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.AddrModel{},
		&models.CouponModel{},
		&models.CollectModel{},
		&models.UserCouponModel{},
		&models.GoodsModel{},
		&models.OrderModel{},
		&models.OrderGoodsModel{},
		&models.LookGoodsModel{},
		&models.MessageModel{},
		&models.CommentModel{},
		&models.SecKillModel{},
		&models.CarModel{},
		&models.OrderCouponModel{},
		&models.Feedback{},
		&models.Item{},
		&models.User{},
	)
	if err != nil {
		logrus.Errorf("表结构迁移失败 %s", err)
		return
	}
	logrus.Infof("表结构迁移成功")
}
