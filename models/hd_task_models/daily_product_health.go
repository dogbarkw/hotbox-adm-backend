package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var DailyProductHealthDal = &DailyProductHealth{}

type DailyProductHealth struct {
	Id               int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	ProductId        int64     `gorm:"column:product_id;type:bigint(20);comment:商品id;NOT NULL" json:"product_id"`
	ProductTitle     string    `gorm:"column:product_title;type:varchar(500);comment:商品title" json:"product_title"`
	NftProductSizeId int64     `gorm:"column:nft_product_size_id;type:bigint(20);comment:nft商品款id;NOT NULL" json:"nft_product_size_id"`
	Ymd              string    `gorm:"column:ymd;type:varchar(12);default:0;comment:日期" json:"ymd"`
	Info             string    `gorm:"column:info;type:json;comment:信息" json:"info"`
	Type             int       `gorm:"column:type;type:tinyint(4);default:0;comment:类型:1 健康, 2 非健康" json:"type"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *DailyProductHealth) TableName() string {
	return "daily_product_health"
}

func (m *DailyProductHealth) GetProductHealthByYmd(ctx context.Context, ymd string) (result []DailyProductHealth, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(m).Where("ymd = ?", ymd).Where("type in (1,2)").Find(&result).Error
	return
}
