package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var HdYopTestUserRateRecordDal = &HdYopTestUserRateRecord{}

// 特殊账号分成记录表
type HdYopTestUserRateRecord struct {
	Id            uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	YopTestUserId int64     `gorm:"column:yop_test_user_id;type:bigint(20);default:0;comment:特殊账号ID;NOT NULL" json:"yop_test_user_id"`
	Rate          int       `gorm:"column:rate;type:int(11);default:0;comment:分成比例;NOT NULL" json:"rate"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdYopTestUserRateRecord) TableName() string {
	return "hotbox_yop_test_user_rate_record"
}

func (m *HdYopTestUserRateRecord) Create(ctx context.Context, payload *HdYopTestUserRateRecord) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Create(&payload).Error
}

func (m *HdYopTestUserRateRecord) GetByParams(ctx context.Context, where map[string][]any, order []string, limit int) (list []HdYopTestUserRateRecord, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query = query.Where(k, v)
	}
	for _, v := range order {
		query.Order(v)
	}
	if limit > 0 {
		query.Limit(limit)
	}

	err = query.Scan(&list).Error
	return
}

func (m *HdYopTestUserRateRecord) Delete(ctx context.Context, userIds []int64) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Unscoped().Delete(&m, "yop_test_user_id", userIds).Error
}
