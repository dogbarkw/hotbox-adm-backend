package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var CmDailyGmvDal = &CmDailyGmv{}

// 每日gmv统计表
type CmDailyGmv struct {
	Id        uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	Ymd       string    `gorm:"column:ymd;type:varchar(128);comment:统计日期;NOT NULL" json:"ymd"`
	Gmv       float64   `gorm:"column:gmv;type:decimal(10,2);default:0;comment:所有gmv;NOT NULL" json:"gmv"`
	RGmv      float64   `gorm:"column:r_gmv;type:decimal(10,2);default:0;comment:排除测试用户的gmv;NOT NULL" json:"r_gmv"`
	UserCnt   uint      `gorm:"column:user_cnt;type:int(11) unsigned;default:0;comment:用户数;NOT NULL" json:"user_cnt"`
	RUserCnt  uint      `gorm:"column:r_user_cnt;type:int(11) unsigned;default:0;comment:排除测试用户的用户数;NOT NULL" json:"r_user_cnt"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *CmDailyGmv) TableName() string {
	return "hotbox_daily_gmv"
}

func (m *CmDailyGmv) FirstOrCreate(ctx context.Context, payload CmDailyGmv) (int64, error) {
	tx := cli.HotDogTaskGormDB.WithContext(ctx).Where("ymd = ?", payload.Ymd).FirstOrCreate(&payload)
	return tx.RowsAffected, tx.Error
}

func (m *CmDailyGmv) UpdateByParams(ctx context.Context, where map[string]any, params map[string]any) error {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query.Where(k, v)
	}
	return query.Updates(params).Error
}

func (m *CmDailyGmv) GetCmDailyGmvList(ctx context.Context, where map[string]any, order []string, pageNum, pageSize int) (list []CmDailyGmv, total int64, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}

	offset := (pageNum - 1) * pageSize
	if pageSize > 0 {
		query.Limit(pageSize)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}
