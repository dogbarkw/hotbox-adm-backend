package models

import (
	"time"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
)

type SaleProductNftSizePfpModel struct {
	PfpId            int64     `gorm:"column:pfp_id" json:"pfp_id"`
	ProductId        int64     `gorm:"column:product_id" json:"product_id"`
	NftProductSizeId int64     `gorm:"column:nft_product_size_id" json:"nft_product_size_id"`
	ReceiverCity     string    `gorm:"column:receiver_city" json:"receiver_city"`
	FrontCover       string    `gorm:"column:front_cover" json:"front_cover"`
	OriginMedia      string    `gorm:"column:origin_media" json:"origin_media"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
	IsDelete         int64     `gorm:"column:is_delete" json:"is_delete"`
}

func (SaleProductNftSizePfpModel) TableName() string {
	return "sale_product_nft_size_pfp"
}

type SaleProductNftSizePfp struct {
	Ctx *gin.Context
}

func (s SaleProductNftSizePfp) GetOneByParams(where map[string]any) (r SaleProductNftSizePfpModel, err error) {
	err = cli.HotDogGormDB.WithContext(s.Ctx).Where(where).First(&r).Error
	return
}
