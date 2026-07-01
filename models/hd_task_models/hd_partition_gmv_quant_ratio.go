package hd_task_models

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/clause"

	"hotbox-adm-backend/cli"
)

// 分区 GMV 量化配比表（各平台独立，表名前缀区分）
type HdPartitionGmvQuantRatio struct {
	Id         int64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	MainId     int64     `gorm:"column:main_id;type:bigint(20);default:0;comment:主分区id;NOT NULL" json:"main_id"`
	ChildId    int64     `gorm:"column:child_id;type:bigint(20);default:0;comment:子分区id;NOT NULL" json:"child_id"`
	QuantRatio float64   `gorm:"column:quant_ratio;type:decimal(10,2);default:0.00;comment:量化配比倍数;NOT NULL" json:"quant_ratio"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdPartitionGmvQuantRatio) TableName() string {
	return "hotbox_partition_gmv_quant_ratio"
}

var HdPartitionGmvQuantRatioDal = &HdPartitionGmvQuantRatio{}

func (m *HdPartitionGmvQuantRatio) GetByPartitionId(ctx context.Context, mainId int64, childId int64) (data HdPartitionGmvQuantRatio, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Where("main_id = ?", mainId).Where("child_id = ?", childId).First(&data).Error
	return
}

func (m *HdPartitionGmvQuantRatio) GetByPartitionIds(ctx context.Context, partitionIds [][]any) (dataMap map[string]HdPartitionGmvQuantRatio, err error) {
	if len(partitionIds) == 0 {
		return dataMap, nil
	}
	dataMap = make(map[string]HdPartitionGmvQuantRatio)
	var data []HdPartitionGmvQuantRatio
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).Where("(main_id, child_id) IN ?", partitionIds).Scan(&data).Error
	if err != nil {
		return dataMap, err
	}
	for _, datum := range data {
		dataMap[fmt.Sprintf("%d-%d", datum.MainId, datum.ChildId)] = datum
	}
	return
}

func (m *HdPartitionGmvQuantRatio) Upsert(ctx context.Context, mainId int64, childId int64, quantRatio float64) error {
	data := HdPartitionGmvQuantRatio{
		MainId:     mainId,
		ChildId:    childId,
		QuantRatio: quantRatio,
	}
	return cli.HotDogTaskGormDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "main_id"}, {Name: "child_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"quant_ratio", "updated_at"}),
		}).
		Create(&data).Error
}
