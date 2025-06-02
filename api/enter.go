package api

import (
	"fast_gin/api/captcha_api"
	carapi "fast_gin/api/car_api"
	commentapi "fast_gin/api/comment_api"
	couponapi "fast_gin/api/coupon_api"
	dataapi "fast_gin/api/data_api"
	goodsapi "fast_gin/api/goods_api"
	"fast_gin/api/image_api"
	msgapi "fast_gin/api/msg_api"
	orderapi "fast_gin/api/order_api"
	seckillapi "fast_gin/api/seckill_api"
	"fast_gin/api/user_api"
	usercenterapi "fast_gin/api/user_center_api"
)

type Api struct {
	UserApi       user_api.UserApi
	ImageApi      image_api.ImageApi
	CaptchaApi    captcha_api.CaptchaApi
	GoodsApi      goodsapi.GoodsApi
	CouponApi     couponapi.CouponApi
	UserCenterApi usercenterapi.UserCenterApi
	CarApi        carapi.CarApi
	OrderApi      orderapi.OrderApi
	CommentApi    commentapi.CommentApi
	MsgApi        msgapi.MsgApi
	SecKillApi    seckillapi.SecKillApi
	DataApi       dataapi.DataApi
}

var App = new(Api)
