package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var AiMatchProductNftReplaceContentDal = &AiMatchProductNftReplaceContent{}

// AiMatchProductNftReplaceContent undefined
type AiMatchProductNftReplaceContent struct {
	ReplaceContentId            int64     `json:"replace_content_id" gorm:"replace_content_id"`
	ReplaceId                   int64     `json:"replace_id" gorm:"replace_id"`
	ReplaceType                 string    `json:"replace_type" gorm:"replace_type"`
	ReplaceSerial               int64     `json:"replace_serial" gorm:"replace_serial"`
	TargetType                  string    `json:"target_type" gorm:"target_type"`
	TargetCount                 int64     `json:"target_count" gorm:"target_count"`
	TargetName                  string    `json:"target_name" gorm:"target_name"`
	ProductId                   int64     `json:"product_id" gorm:"product_id"`
	NftProductSizeId            int64     `json:"nft_product_size_id" gorm:"nft_product_size_id"`
	CouponId                    int64     `json:"coupon_id" gorm:"coupon_id"`
	PropId                      int64     `json:"prop_id" gorm:"prop_id"`
	CreateTime                  time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime                  time.Time `json:"update_time" gorm:"update_time"`
	IsDelete                    int64     `json:"is_delete" gorm:"is_delete"`
	GroupId                     int64     `json:"group_id" gorm:"group_id"`
	InscriptionProductId        int64     `json:"inscription_product_id" gorm:"inscription_product_id"`                   // 铭文藏品 ID
	InscriptionProductNftSizeId int64     `json:"inscription_product_nft_size_id" gorm:"inscription_product_nft_size_id"` // 铭文源文件 ID
	AutoTtb                     int8      `json:"auto_ttb" gorm:"auto_ttb"`                                               // 活动结束是否自动进入时光宝库 1=是
	AutoTtbStatus               int8      `json:"auto_ttb_status" gorm:"auto_ttb_status"`                                 // 自动进入时光宝库状态 1=待进 2=已进
}

// TableName 表名称
func (*AiMatchProductNftReplaceContent) TableName() string {
	return "ai_match_product_nft_replace_content"
}

func (a *AiMatchProductNftReplaceContent) GetByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftReplaceContent, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (a *AiMatchProductNftReplaceContent) GetByParamsLimit(ctx context.Context, where map[string]any, limit int) (list []AiMatchProductNftReplaceContent, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Limit(limit).Scan(&list).Error
	return
}

func (a *AiMatchProductNftReplaceContent) One(ctx context.Context, replaceContentId int) (resp AiMatchProductNftReplaceContent, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("replace_content_id", replaceContentId).First(&resp).Error
	return
}
