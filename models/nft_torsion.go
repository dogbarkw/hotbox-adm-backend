package models

import "time"

var NftTorsionDal = &NftTorsion{}

// NftTorsion undefined
type NftTorsion struct {
	ID                int64     `json:"id" gorm:"id"`                 // 主键id
	ProductId         int64     `json:"product_id" gorm:"product_id"` // 产品id
	OrderId           int64     `json:"order_id" gorm:"order_id"`     // 订单id
	FromUser          int64     `json:"from_user" gorm:"from_user"`   // 来自于
	ToUser            int64     `json:"to_user" gorm:"to_user"`       // 赠予
	TradingHash       string    `json:"trading_hash" gorm:"trading_hash"`
	Status            string    `json:"status" gorm:"status"`     // 状态
	OrderNo           string    `json:"order_no" gorm:"order_no"` // 交易记录
	CreateTime        time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime        time.Time `json:"update_time" gorm:"update_time"`
	IsDelete          int64     `json:"is_delete" gorm:"is_delete"`
	TradeAmount       float64   `json:"trade_amount" gorm:"trade_amount"`
	TransactionHash   string    `json:"transaction_hash" gorm:"transaction_hash"`
	TransactionStatus int64     `json:"transaction_status" gorm:"transaction_status"`
	TrxHash           string    `json:"trx_hash" gorm:"trx_hash"`     // hotdog_chain_交易hash
	TrxStatus         int64     `json:"trx_status" gorm:"trx_status"` // hotdog_chain_交易hash状态
	ExaToUserId       int64     `json:"exa_to_user_id" gorm:"exa_to_user_id"`
	ExaFromUserId     int64     `json:"exa_from_user_id" gorm:"exa_from_user_id"`
}

// TableName 表名称
func (*NftTorsion) TableName() string {
	return "nft_torsion"
}
