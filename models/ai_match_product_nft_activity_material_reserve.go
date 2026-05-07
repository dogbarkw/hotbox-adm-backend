package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/gorm"
)

var ActivityMaterialReserveDal = &AiMatchProductNftActivityMaterialReserve{}

type AiMatchProductNftActivityMaterialReserve struct {
	Id                 int64          `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	ActivityId         int64          `gorm:"column:activity_id;type:bigint(20);comment:活动id;NOT NULL" json:"activity_id"`
	ActivityType       int            `gorm:"column:activity_type;type:int(11);comment:活动类型;NOT NULL" json:"activity_type"`
	ExecStatus         int            `gorm:"column:exec_status;type:int(11);default:0;comment:任务状态 0=待执行 1=执行成功 -1=执行失败;NOT NULL" json:"exec_status"`
	ExecTime           int64          `gorm:"column:exec_time;type:bigint(20);default:0;comment:预计执行时间" json:"exec_time"`
	ExecEndTime        int64          `gorm:"column:exec_end_time;type:bigint(20);default:0;comment:实际结束时间" json:"exec_end_time"`
	ActivityReserveNum int64          `gorm:"column:activity_reserve_num;type:int(11);default:0;comment:活动预留份数;NOT NULL" json:"activity_reserve_num"`
	Remark             string         `gorm:"column:remark;type:varchar(64);comment:备注;NOT NULL" json:"remark"`
	UserId             int64          `gorm:"column:user_id;comment: 后台用户" json:"user_id"`
	UserName           string         `gorm:"column:user_name;comment:后台用户名称" json:"user_name"`
	CreatedAt          time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (m *AiMatchProductNftActivityMaterialReserve) TableName() string {
	return "ai_match_product_nft_activity_material_reserve"
}

func (m *AiMatchProductNftActivityMaterialReserve) Create(ctx context.Context, input *AiMatchProductNftActivityMaterialReserve) error {
	return cli.HotDogGormDB.WithContext(ctx).Create(&input).Error
}

func (m *AiMatchProductNftActivityMaterialReserve) One(ctx context.Context, id int64) (resp AiMatchProductNftActivityMaterialReserve, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Where("id", id).First(&resp).Error
	return
}

func (m *AiMatchProductNftActivityMaterialReserve) GetByParams(ctx context.Context, where map[string]any) (list []AiMatchProductNftActivityMaterialReserve, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(m)
	for k, v := range where {
		query = query.Where(k, v)
	}
	err = query.Scan(&list).Error
	return
}

func (m *AiMatchProductNftActivityMaterialReserve) UpdateByParams(where, payload map[string]any) (affectRow int64, err error) {
	query := cli.HotDogGormDB.Model(&AiMatchProductNftActivityMaterialReserve{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	query = query.Updates(payload)
	affectRow = query.RowsAffected
	err = query.Error
	return
}
