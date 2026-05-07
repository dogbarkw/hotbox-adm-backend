package models

import "time"

type AiMatchProductOtherPropNumModel struct {
	OtherUuid   string    `gorm:"column:other_uuid" json:"other_uuid"`
	BenefitType string    `gorm:"column:benefit_type" json:"benefit_type"`
	UserId      uint64    `gorm:"column:user_id" json:"user_id"`
	PropId      int64     `gorm:"column:prop_id" json:"prop_id"`
	Number      uint32    `gorm:"column:number" json:"number"`
	Version     uint32    `gorm:"column:version" json:"version"`
	IsDelete    uint32    `gorm:"column:is_delete" json:"is_delete"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}
