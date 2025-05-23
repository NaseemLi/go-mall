package middleware

import (
	"errors"
	"fast_gin/global"
	"fast_gin/models"
	"fast_gin/models/ctype"
	"fast_gin/service/redis_ser"
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.CheckToken(token)
	if err != nil {
		res.FailWithMsg("认证失败", c)
		c.Abort()
		return
	}
	if redis_ser.HasLogout(token) {
		res.FailWithMsg("当前登录已注销", c)
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.CheckToken(token)
	if err != nil {
		res.FailWithMsg("认证失败", c)
		c.Abort()
		return
	}
	if redis_ser.HasLogout(token) {
		res.FailWithMsg("当前登录已注销", c)
		c.Abort()
		return
	}

	if claims.RoleID != ctype.AdminRole {
		res.FailWithMsg("角色认证失败", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Next()
}

func GetAuth(c *gin.Context) (cl *jwts.MyClaims) {
	cl = new(jwts.MyClaims)
	_claims, ok := c.Get("claims")
	if !ok {
		return
	}
	cl, ok = _claims.(*jwts.MyClaims)
	if !ok {
		return
	}
	return cl
}

func GetUser(c *gin.Context) (user models.UserModel, err error) {
	claims := GetAuth(c)
	if claims == nil || claims.UserID == 0 {
		err = errors.New("用户信息非法")
		return
	}
	err = global.DB.First(&user, claims.UserID).Error
	return
}
