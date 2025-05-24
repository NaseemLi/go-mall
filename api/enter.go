package api

import (
	"fast_gin/api/captcha_api"
	couponapi "fast_gin/api/coupon_api"
	goodsapi "fast_gin/api/goods_api"
	"fast_gin/api/image_api"
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
}

var App = new(Api)
