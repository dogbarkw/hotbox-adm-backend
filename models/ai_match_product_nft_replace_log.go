package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/datatypes"
)

var AiMatchProductNftReplaceLogDal = AiMatchProductNftReplaceLog{}

// AiMatchProductNftReplaceLog 置换活动日志表
type AiMatchProductNftReplaceLog struct {
	ID         int64                                                      `json:"id" gorm:"id"`                 // 主键ID
	ReplaceId  int64                                                      `json:"replace_id" gorm:"replace_id"` // 活动ID
	UserId     int64                                                      `json:"user_id" gorm:"user_id"`       // 用户ID
	Remark     string                                                     `json:"remark" gorm:"remark"`         // 备注
	ResultInfo datatypes.JSONSlice[AiMatchProductNftReplaceLogResultInfo] `json:"result_info" gorm:"result_info"`
	CreatedAt  time.Time                                                  `json:"created_at" gorm:"created_at"` // 创建时间
	UpdatedAt  time.Time                                                  `json:"updated_at" gorm:"updated_at"` // 更新时间
}

type AiMatchProductNftReplaceLogResultInfo struct {
	// Sn          string `json:"sn"`
	Num         int    `json:"num"`
	ResultType  string `json:"result_type"`
	ProductName string `json:"product_name"`
}

// TableName 表名称
func (*AiMatchProductNftReplaceLog) TableName() string {
	return "ai_match_product_nft_replace_log"
}

func (*AiMatchProductNftReplaceLog) GetByReplaceLogListByReplaceId(ctx context.Context, replaceId int) (result []AiMatchProductNftReplaceLog, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(AiMatchProductNftReplaceLog{}).Where("replace_id", replaceId).Find(&result).Error
	return
}
