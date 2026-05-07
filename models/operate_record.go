package models

import (
	"time"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
)

// AiMatchBackendOperateRecord  后台操作记录表
type AiMatchBackendOperateRecord struct {
	ID          int64     `gorm:"column:id" json:"id"`                     //  主键id
	UserId      int64     `gorm:"column:user_id" json:"user_id"`           //  操作人
	Username    string    `gorm:"column:username" json:"username"`         //  操作人昵称
	Remark      string    `gorm:"column:remark" json:"remark"`             //  操作记录
	Scenes      int64     `gorm:"column:scenes" json:"scenes"`             //  场景  # 1 修改藏品剩余数量/修改理论值 2修改合成库存/自增自减合成库存 / 3变更战国通自动化状态 / 4回收道具 / 67 回收藏品 /  68 空投道具 / 69空投藏品 / 76 摧毁藏品 / 77 操作特殊用户
	AssociateId int64     `gorm:"column:associate_id" json:"associate_id"` //  关联id
	RequestData string    `gorm:"column:request_data" json:"request_data"` //  请求参数
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     //  创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`     //  更新时间
}

func (AiMatchBackendOperateRecord) TableName() string {
	return "ai_match_backend_operate_record"
}

type OperateRecord struct {
	Ctx *gin.Context
}

func (o OperateRecord) CreateRecord(t AiMatchBackendOperateRecord) error {
	return cli.HotDogGormDB.WithContext(o.Ctx).Save(&t).Error
}
