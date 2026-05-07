package form

import "time"

type Sweepstakes struct {
	OtherUuid string `json:"other_uuid"  binding:"required" desc:"抽奖id"`
}

type SpinnerHis struct {
	UserId uint64 `json:"user_id" desc:"用户Id" binding:"required"`
}

type FixBoxReq struct {
	UserId uint64 `json:"user_id" desc:"用户Id" binding:"required"`
	Level  int    `json:"level" desc:"等级" binding:"required"`
	BoxNum int    `json:"box_num" desc:"盒子" binding:"required"`
}

type SearchHashReq struct {
	Hash string `json:"hash" desc:"hash" binding:"required"`
}

type NftTorsionOrderRes struct {
	OrderId uint64 `json:"order_id"`
}

type SysUserRes struct {
	LoginName string `json:"login_name"`
}

type NftSizeRes struct {
	ProductTitle string `json:"product_title"`
	ProductId    uint64 `json:"product_id"`
}

type SaleCalendartProductRes struct {
	AuthorName string `json:"author_name"`
}

type SearchNFTHashResp struct {
	ProductName  string    `json:"product_name"`
	Artist       string    `json:"artist"`
	Receiver     string    `json:"receiver"`
	SerialNumber string    `json:"serial_number"`
	ReceiveTime  time.Time `json:"receive_time"`
	TrxHash      string    `json:"trx_hash"`
	BlockHeight  uint64    `json:"block_height"`
	BlockHash    string    `json:"block_hash"`
}

type SearchWalletHashResp struct {
	ProductName  string `json:"product_name"`
	SerialNumber string `json:"serial_number"`
	ConvUrl      string `json:"conv_url"`
	ReceiveTime  string `json:"receive_time"`
	TrxHash      string `json:"trx_hash"`
}

type TransactionResp struct {
	Height    uint64 `gorm:"column:height;comment:'高度'"`
	BlockHash string `gorm:"column:block_hash;comment:'区块哈希'"`
}

type TmpData struct {
	TodayVolume               uint64  `json:"today_volume"`
	TotalMarketCapitalization float32 `json:"total_market_capitalization"`
	VolumeTopGainers          string  `json:"volume_top_gainers"`
}
