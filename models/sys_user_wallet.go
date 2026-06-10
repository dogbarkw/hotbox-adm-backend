package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var SysUserWalletDal = &SysUserWallet{}

// 用户钱包
type SysUserWallet struct {
	UserId        int64     `gorm:"column:user_id;type:bigint(20);primary_key;comment:用户user_id" json:"user_id"`
	TotalBalance  float64   `gorm:"column:total_balance;type:decimal(10,2) unsigned;default:0.00;comment:钱包总可用余额" json:"total_balance"`
	FreezeBalance float64   `gorm:"column:freeze_balance;type:decimal(10,2);default:0.00;comment:钱包冻结余额" json:"freeze_balance"`
	ForceFreeze   int64     `gorm:"column:force_freeze;type:bigint(20);default:0;comment:强制冻结金额" json:"force_freeze"`
	UpdateTime    time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *SysUserWallet) TableName() string {
	return "sys_user_wallet"
}

func (m *SysUserWallet) GetUserWallet(ctx context.Context, userId int64) (data SysUserWallet, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Where("user_id", userId).
		First(&data).Error
	return
}

func (m *SysUserWallet) GetUserWallets(ctx context.Context, userIds []int64) (data []SysUserWallet, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Where("user_id in ?", userIds).
		Find(&data).Error
	return
}
