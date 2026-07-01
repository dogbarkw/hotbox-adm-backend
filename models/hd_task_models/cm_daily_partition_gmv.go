package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var HdDailyNftCategoryGmvDal = &HdDailyNftCategoryGmv{}

// 每日分区gmv统计表
type HdDailyNftCategoryGmv struct {
	Id           uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	Ymd          string    `gorm:"column:ymd;type:varchar(128);comment:统计日期;NOT NULL" json:"ymd"`
	CategoryPath string    `gorm:"column:category_path;type:varchar(128);comment:分区路径;NOT NULL" json:"category_path"`
	Category     string    `gorm:"column:category;type:varchar(128);comment:分区;NOT NULL" json:"category"`
	Rk           uint      `gorm:"column:rk;type:int(11) unsigned;default:0;comment:排名;NOT NULL" json:"rk"`
	Gmv          float64   `gorm:"column:gmv;type:decimal(10,2);default:0.00;comment:所有gmv;NOT NULL" json:"gmv"`
	UserCnt      uint      `gorm:"column:user_cnt;type:int(11) unsigned;default:0;comment:用户数;NOT NULL" json:"user_cnt"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdDailyNftCategoryGmv) TableName() string {
	return "hotbox_hd_daily_nft_category_gmv"
}

func (m *HdDailyNftCategoryGmv) FirstOrCreate(ctx context.Context, payload HdDailyNftCategoryGmv) (int64, error) {
	tx := cli.HotDogTaskGormDB.WithContext(ctx).Where("ymd = ? and category_path = ?", payload.Ymd, payload.CategoryPath).FirstOrCreate(&payload)
	return tx.RowsAffected, tx.Error
}

func (m *HdDailyNftCategoryGmv) UpdateByParams(ctx context.Context, where map[string]any, params map[string]any) error {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query.Where(k, v)
	}
	return query.Updates(params).Error
}

func (m *HdDailyNftCategoryGmv) GetGMVByParams(ctx context.Context, where map[string]any) (gmv float64, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Select("ifnull(sum(gmv), 0) as gmv").Scan(&gmv).Error
	return
}

func (m *HdDailyNftCategoryGmv) GetDailyNftCategoryGmvList(ctx context.Context, where map[string]any, order []string, pageNum, pageSize int) (list []HdDailyNftCategoryGmv, total int64, err error) {
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
