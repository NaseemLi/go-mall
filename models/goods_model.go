package models

type GoodsModel struct {
	Model
	Title            string            `gorm:"size:64" json:"title"`
	VideoPath        *string           `gorm:"size:256" json:"videoPath"`
	Images           []string          `gorm:"type:longtext;serializer:json" json:"images"`          //主图
	Price            int               `json:"price"`                                                //价格单位:分
	Inventory        *int              `json:"inventory"`                                            //库存
	Category         string            `json:"category"`                                             //分类
	Abstract         string            `json:"abstract"`                                             //商品简介
	GoodsConfigList  []GoodsConfig     `gorm:"type:longtext;serializer:json" json:"goodsConfigList"` //商品配置
	SecKill          bool              `json:"secKill"`                                              //是否参与秒杀
	Coupon           bool              `json:"coupon"`                                               //是否开启优惠券
	Status           int8              `json:"status"`                                               //商品状态
	LookCount        int               `json:"lookCount"`                                            //浏览量
	CommentCount     int               `json:"commentCount"`                                         //评论数
	UserBuyGoodsList []OrderGoodsModel `gorm:"foreignKey:GoodsID" json:"-"`                          //购买的次数
}

type GoodsConfig struct {
	Title   string           `json:"title"`   //描述
	SubList []GoodsSubConfig `json:"subList"` //子配置
}

type GoodsSubConfig struct {
	Title string `json:"title"`
	Image string `json:"image"`
}
