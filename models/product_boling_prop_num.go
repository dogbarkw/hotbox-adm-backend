package models

import "time"

type AiMatchProductBolingPropNumModel struct {
	BolingUuid string    `gorm:"column:boling_uuid" json:"boling_uuid"`
	UserId     int64     `gorm:"column:user_id" json:"user_id"`
	PropId     int64     `gorm:"column:prop_id" json:"prop_id"`
	Number     int64     `gorm:"column:number" json:"number"`
	FreezeNum  int64     `gorm:"column:freeze_num" json:"freeze_num"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	IsDelete   int8      `gorm:"column:is_delete;default:0" json:"is_delete"`
}

func (AiMatchProductBolingPropNumModel) TableName() string {
	return "ai_match_product_boling_prop_num"
}
