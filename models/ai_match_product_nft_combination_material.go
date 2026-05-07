package models

import (
	"context"

	"hotbox-adm-backend/cli"

	"gorm.io/datatypes"
)

var CombinationMaterialDal = &AiMatchProductNftCombinationMaterial{}

type AiMatchProductNftCombinationMaterial struct {
	MaterialUuid  string         `gorm:"column:material_uuid" json:"material_uuid"`
	CombinationId uint64         `gorm:"column:combination_id" json:"combination_id"`
	MaterialNum   uint32         `gorm:"column:material_num" json:"material_num"`
	MaterialPic   string         `gorm:"column:material_pic" json:"material_pic"`
	Weight        uint32         `gorm:"column:weight" json:"weight"`
	MaterialInfo  datatypes.JSON `gorm:"column:material_info" json:"material_info"`
	MaterialType  uint8          `gorm:"column:material_type" json:"material_type"`
	CreateTime    string         `gorm:"column:create_time" json:"create_time"`
	UpdateTime    string         `gorm:"column:update_time" json:"update_time"`
	IsDelete      uint8          `gorm:"column:is_delete" json:"is_delete"`
}

func (*AiMatchProductNftCombinationMaterial) TableName() string {
	return "ai_match_product_nft_combination_material"
}

func (a *AiMatchProductNftCombinationMaterial) GetByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftCombinationMaterial, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (a *AiMatchProductNftCombinationMaterial) GetByParamsLimit(ctx context.Context, where map[string]any, limit int) (list []AiMatchProductNftCombinationMaterial, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Limit(limit).Scan(&list).Error
	return
}

func (a *AiMatchProductNftCombinationMaterial) One(ctx context.Context, materialUuid string) (resp AiMatchProductNftCombinationMaterial, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("material_uuid", materialUuid).First(&resp).Error
	return
}
