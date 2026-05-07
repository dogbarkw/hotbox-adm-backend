package models

import (
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/form"

	"github.com/gin-gonic/gin"
)

type AiMatchProductNftSecondModel struct {
	ID                 int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                         //  主键id
	OrderId            int64     `gorm:"column:order_id" db:"order_id" json:"order_id" form:"order_id"` //  订单id
	SecondType         string    `gorm:"column:second_type" db:"second_type" json:"second_type" form:"second_type"`
	UserId             int64     `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`
	BuyerUserId        int64     `gorm:"column:buyer_user_id" db:"buyer_user_id" json:"buyer_user_id" form:"buyer_user_id"`
	SellerUserId       int64     `gorm:"column:seller_user_id" db:"seller_user_id" json:"seller_user_id" form:"seller_user_id"`
	ReceiverProvince   string    `gorm:"column:receiver_province" db:"receiver_province" json:"receiver_province" form:"receiver_province"`
	ReceiverCity       string    `gorm:"column:receiver_city" db:"receiver_city" json:"receiver_city" form:"receiver_city"`
	ReceiverRegion     string    `gorm:"column:receiver_region" db:"receiver_region" json:"receiver_region" form:"receiver_region"`
	OriginPrice        float64   `gorm:"column:origin_price" db:"origin_price" json:"origin_price" form:"origin_price"` //  原价
	Price              float64   `gorm:"column:price" db:"price" json:"price" form:"price"`                             //  出售价
	ProductId          int64     `gorm:"column:product_id" db:"product_id" json:"product_id" form:"product_id"`         //  产品id
	SizeId             int64     `gorm:"column:size_id" db:"size_id" json:"size_id" form:"size_id"`
	NftProductSizeId   int64     `gorm:"column:nft_product_size_id" db:"nft_product_size_id" json:"nft_product_size_id" form:"nft_product_size_id"`
	ProductTitle       string    `gorm:"column:product_title" db:"product_title" json:"product_title" form:"product_title"`
	ProductPicture     string    `gorm:"column:product_picture" db:"product_picture" json:"product_picture" form:"product_picture"`
	Status             string    `gorm:"column:status" db:"status" json:"status" form:"status"` //  状态
	CreateTime         time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
	IsDelete           int64     `gorm:"column:is_delete" db:"is_delete" json:"is_delete" form:"is_delete"`
	Version            int64     `gorm:"column:version" db:"version" json:"version" form:"version"`
	IsWallet           int64     `gorm:"column:is_wallet" db:"is_wallet" json:"is_wallet" form:"is_wallet"`
	OrderSeId          int64     `gorm:"column:order_se_id" db:"order_se_id" json:"order_se_id" form:"order_se_id"`
	IsLock             int64     `gorm:"column:is_lock" db:"is_lock" json:"is_lock" form:"is_lock"`
	ActiveStartTime    int64     `gorm:"column:active_start_time" db:"active_start_time" json:"active_start_time" form:"active_start_time"` //  活动开始时间
	PropUserUuid       string    `gorm:"column:prop_user_uuid" db:"prop_user_uuid" json:"prop_user_uuid" form:"prop_user_uuid"`
	UserPicUrl         string    `gorm:"column:user_pic_url" db:"user_pic_url" json:"user_pic_url" form:"user_pic_url"`
	LoginName          string    `gorm:"column:login_name" db:"login_name" json:"login_name" form:"login_name"`
	BuyerUserPicUrl    string    `gorm:"column:buyer_user_pic_url" db:"buyer_user_pic_url" json:"buyer_user_pic_url" form:"buyer_user_pic_url"`
	BuyerUserLoginName string    `gorm:"column:buyer_user_login_name" db:"buyer_user_login_name" json:"buyer_user_login_name" form:"buyer_user_login_name"`
	ChangeType         int64     `gorm:"column:change_type" db:"change_type" json:"change_type" form:"change_type"`
	LandId             int64     `gorm:"column:land_id" db:"land_id" json:"land_id" form:"land_id"`                             //  土地id
	MarketType         string    `gorm:"column:market_type" db:"market_type" json:"market_type" form:"market_type"`             //  土地id
	SellContent        int64     `gorm:"column:sell_content" db:"sell_content" json:"sell_content" form:"sell_content"`         //  挂单内容：1藏品 / 2藏品+道具
	Level              int64     `gorm:"column:level" db:"level" json:"level" form:"level"`                                     //  打包资格等级 1/2/3....
	AgentOrderId       int64     `gorm:"column:agent_order_id" db:"agent_order_id" json:"agent_order_id" form:"agent_order_id"` //  代理订单id
	SendPoint          int64     `gorm:"column:send_point" db:"send_point" json:"send_point" form:"send_point"`                 //  是否发放消费积分
	IsMatch            int64     `gorm:"column:is_match" db:"is_match" json:"is_match" form:"is_match"`                         //  0=未被撮合处理 1=已被撮合处理
}

func (AiMatchProductNftSecondModel) TableName() string {
	return "ai_match_product_nft_second"
}

type AiMatchProductNftSecond struct {
	Ctx *gin.Context
}

func (n AiMatchProductNftSecond) GetProductNftSecondByOrderIds(orderIds []int64, req form.NftRecyclingByCountReq) (list []AiMatchProductNftSecondModel, err error) {
	err = cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductNftSecondModel{}).
		Where("order_id", orderIds).
		Where("status", "on_shelf").
		Limit(req.Count).Order("id asc").Scan(&list).Error
	return
}

func (n AiMatchProductNftSecond) UpdateNftProductStatusOffByOrderIds(orderIds []int64) (err error) {
	err = cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductNftSecondModel{}).
		Where("order_id IN (?)", orderIds).
		Where("is_delete", 0).
		Where("status", "on_shelf").
		Update("status", "off_shelf").Error
	return
}
