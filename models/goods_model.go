package models

type GoodsModel struct {
	Model
	Title           string        `gorm:"size:64" json:"title"`
	VideoPath       *string       `gorm:"size:256" json:"videoPath"`
	Images          []string      `gorm:"type:longtext;serializer:json" json:"images"`          //主图
	Price           int           `json:"price"`                                                //价格单位:分
	Inventory       *int          `json:"inventory"`                                            //库存
	Category        string        `json:"category"`                                             //分类
	Abstract        string        `json:"abstract"`                                             //商品简介
	GoodsConfigList []GoodsConfig `gorm:"type:longtext;serializer:json" json:"goodsConfigList"` //商品配置
	Seckill         bool          `json:"seckill"`                                              //是否参与秒杀
	Coupon          bool          `json:"coupon"`                                               //是否开启优惠券
}

type GoodsConfig struct {
	Title   string           `json:"title"`   //描述
	SubList []GoodsSubConfig `json:"subList"` //子配置
}

type GoodsSubConfig struct {
	Title string `json:"title"`
	Image string `json:"image"`
}
