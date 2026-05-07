package dg_models

import (
	"time"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActivityModel struct {
	ID                   int64     `gorm:"column:id" json:"id"`                               //  id
	ActivityTitle        string    `gorm:"column:activity_title" json:"activity_title"`       //  活动title
	ActivityId           string    `gorm:"column:activity_id" json:"activity_id"`             //  活动id
	ActivityType         int       `gorm:"column:activity_type" json:"activity_type"`         //  活动类型 1 合成，2升级, 3 置换, 4 分解
	ActivityStartTs      int64     `gorm:"column:activity_start_ts" json:"activity_start_ts"` //  活动开始时间戳
	ActivityEndTs        int64     `gorm:"column:activity_end_ts" json:"activity_end_ts"`     //  活动结束时间戳
	Status               int64     `gorm:"column:status" json:"status"`                       //  0 下架状态, 1 已经上架
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsDelete             int64     `gorm:"column:is_delete" json:"is_delete"`
	Raid                 int64     `gorm:"column:raid" json:"raid"`
	DeActivityId         int64     `gorm:"de_activity_id" json:"-"`
	ActualActivityEndTs  int64     `gorm:"actual_activity_end_ts" json:"actual_activity_end_ts"`
	OnShelfTs            int64     `gorm:"on_shelf_ts" json:"-"`
	ShortName            string    `gorm:"short_name" json:"_"`
	ActivityNumCountFlag int       `gorm:"column:activity_num_count_flag;type:tinyint(4);default:0;comment:状态 0:未检查, 1:已检查;NOT NULL" json:"activity_num_count_flag"`
	Tm                   int64     `gorm:"tm; comment:定时时间戳" json:"tm"`
}

func (ActivityModel) TableName() string {
	return "dog_activity_list"
}

type Activity struct {
	Ctx *gin.Context
}

func (a Activity) GetByActivityIdAndType(activityId int64, activityType int64) (result ActivityModel, err error) {
	err = cli.GormDB.WithContext(a.Ctx).
		Model(&ActivityModel{}).
		Where("de_activity_id", activityId).
		Where("activity_type", activityType).
		First(&result).Error
	return result, err
}

func (a Activity) GetActivityByParams(where map[string][]any, order []string) (list []ActivityModel, err error) {
	query := cli.GormDB.WithContext(a.Ctx).
		Model(&ActivityModel{})

	for k, v := range where {
		query = query.Where(k, v...)
	}
	for _, v := range order {
		query = query.Order(v)
	}
	err = query.Scan(&list).Error
	if err != nil {
		logrus.Errorf("GetActivityList error: %v", err)
		return nil, err
	}
	return
}

func (a Activity) GetActivityByParamsV2(where map[string][]any, order []string, limit, offset int) (list []ActivityModel, err error) {
	query := cli.GormDB.WithContext(a.Ctx).
		Model(&ActivityModel{})

	for k, v := range where {
		query = query.Where(k, v...)
	}
	for _, v := range order {
		query = query.Order(v)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	err = query.Scan(&list).Error
	if err != nil {
		logrus.Errorf("GetActivityList error: %v", err)
		return nil, err
	}
	return
}
