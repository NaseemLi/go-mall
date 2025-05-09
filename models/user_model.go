package models

import "fast_gin/models/ctype"

type UserModel struct {
	Model
	Username string     `gorm:"size:16" json:"username"`
	Nickname string     `gorm:"size:32" json:"nickname"`
	Password string     `gorm:"size:64" json:"-"`
	Avatar   string     `json:"avatar"`
	RoleID   ctype.Role `json:"roleID"` // 1 管理员 2 普通用户
}
