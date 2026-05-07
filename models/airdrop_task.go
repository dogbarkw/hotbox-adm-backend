package models

import (
	"time"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
)

type AirdropTaskModel struct {
	Id               uint64    `gorm:"column:id" json:"id"`
	OperatorId       uint64    `gorm:"column:operator_id" json:"operator_id"`
	OperatorName     string    `gorm:"column:operator_name" json:"operator_name"`
	Title            string    `gorm:"column:title" json:"title"`
	BenefitType      string    `gorm:"column:benefit_type" json:"benefit_type"`
	PropId           uint64    `gorm:"column:prop_id" json:"prop_id"`
	PropName         string    `gorm:"column:prop_name" json:"prop_name"`
	Source           uint32    `gorm:"column:source" json:"source"`
	UserList         string    `gorm:"column:user_list" json:"user_list"`
	Filename         string    `gorm:"column:filename" json:"filename"`
	FileUrl          string    `gorm:"column:file_url" json:"file_url"`
	SnapshotDataId   uint64    `gorm:"column:snapshot_data_id" json:"snapshot_data_id"`
	SnapshotDataName string    `gorm:"column:snapshot_data_name" json:"snapshot_data_name"`
	DropUserNum      uint64    `gorm:"column:drop_user_num" json:"drop_user_num"`
	SendNum          uint32    `gorm:"column:send_num" json:"send_num"`
	DropCount        uint64    `gorm:"column:drop_count" json:"drop_count"`
	DropTime         null.Time `gorm:"column:drop_time" json:"drop_time"`
	Status           uint32    `gorm:"column:status;default:1" json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	AirdropType      string    `gorm:"column:airdrop_type" json:"airdrop_type"`
	ProductSizeId    uint64    `gorm:"column:product_size_id" json:"product_size_id"`
	LogId            string    `gorm:"column:log_id" json:"log_id"`
}

func (*AirdropTaskModel) TableName() string {
	return "airdrop_task"
}

type AirdropTask struct {
	Ctx *gin.Context
}

func (a AirdropTask) AddAirdropTask(payload AirdropTaskModel) (AirdropTaskModel, error) {
	err := cli.HotDogGormDB.WithContext(a.Ctx).Save(&payload).Error
	if err != nil {
		return payload, err
	}
	return payload, err
}
