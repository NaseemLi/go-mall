package models

type UserLoginModel struct {
	Model
	UserID uint   `json:"userID"`
	Ip     string `gorm:"size:32" json:"ip"`
	Addr   string `gorm:"size:64" json:"addr"`
	Ua     string `gorm:"size:128" json:"ua"`
}
