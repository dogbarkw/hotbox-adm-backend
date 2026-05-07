package models

import (
	"context"
	"time"

	"gorm.io/gorm/clause"

	"hotbox-adm-backend/cli"

	"gorm.io/gorm"
)

var TableScanRecordDal = &TableScanRecord{}

// 数据表扫描记录
type TableScanRecord struct {
	Id        int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Table     string    `gorm:"column:table;type:varchar(64);comment:表名;NOT NULL" json:"table"`
	LastId    int64     `gorm:"column:last_id;type:bigint(20);default:0;comment:最后一个处理的ID;NOT NULL" json:"last_id"`
	TaskType  int       `gorm:"column:task_type;type:int(11);default:0;comment:任务类型，1统计艺术家活动数;NOT NULL" json:"task_type"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (m *TableScanRecord) TableName() string {
	return "table_scan_record"
}

func (m *TableScanRecord) Save(tx *gorm.DB, data TableScanRecord) error {
	return tx.Model(m).Save(&data).Error
}

func (m *TableScanRecord) GetLastScanId(ctx context.Context, table string, taskType int) (LastId int64, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Select("last_id").
		Where("table", table).
		Where("task_type", taskType).
		Order("id DESC").
		Take(&LastId).Error
	return
}

func (m *TableScanRecord) UpdateLastScanId(tx *gorm.DB, dm *TableScanRecord) error {
	// return tx.Model(&m).Where("table", table).Where("task_type", taskType).Update("last_id", lastId).Error
	return tx.Model(&m).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "table"}, {Name: "task_type"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"last_id": dm.LastId,
		}),
	}).Create(&dm).Error
}
