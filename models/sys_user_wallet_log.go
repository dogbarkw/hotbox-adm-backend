package models

import (
	"context"

	"hotbox-adm-backend/cli"
)

import (
	"time"
)

var SysUserWalletLogDal = &SysUserWalletLog{}

// 用户钱包流水记录表
type SysUserWalletLog struct {
	Id                  string    `gorm:"column:id;type:varchar(255);primary_key;comment:订单id流水号" json:"id"`
	UserId              int64     `gorm:"column:user_id;type:bigint(20);comment:提现用户id;NOT NULL" json:"user_id"`
	ActionType          uint      `gorm:"column:action_type;type:smallint(5) unsigned;default:1;comment:业务类型，1：提现" json:"action_type"`
	Fee                 float64   `gorm:"column:fee;type:decimal(10,2);default:0.00;comment:变动的金额，正负数。;NOT NULL" json:"fee"`
	OriginalAccountJson string    `gorm:"column:original_account_json;type:varchar(1000);comment:账户变更前的数据 json存储" json:"original_account_json"`
	DisposeAccountJson  string    `gorm:"column:dispose_account_json;type:varchar(1000);comment:账户变更后的数据 json存储" json:"dispose_account_json"`
	State               int       `gorm:"column:state;type:int(11);comment:流水的状态1: 提现中 2: 提现失败 3: 已退回 4: 商品收益 5: 支付宝提现" json:"state"`
	Remark              string    `gorm:"column:remark;type:varchar(255);comment:提现备注" json:"remark"`
	ReviewUserId        int64     `gorm:"column:review_user_id;type:bigint(20);comment:审核用户id" json:"review_user_id"`
	CreateTime          time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:订单创建时间;NOT NULL" json:"create_time"`
	ReviewTime          time.Time `gorm:"column:review_time;type:timestamp;comment:订单审核时间" json:"review_time"`
	UpdateTime          time.Time `gorm:"column:update_time;type:timestamp;default:0000-00-00 00:00:00;comment:订单最后一次更新时间;NOT NULL" json:"update_time"`
	NickName            string    `gorm:"column:nick_name;type:varchar(255);comment:提现到的账户昵称" json:"nick_name"`
	ProductName         string    `gorm:"column:product_name;type:varchar(255);comment:商品名称，特殊情况是用户提现" json:"product_name"`
	AccountUserId       string    `gorm:"column:account_user_id;type:varchar(255);comment:提现到支付宝唯一id" json:"account_user_id"`
	UserMobile          string    `gorm:"column:user_mobile;type:varchar(255);comment:提现用户账号" json:"user_mobile"`
	PayOutBizNo         string    `gorm:"column:pay_out_biz_no;type:varchar(255);comment:提现返回的商户订单号" json:"pay_out_biz_no"`
	PayOrderId          string    `gorm:"column:pay_order_id;type:varchar(255);comment:提现返回的商户订单id" json:"pay_order_id"`
	PayFundOrderId      string    `gorm:"column:pay_fund_order_id;type:varchar(255);comment:支付宝支付资金流水号" json:"pay_fund_order_id"`
	PayStatus           string    `gorm:"column:pay_status;type:varchar(255);comment:转账单据的状态" json:"pay_status"`
	PayTransDate        time.Time `gorm:"column:pay_trans_date;type:timestamp;comment:订单支付时间" json:"pay_trans_date"`
	ActionMethod        int       `gorm:"column:action_method;type:int(11);default:0" json:"action_method"`
	UserBankUuid        string    `gorm:"column:user_bank_uuid;type:varchar(32);default:0" json:"user_bank_uuid"`
	YiRequestId         string    `gorm:"column:yi_request_id;type:varchar(64)" json:"yi_request_id"`
}

func (m *SysUserWalletLog) TableName() string {
	return "sys_user_wallet_log"
}

func (m *SysUserWalletLog) SumTestUserIncome(ctx context.Context, userId int64, startTime time.Time) (total float64, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&m).
		Where("user_id", userId).
		Where("fee > ?", 0).
		Where("create_time >= ?", startTime).
		Select("IFNULL(SUM(fee),0) as total").
		First(&total).Error
	return
}
