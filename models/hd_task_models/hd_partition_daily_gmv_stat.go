package hd_task_models

import (
	"context"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/gorm/clause"
)

var HdPartitionDailyGmvStatDal = &HdPartitionDailyGmvStat{}

// 分区每日gmv统计表
type HdPartitionDailyGmvStat struct {
	Id          int64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	MainId      int64     `gorm:"column:main_id;type:bigint(20);default:0;comment:主分区id;NOT NULL" json:"main_id"`
	ChildId     int64     `gorm:"column:child_id;type:bigint(20);default:0;comment:子分区id;NOT NULL" json:"child_id"`
	TargetGmv   float64   `gorm:"column:target_gmv;type:decimal(10,2);default:0.00;comment:日均gmv目标;NOT NULL" json:"target_gmv"`
	CurrentGmv  float64   `gorm:"column:current_gmv;type:decimal(10,2);default:0.00;comment:当前gmv;NOT NULL" json:"current_gmv"`
	ShareIncome float64   `gorm:"column:share_income;type:decimal(10,2);default:0.00;comment:当前分账;NOT NULL" json:"share_income"`
	PreGmv      float64   `gorm:"column:pre_gmv;type:decimal(10,2);default:0.00;comment:前一天的gmv;NOT NULL" json:"pre_gmv"`
	Date        int       `gorm:"column:date;type:int(11);default:0;comment:统计日期;NOT NULL" json:"date"`
	Status      int       `gorm:"column:status;type:int(11);default:0;comment:账号补偿状态 1恢复(分账中) 2暂停(正常比例);NOT NULL" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdPartitionDailyGmvStat) TableName() string {
	return "hotbox_partition_daily_gmv_stat"
}

func (m *HdPartitionDailyGmvStat) Create(ctx context.Context, data *HdPartitionDailyGmvStat) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Create(&data).Error
}

func (m *HdPartitionDailyGmvStat) GetPartitionTodayGmvStatList(ctx context.Context, partitionIds [][]any) (dataMap map[string]HdPartitionDailyGmvStat, err error) {
	var data []HdPartitionDailyGmvStat
	dataMap = make(map[string]HdPartitionDailyGmvStat, len(partitionIds))
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).Where("(main_id, child_id) IN ?", partitionIds).
		Where("date = ?", time.Now().Format("20060102")).
		Where("target_gmv > ?", 0).
		Scan(&data).Error
	if err != nil {
		return dataMap, err
	}
	for _, datum := range data {
		dataMap[fmt.Sprintf("%d-%d", datum.MainId, datum.ChildId)] = datum
	}
	return
}

func (m *HdPartitionDailyGmvStat) GetByParams(ctx context.Context, params map[string]any, mainId []int64, childId []int64) (data []HdPartitionDailyGmvStat, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range params {
		query = query.Where(k, v)
	}
	if mainId != nil && len(mainId) > 0 {
		query = query.Where("main_id IN ?", mainId)
	}
	if childId != nil && len(childId) > 0 {
		query = query.Where("child_id IN ?", childId)
	}
	err = query.Scan(&data).Error
	return
}

func (m *HdPartitionDailyGmvStat) GetOneByParams(ctx context.Context, params map[string]any) (data HdPartitionDailyGmvStat, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range params {
		query = query.Where(k, v)
	}
	err = query.First(&data).Error
	return
}

func (m *HdPartitionDailyGmvStat) UpdateByParams(ctx context.Context, where map[string]any, params map[string]any) error {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query = query.Where(k, v)
	}
	return query.Updates(params).Error
}

func (m *HdPartitionDailyGmvStat) DeleteByPartitionId(ctx context.Context, mainId int64, childId int64) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).
		Where("main_id", mainId).
		Where("child_id", childId).
		Where("date >= ?", time.Now().Format("20060102")).Delete(&m).Error
}

type WeeklyGmvStatData struct {
	MainId            int64   `json:"main_id"`
	ChildId           int64   `json:"child_id"`
	WeeklyGmv         float64 `json:"weekly_gmv"`
	WeeklyShareIncome float64 `json:"weekly_share_income"`
}

func (m *HdPartitionDailyGmvStat) SumWeeklyGmvStat(ctx context.Context, startDate int, endDate int) (dataMap map[string]WeeklyGmvStatData, err error) {
	dataMap = make(map[string]WeeklyGmvStatData)
	var list []WeeklyGmvStatData
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	err = query.Where("date >= ?", startDate).
		Where("date <= ?", endDate).
		Select("main_id, child_id, sum(current_gmv  + pre_gmv) as weekly_gmv, sum(share_income) as weekly_share_income").
		Group("main_id, child_id").Scan(&list).Error
	if err != nil {
		return
	}
	for _, stat := range list {
		dataMap[fmt.Sprintf("%d-%d", stat.MainId, stat.ChildId)] = WeeklyGmvStatData{
			MainId:            stat.MainId,
			ChildId:           stat.ChildId,
			WeeklyGmv:         stat.WeeklyGmv,
			WeeklyShareIncome: stat.WeeklyShareIncome,
		}
	}
	return
}

func (m *HdPartitionDailyGmvStat) UpSertNextGmv(ctx context.Context, data *HdPartitionDailyGmvStat) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "main_id"}, {Name: "child_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"target_gmv", "pre_gmv", "status"}),
		}).
		Create(&data).Error
}
