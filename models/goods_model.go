package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GoodsModel struct {
	Model
	Title            string            `gorm:"size:64" json:"title"`
	VideoPath        *string           `gorm:"size:256" json:"videoPath"`
	Images           []string          `gorm:"type:longtext;serializer:json" json:"images"`          //主图
	Price            int               `json:"price"`                                                //价格单位:分
	Inventory        *int              `json:"inventory"`                                            //库存
	Category         string            `json:"category"`                                             //分类
	Abstract         string            `json:"abstract"`                                             //商品简介
	GoodsConfigList  []GoodsConfig     `gorm:"type:longtext;serializer:json" json:"goodsConfigList"` //商品配置
	SecKill          bool              `json:"secKill"`                                              //是否参与秒杀
	Coupon           bool              `json:"coupon"`                                               //是否开启优惠券
	Status           int8              `json:"status"`                                               //商品状态
	LookCount        int               `json:"lookCount"`                                            //浏览量
	CommentCount     int               `json:"commentCount"`                                         //评论数
	UserBuyGoodsList []OrderGoodsModel `gorm:"foreignKey:GoodsID" json:"-"`                          //购买的次数
}

type GoodsConfig struct {
	Title   string           `json:"title"`   //描述
	SubList []GoodsSubConfig `json:"subList"` //子配置
}

type GoodsSubConfig struct {
	Title string `json:"title"`
	Image string `json:"image"`
}

func (g GoodsModel) BeforeDelete(tx *gorm.DB) (err error) {
	// 删除秒杀
	if err = tx.Where("goods_id = ?", g.ID).Delete(&SecKillModel{}).Error; err != nil {
		logrus.Errorf("删除秒杀商品失败: goods_id=%d, err=%v", g.ID, err)
		return err
	}
	logrus.Infof("删除秒杀商品: goods_id=%d", g.ID)

	// 删除优惠券
	if err = tx.Where("goods_id = ?", g.ID).Delete(&CouponModel{}).Error; err != nil {
		logrus.Errorf("删除优惠券失败: goods_id=%d, err=%v", g.ID, err)
		return err
	}
	logrus.Infof("删除优惠券: goods_id=%d", g.ID)

	// 删除购物车
	if err = tx.Where("goods_id = ?", g.ID).Delete(&OrderGoodsModel{}).Error; err != nil {
		logrus.Errorf("删除购物车商品失败: goods_id=%d, err=%v", g.ID, err)
		return err
	}
	logrus.Infof("删除购物车商品: goods_id=%d", g.ID)

	return nil
}
