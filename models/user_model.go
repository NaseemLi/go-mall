package models

import (
	"fast_gin/models/ctype"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserModel struct {
	Model
	Username string     `gorm:"size:16" json:"username"`
	Nickname string     `gorm:"size:32" json:"nickname"`
	Password string     `gorm:"size:64" json:"-"`
	Avatar   string     `json:"avatar"`
	RoleID   ctype.Role `json:"roleID"` // 1 管理员 2 普通用户
}

func (u UserModel) BeforeDelete(tx *gorm.DB) (err error) {
	//地址表
	var addrList []AddrModel
	tx.Find(&addrList, "user_id = ?", u.ID)
	if len(addrList) > 0 {
		tx.Delete(&addrList)
		logrus.Infof("删除用户%s 的地址表数据", u.Username)
		logrus.Infof("删除多少条地址表数据: %d", len(addrList))
	}
	//收藏表
	var favoriteList []CollectModel
	tx.Find(&favoriteList, "user_id = ?", u.ID)
	if len(favoriteList) > 0 {
		tx.Delete(&favoriteList)
		logrus.Infof("删除用户%s 的收藏表数据", u.Username)
		logrus.Infof("删除多少条收藏表数据: %d", len(favoriteList))
	}
	//购物车表
	var cartList []CarModel
	tx.Find(&cartList, "user_id = ?", u.ID)
	if len(cartList) > 0 {
		tx.Delete(&cartList)
		logrus.Infof("删除用户%s 的购物车表数据", u.Username)
		logrus.Infof("删除多少条购物车表数据: %d", len(cartList))
	}
	//评论表
	var commentList []CommentModel
	tx.Find(&commentList, "user_id = ?", u.ID)
	if len(commentList) > 0 {
		tx.Delete(&commentList)
		logrus.Infof("删除用户%s 的评论表数据", u.Username)
		logrus.Infof("删除多少条评论表数据: %d", len(commentList))
	}
	//订单 订单商品 订单优惠卷表
	var orderList []OrderModel
	tx.Find(&orderList, "user_id = ?", u.ID)
	if len(orderList) > 0 {
		tx.Unscoped().Delete(&orderList)
		logrus.Infof("删除用户%s 的订单表数据", u.Username)
		logrus.Infof("删除多少条订单表数据: %d", len(orderList))
	}

	var orderGoodsList []OrderGoodsModel
	tx.Find(&orderGoodsList, "user_id = ?", u.ID)
	if len(orderGoodsList) > 0 {
		tx.Delete(&orderGoodsList)
		logrus.Infof("删除用户%s 的订单商品表数据", u.Username)
		logrus.Infof("删除多少条订单商品表数据: %d", len(orderGoodsList))
	}

	var orderCouponList []OrderCouponModel
	tx.Find(&orderCouponList, "user_id = ?", u.ID)
	if len(orderCouponList) > 0 {
		tx.Delete(&orderCouponList)
		logrus.Infof("删除用户%s 的订单优惠卷表数据", u.Username)
		logrus.Infof("删除多少条订单优惠卷表数据: %d", len(orderCouponList))
	}

	//用户优惠卷
	var userCouponList []UserCouponModel
	tx.Find(&userCouponList, "user_id = ?", u.ID)
	if len(userCouponList) > 0 {
		tx.Delete(&userCouponList)
		logrus.Infof("删除用户%s 的用户优惠卷表数据", u.Username)
		logrus.Infof("删除多少条用户优惠卷表数据: %d", len(userCouponList))
	}
	//商品浏览
	var browseList []LookGoodsModel
	tx.Find(&browseList, "user_id = ?", u.ID)
	if len(browseList) > 0 {
		tx.Delete(&browseList)
		logrus.Infof("删除用户%s 的商品浏览表数据", u.Username)
		logrus.Infof("删除多少条商品浏览表数据: %d", len(browseList))
	}
	//系统消息
	var messageList []MessageModel
	tx.Find(&messageList, "user_id = ?", u.ID)
	if len(messageList) > 0 {
		tx.Unscoped().Delete(&messageList)
		logrus.Infof("删除用户%s 的系统消息表数据", u.Username)
		logrus.Infof("删除多少条系统消息表数据: %d", len(messageList))
	}
	logrus.Infof("删除用户%s 完成\n", u.Username)
	return nil
}
