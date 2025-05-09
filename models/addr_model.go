package models

type AddrModel struct {
	Model
	UserID     uint      `json:"userID"` //用户ID
	UserModel  UserModel `gorm:"foreignKey:UserID" json:"-"`
	Name       string    `gorm:"size:16" json:"name"`
	Tel        string    `gorm:"size:16" json:"tel"`
	Addr       string    `gorm:"size:32" json:"addr"`
	DetailAddr string    `gorm:"size:64" json:"detailAddr"` //详细地址
	IsDefault  bool      `json:"isDefault"`                 //是否默认地址
}
