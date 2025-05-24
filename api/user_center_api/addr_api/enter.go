package addrapi

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/models"
	"fast_gin/service/common"
	"fast_gin/utils/res"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AddrApi struct {
}

type AddrCreateRequest struct {
	Name       string `json:"name" binding:"required,max=16"`
	Tel        string `json:"tel" binding:"required,max=16"`
	Addr       string `json:"addr" binding:"required,max=32"`
	DetailAddr string `json:"detailAddr" binding:"required,max=64"`
}

func (AddrApi) AddrCreateView(c *gin.Context) {
	cr := middleware.GetBind[AddrCreateRequest](c)
	user, err := middleware.GetUser(c)
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	//四个信息不能完全一样
	var model models.AddrModel
	err = global.DB.Take(&model, "name = ? AND tel = ? AND addr = ? AND detail_addr = ?", cr.Name, cr.Tel, cr.Addr, cr.DetailAddr).Error
	if err == nil {
		res.FailWithMsg("地址已存在", c)
		return
	}

	//如果是第一次插入就是默认地址
	addrModel := models.AddrModel{
		UserID:     user.ID,
		Name:       cr.Name,
		Tel:        cr.Tel,
		Addr:       cr.Addr,
		DetailAddr: cr.DetailAddr,
	}
	err = global.DB.Take(&model, "user_id = ?", user.ID).Error
	if err == nil {
		model.IsDefault = true
	}

	err = global.DB.Create(&addrModel).Error
	if err != nil {
		res.FailWithMsg("地址创建失败", c)
		return
	}

	res.OkWithMsg("地址创建成功", c)
}

func (AddrApi) AddrListView(c *gin.Context) {
	var cr = middleware.GetBind[models.PageInfo](c)

	claims := middleware.GetAuth(c)

	list, count, _ := common.QueryList(models.AddrModel{
		UserID: claims.UserID,
	}, common.QueryOption{
		PageInfo: cr,
	})

	res.OkWithList(list, count, c)
}

type AddrUpdateRequest struct {
	ID         uint   `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required,max=16"`
	Tel        string `json:"tel" binding:"required,max=16"`
	Addr       string `json:"addr" binding:"required,max=32"`
	DetailAddr string `json:"detailAddr" binding:"required,max=64"`
}

func (AddrApi) AddrUpdateView(c *gin.Context) {
	cr := middleware.GetBind[AddrUpdateRequest](c)
	user, err := middleware.GetUser(c)
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	var model models.AddrModel
	err = global.DB.Take(&model, "id = ? AND user_id = ?", cr.ID, user.ID).Error
	if err != nil {
		res.FailWithMsg("地址不存在", c)
		return
	}

	//完全没改
	if model.Name == cr.Name && model.Tel == cr.Tel && model.Addr == cr.Addr && model.DetailAddr == cr.DetailAddr {
		res.FailWithMsg("收货地址未修改", c)
		return
	}

	fmt.Println("test")

	//改了之后和之前的一样
	var _model models.AddrModel
	err = global.DB.Take(&_model, "name = ? AND tel = ? AND addr = ? AND detail_addr = ? and id != ?", cr.Name, cr.Tel, cr.Addr, cr.DetailAddr, cr.ID).Error
	if err == nil {
		res.FailWithMsg("此次修改与之前配置重复", c)
		return
	}

	addrModel := models.AddrModel{
		Name:       cr.Name,
		Tel:        cr.Tel,
		Addr:       cr.Addr,
		DetailAddr: cr.DetailAddr,
	}

	err = global.DB.Model(&model).Updates(addrModel).Error
	if err != nil {
		res.FailWithMsg("地址修改失败", c)
		return
	}

	res.OkWithMsg("地址修改成功", c)
}
