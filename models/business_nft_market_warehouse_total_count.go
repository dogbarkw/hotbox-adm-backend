package models

import (
	"context"

	"hotbox-adm-backend/cli"
)

type BusinessNftMarketWarehouseTotalCount struct {
	ctx           context.Context
	ProductId     int64  `json:"product_id" gorm:"product_id"`           // 产品id
	ProductSizeId int64  `json:"product_size_id" gorm:"product_size_id"` // nft款式id
	ProductTitle  string `json:"product_title" gorm:"product_title"`     // 产品标题
	NftCount      int64  `json:"nft_count" gorm:"nft_count"`             // 剩余份数
	RemainCount   int64  `json:"remain_count" gorm:"remain_count"`       // 流通份数
}

// TableName 表名称
func (*BusinessNftMarketWarehouseTotalCount) TableName() string {
	return "business_nft_market_warehouse_total_count"
}

func NewBusinessNftMarketWarehouseTotalCount(ctx context.Context) *BusinessNftMarketWarehouseTotalCount {
	return &BusinessNftMarketWarehouseTotalCount{
		ctx: ctx,
	}
}

func (b *BusinessNftMarketWarehouseTotalCount) GetByProductIdAndSizeId(productId, productSizeId int64) (data BusinessNftMarketWarehouseTotalCount, err error) {
	err = cli.HotDogGormDB.WithContext(b.ctx).Model(&data).Where("product_id = ? AND product_size_id = ?", productId, productSizeId).First(&data).Error
	return
}
