package models

import (
	"fast_gin/models/ctype"
	"fmt"

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
	//todo:联动删除关联信息
	fmt.Printf("删除用户%s 之前\n", u.Username)
	return nil
}
