package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var AiMatchProductNftCombinationLogDal = AiMatchProductNftCombinationLog{}

// AiMatchProductNftCombinationLog 合成活动日志表
type AiMatchProductNftCombinationLog struct {
	ID        int64     `json:"id" gorm:"id"`                 // 主键ID
	CombineId int64     `json:"combine_id" gorm:"combine_id"` // 活动ID
	UserId    int64     `json:"user_id" gorm:"user_id"`       // 用户ID
	Remark    string    `json:"remark" gorm:"remark"`         // 备注
	CreatedAt time.Time `json:"created_at" gorm:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"` // 更新时间
}

// TableName 表名称
func (*AiMatchProductNftCombinationLog) TableName() string {
	return "ai_match_product_nft_combination_log"
}

func (*AiMatchProductNftCombinationLog) GetSuccessCountByCombineId(ctx context.Context, combineId int) (result int64, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(AiMatchProductNftCombinationLog{}).Where("combine_id", combineId).Count(&result).Error
	return
}
