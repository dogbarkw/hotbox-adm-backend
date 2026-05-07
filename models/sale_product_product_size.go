package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

type SaleCalendarProductSizeModel struct {
	ID            int64     `gorm:"column:id" json:"id"`
	ProductId     int64     `gorm:"column:product_id" json:"product_id"` //  产品id
	Size          string    `gorm:"column:size" json:"size"`             //  产品尺寸
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"update_time"`
	IsDelete      int64     `gorm:"column:is_delete" json:"is_delete"` //  是否已删除,0:未删除,1:已删除
	SpuId         string    `gorm:"column:spu_id" json:"spu_id"`
	Price         float64   `gorm:"column:price" json:"price"`             //  价格
	StockCount    int64     `gorm:"column:stock_count" json:"stock_count"` //  库存
	SkuId         string    `gorm:"column:sku_id" json:"sku_id"`
	DiffPrice     float64   `gorm:"column:diff_price" json:"diff_price"`         //  差价
	Weight        int64     `gorm:"column:weight" json:"weight"`                 //  权重
	TotalCount    int64     `gorm:"column:total_count" json:"total_count"`       //  总量库存
	CostPrice     float64   `gorm:"column:cost_price" json:"cost_price"`         //  供货价
	DiscountPrice float64   `gorm:"column:discount_price" json:"discount_price"` //  折扣价
	OriginalPrice float64   `gorm:"column:original_price" json:"original_price"` //  划线价
	RedeemPrice   float64   `gorm:"column:redeem_price" json:"redeem_price"`     //  兑换金额
	RedeemPoint   int64     `gorm:"column:redeem_point" json:"redeem_point"`     //  兑换积分
}

func (SaleCalendarProductSizeModel) TableName() string {
	return "sale_calendar_product_size"
}

type SaleCalendarProductSize struct{}

func NewSaleCalendarProductSize() *SaleCalendarProductSize {
	return &SaleCalendarProductSize{}
}

func (s *SaleCalendarProductSize) GetSaleCalendarProductSizeByParams(ctx context.Context, where map[string]any, order []string, limit *int) (list []SaleCalendarProductSizeModel, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&SaleCalendarProductSizeModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}
