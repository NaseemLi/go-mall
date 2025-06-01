package api

import (
	"fast_gin/api/captcha_api"
	carapi "fast_gin/api/car_api"
	commentapi "fast_gin/api/comment_api"
	couponapi "fast_gin/api/coupon_api"
	goodsapi "fast_gin/api/goods_api"
	"fast_gin/api/image_api"
	orderapi "fast_gin/api/order_api"
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
}

var App = new(Api)
