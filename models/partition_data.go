package models

import (
	"context"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
)

var PartitionDataDal = &PartitionData{}

// 分区表(数据同学使用)
type PartitionData struct {
	Id         uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	MainId     int64     `gorm:"column:main_id;type:int(11);default:0;comment:主分区id;NOT NULL" json:"main_id"`
	ChildId    int64     `gorm:"column:child_id;type:int(11);default:0;comment:子分区id;NOT NULL" json:"child_id"`
	MainTitle  string    `gorm:"column:main_title;type:varchar(100);comment:主分区标题;NOT NULL" json:"main_title"`
	ChildTitle string    `gorm:"column:child_title;type:varchar(100);comment:子分区标题;NOT NULL" json:"child_title"`
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(4);default:0;NOT NULL" json:"is_delete"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
}

func (PartitionData) TableName() string {
	return "partition_data"
}

func (n PartitionData) GetPartitionData(ctx context.Context, partitionIds [][]any) (result map[string]PartitionData, err error) {
	result = make(map[string]PartitionData, len(partitionIds))
	var list []PartitionData
	err = cli.HotDogGormDB.WithContext(ctx).Model(&PartitionData{}).
		Where("(main_id,child_id) in ?", partitionIds).
		Where("is_delete", 0).
		Scan(&list).Error
	if err != nil {
		return result, err
	}
	for _, v := range list {
		result[fmt.Sprintf("%d-%d", v.MainId, v.ChildId)] = v
	}
	return
}

func (n PartitionData) GetPartitionList(ctx context.Context, mainIds []int64) (result map[int64][]PartitionData, list []PartitionData, err error) {
	result = make(map[int64][]PartitionData, len(mainIds))
	err = cli.HotDogGormDB.WithContext(ctx).Model(&PartitionData{}).
		Where("main_id in ?", mainIds).
		Where("is_delete = ?", 0).
		Where("child_id > ?", 0).
		Scan(&list).Error
	if err != nil {
		return result, list, err
	}
	for _, v := range list {
		result[v.MainId] = append(result[v.MainId], v)
	}
	return
}

func (n PartitionData) GetByParams(ctx context.Context, where map[string]any, order []string) (list []PartitionData, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&PartitionData{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	for _, s := range order {
		query.Order(s)
	}

	err = query.Scan(&list).Error
	return
}
