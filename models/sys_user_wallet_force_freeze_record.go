package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var SysUserWalletForceFreezeRecordDal = &SysUserWalletForceFreezeRecord{}

// 用户钱包强制冻结记录表
type SysUserWalletForceFreezeRecord struct {
	Id            uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	UserId        int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`
	Amount        int64     `gorm:"column:amount;type:bigint(20)" json:"amount"` // 冻结金额(分), 正数冻结负数解冻
	Remark        string    `gorm:"column:remark;type:varchar(255)" json:"remark"`
	OpName        string    `gorm:"column:op_name;type:varchar(255)" json:"op_name"`
	RealName      string    `gorm:"column:real_name;type:varchar(255)" json:"real_name"`
	CanUseBalance int64     `gorm:"column:can_use_balance;type:bigint(20)" json:"can_use_balance"`
	Mobile        string    `gorm:"column:mobile;type:varchar(255)" json:"mobile"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:timestamp" json:"update_time"`
}

func (m *SysUserWalletForceFreezeRecord) TableName() string {
	return "sys_user_wallet_force_freeze_record"
}

// SumFreezeByUserIdsAndTimeRange 按用户ID和时间范围汇总冻结金额(只统计冻结操作)
func (m *SysUserWalletForceFreezeRecord) SumFreezeByUserIdsAndTimeRange(ctx context.Context, userIds []int64, startTime, endTime time.Time) (total int64, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Where("user_id in ?", userIds).
		Where("amount > 0").
		Where("create_time >= ?", startTime).
		Where("create_time <= ?", endTime).
		Select("IFNULL(SUM(amount), 0) as total").
		First(&total).Error
	return
}

// SumFreezeGroupByUserAndTimeRange 按用户ID分组汇总时间范围内的冻结金额
func (m *SysUserWalletForceFreezeRecord) SumFreezeGroupByUserAndTimeRange(ctx context.Context, userIds []int64, startTime, endTime time.Time) (result map[int64]int64, err error) {
	type sumResult struct {
		UserId int64
		Total  int64
	}
	var list []sumResult
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Where("user_id in ?", userIds).
		Where("amount > 0").
		Where("create_time >= ?", startTime).
		Where("create_time <= ?", endTime).
		Select("user_id, IFNULL(SUM(amount), 0) as total").
		Group("user_id").
		Find(&list).Error
	if err != nil {
		return
	}
	result = make(map[int64]int64, len(list))
	for _, item := range list {
		result[item.UserId] = item.Total
	}
	return
}
