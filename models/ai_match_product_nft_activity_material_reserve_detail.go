package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/gorm"
)

var ActivityMaterialReserveDetailDal = &AiMatchProductNftActivityMaterialReserveDetail{}

type AiMatchProductNftActivityMaterialReserveDetail struct {
	Id            int            `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	ReserveId     int64          `gorm:"column:reserve_id;type:int(11);comment:材料预留主表id;NOT NULL" json:"reserve_id"`
	ActivityId    int64          `gorm:"column:activity_id;type:bigint(20);comment:活动id;NOT NULL" json:"activity_id"`
	ActivityType  int            `gorm:"column:activity_type;type:int(11);comment:活动类型 1=合成;NOT NULL" json:"activity_type"`
	MaterialId    int64          `gorm:"column:material_id" json:"material_id"`
	MaterialUuid  string         `gorm:"column:material_uuid;type:varchar(32);comment:材料组id;NOT NULL" json:"material_uuid"`
	MaterialType  string         `gorm:"column:material_type;type:varchar(32);comment:材料类型 nft=藏品 prop=道具 product=藏品;NOT NULL" json:"material_type"`
	MaterialName  string         `gorm:"column:material_name;type:varchar(64);comment:材料名称;NOT NULL" json:"material_name"`
	ProductId     uint64         `gorm:"column:product_id;type:bigint(20);comment:藏品id" json:"product_id"`
	ProductSizeId uint64         `gorm:"column:product_size_id;type:bigint(20);comment:藏品尺码id" json:"product_size_id"`
	PropId        uint64         `gorm:"column:prop_id;type:bigint(20);comment:道具id" json:"prop_id"`
	ReserveNum    int64          `gorm:"column:reserve_num;type:int(11);default:0;comment:预留份数;NOT NULL" json:"reserve_num"`
	Status        int            `gorm:"column:status;type:tinyint(3);default:0;comment: 0=待执行 1=执行成功 -1=执行失败" json:"status"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (m *AiMatchProductNftActivityMaterialReserveDetail) TableName() string {
	return "ai_match_product_nft_activity_material_reserve_detail"
}

func (*AiMatchProductNftActivityMaterialReserveDetail) GetByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftActivityMaterialReserveDetail, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftActivityMaterialReserveDetail{})
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

// 合成
type AiMatchProductNftActivityMaterialReserveDetailJoinCombination struct {
	AiMatchProductNftActivityMaterialReserveDetail
	OnSaleStatus int64 `gorm:"column:on_sale_status" db:"on_sale_status" json:"on_sale_status" form:"on_sale_status"`
}

func (*AiMatchProductNftActivityMaterialReserveDetail) GetReserveDetailJoinCombinationByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftActivityMaterialReserveDetailJoinCombination, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftActivityMaterialReserveDetail{}).
		Select("ai_match_product_nft_activity_material_reserve_detail.*, ai_match_product_nft_combination.on_sale_status").
		Joins("JOIN ai_match_product_nft_combination ON ai_match_product_nft_activity_material_reserve_detail.activity_id = ai_match_product_nft_combination.id")
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

// 置换、分解
type AiMatchProductNftActivityMaterialReserveDetailJoinReplace struct {
	AiMatchProductNftActivityMaterialReserveDetail
	OnSaleStatus int64 `gorm:"column:on_sale_status" db:"on_sale_status" json:"on_sale_status" form:"on_sale_status"`
}

func (*AiMatchProductNftActivityMaterialReserveDetail) GetReserveDetailJoinReplaceByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftActivityMaterialReserveDetailJoinReplace, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftActivityMaterialReserveDetail{}).
		Select("ai_match_product_nft_activity_material_reserve_detail.*, ai_match_product_nft_replace.on_sale_status").
		Joins("JOIN ai_match_product_nft_replace ON ai_match_product_nft_activity_material_reserve_detail.activity_id = ai_match_product_nft_replace.replace_id")
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (*AiMatchProductNftActivityMaterialReserveDetail) Create(ctx context.Context, input []AiMatchProductNftActivityMaterialReserveDetail) error {
	return cli.HotDogGormDB.WithContext(ctx).Create(&input).Error
}

func (*AiMatchProductNftActivityMaterialReserveDetail) UpdateByParams(where, payload map[string]any) (affectRow int64, err error) {
	query := cli.HotDogGormDB.Model(&AiMatchProductNftActivityMaterialReserveDetail{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	query = query.Updates(payload)
	affectRow = query.RowsAffected
	err = query.Error
	return
}

type GetJoinMainWithParamsItem struct {
	AiMatchProductNftActivityMaterialReserveDetail
	ExecTime     int64 `json:"exec_time"`     // 预计执行时间
	ExecEndTime  int64 `json:"exec_end_time"` // 实际结束时间（用于判断是否缓慢执行）
	ActivityType int   `json:"activity_type"` // 活动类型
}

func (m *AiMatchProductNftActivityMaterialReserveDetail) GetJoinMainWithParams(where map[string]any) (result []GetJoinMainWithParamsItem, err error) {
	query := cli.HotDogGormDB.Table(m.TableName())
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Select("ai_match_product_nft_activity_material_reserve_detail.*, r.exec_time, r.exec_end_time, r.activity_type").
		Joins(`
		INNER JOIN
		ai_match_product_nft_activity_material_reserve r 
		ON 
		r.id = ai_match_product_nft_activity_material_reserve_detail.reserve_id`).
		Where("ai_match_product_nft_activity_material_reserve_detail.deleted_at is NULL").
		Scan(&result).
		Error
	return
}
