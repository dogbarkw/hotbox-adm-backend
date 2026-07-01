package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var HdPartitionIncomeRecordDal = &HdPartitionIncomeRecord{}

// 分区进账记录表
type HdPartitionIncomeRecord struct {
	Id            uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	MainId        int64      `gorm:"column:main_id;type:bigint(20);default:0;comment:主分区id;NOT NULL" json:"main_id"`
	ChildId       int64      `gorm:"column:child_id;type:bigint(20);default:0;comment:子分区id;NOT NULL" json:"child_id"`
	YopTestUserId int64      `gorm:"column:yop_test_user_id;type:bigint(20);default:0;comment:特殊账号ID;NOT NULL" json:"yop_test_user_id"`
	Fee           float64    `gorm:"column:fee;type:decimal(10,2);default:0.00;comment:金额;NOT NULL" json:"fee"`
	Rate          int        `gorm:"column:rate;type:int(11);default:0;comment:当前分成比例;NOT NULL" json:"rate"`
	Income        float64    `gorm:"column:income;type:decimal(10,2);default:0.00;comment:进账;NOT NULL" json:"income"`
	IncomeTime    *time.Time `gorm:"column:income_time;type:datetime;comment:进账时间" json:"income_time"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdPartitionIncomeRecord) TableName() string {
	return "hotbox_partition_income_record"
}

func (m *HdPartitionIncomeRecord) BatchCreate(ctx context.Context, payload []*HdPartitionIncomeRecord) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).CreateInBatches(&payload, 30).Error
}

type PartitionIncomeStat struct {
	MainId  int64   `json:"main_id"`
	ChildId int64   `json:"child_id"`
	Income  float64 `json:"income"`
}

func (m *HdPartitionIncomeRecord) SumPartitionIncomeByTimeRange(ctx context.Context, partitionIds [][]any, start time.Time, end time.Time) (list []PartitionIncomeStat, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).Where("income_time >= ?", start).Where("income_time < ?", end)
	if len(partitionIds) > 0 {
		query = query.Where("(main_id,child_id) in ?", partitionIds)
	}
	err = query.Select("main_id,child_id,sum(income) as income").Group("main_id,child_id").Scan(&list).Error
	return
}
