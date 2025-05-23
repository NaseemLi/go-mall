package user_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/utils/captcha"
	"fast_gin/utils/pwd"
	"fast_gin/utils/random"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,max=16" label:"用户名"`
	Password    string `json:"password" binding:"required,max=64" label:"密码"`
	RePassword  string `json:"rePassword" binding:"required,max=64" label:"确认密码"`
	CaptchaID   string `json:"captchaID"`
	CaptchaCode string `json:"captchaCode"`
}

func (UserApi) RegisterView(c *gin.Context) {
	cr := middleware.GetBind[RegisterRequest](c)

	if global.Config.Site.Login.Captcha {
		if cr.CaptchaID == "" || cr.CaptchaCode == "" {
			res.FailWithMsg("请输入图片验证码", c)
			return
		}
		if !captcha.CaptchaStore.Verify(cr.CaptchaID, cr.CaptchaCode, true) {
			res.FailWithMsg("图片验证码验证失败", c)
			return
		}
	}

	var user models.UserModel
	err := global.DB.Take(&user, "username = ?", cr.Username).Error
	if err == nil {
		res.FailWithMsg("用户名已经存在", c)
		return
	}

	hashPwd := pwd.GenerateFromPassword(cr.Password)

	user = models.UserModel{
		Username: cr.Username,
		Nickname: "注册用户_" + random.RandStr(5),
		Password: hashPwd,
		RoleID:   ctype.UserRole,
	}

	err = global.DB.Create(&user).Error
	if err != nil {
		res.FailWithMsg("用户注册失败", c)
		return
	}

	// //判断是否有新用户优惠卷
	// var couponList []models.CouponModel
	// global.DB.Find(&couponList, "`type` = ? and `receive` != `num`", ctype.CouponNewUserType)
	// if len(couponList) > 0 {
	// 	//给用户发优惠卷
	// 	var userCouponList []models.UserCouponModel
	// 	for _, couponModel := range couponList {
	// 		userCouponList = append(userCouponList, models.UserCouponModel{
	// 			UserID:   user.ID,
	// 			CouponID: couponModel.ID,
	// 			Status:   ctype.CouponStatusNotUsed,
	// 			EndTime:  time.Now().Add(time.Duration(couponModel.Validity) * time.Hour),
	// 		})
	// 	}
	// 	if len(userCouponList) > 0 {
	// 		global.DB.Create(&userCouponList)
	// 		//增加对应新用户优惠卷数量
	// 		global.DB.Model(&couponList).Update("receive", gorm.Expr("receive + 1"))
	// 		logrus.Infof("添加用户优惠卷成功, 创建了 %d 张优惠卷", len(userCouponList))
	// 	}
	// }
	res.OkWithMsg("注册成功", c)
}
