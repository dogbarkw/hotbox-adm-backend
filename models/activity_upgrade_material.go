package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/datatypes"
)

var ActivityUpgradeMaterialDal = &ActivityUpgradeMaterial{}

// ActivityUpgradeMaterial 升级活动材料表
type ActivityUpgradeMaterial struct {
	ID            int64          `json:"id" gorm:"id"`                           // 主键ID
	ActivityId    int64          `json:"activity_id" gorm:"activity_id"`         // 活动ID
	Number        int64          `json:"number" gorm:"number"`                   // 材料所需数量
	Image         string         `json:"image" gorm:"image"`                     // 材料封面图
	Weight        int64          `json:"weight" gorm:"weight"`                   // 材料权重
	Detail        datatypes.JSON `json:"detail" gorm:"detail"`                   // 材料详情
	LimitType     int8           `json:"limit_type" gorm:"limit_type"`           // 数量要求类型 0=持有组内藏品总数量满足要求 1=持有任意藏品满足数量要求
	DestroyMinNum int64          `json:"destroy_min_num" gorm:"destroy_min_num"` // 销毁数量最小值
	DestroyMaxNum int64          `json:"destroy_max_num" gorm:"destroy_max_num"` // 销毁数量最大值
	CreatedAt     time.Time      `json:"created_at" gorm:"created_at"`           // 创建时间
	UpdatedAt     time.Time      `json:"updated_at" gorm:"updated_at"`           // 更新时间
	AutoTtbStatus int8           `json:"auto_ttb_status" gorm:"auto_ttb_status"` // 自动进入时光宝库状态 1=待进 2=已进
}

// TableName 表名称
func (*ActivityUpgradeMaterial) TableName() string {
	return "activity_upgrade_material"
}

func (a *ActivityUpgradeMaterial) GetByParams(ctx context.Context, where map[string]any) (list []ActivityUpgradeMaterial, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (a *ActivityUpgradeMaterial) GetByParamsLimit(ctx context.Context, where map[string]any, limit int) (list []ActivityUpgradeMaterial, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Limit(limit).Scan(&list).Error
	return
}

func (a *ActivityUpgradeMaterial) One(ctx context.Context, id int) (resp ActivityUpgradeMaterial, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).First(&resp).Error
	return
}
