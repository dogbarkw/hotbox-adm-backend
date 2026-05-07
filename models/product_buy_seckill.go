package models

// ProductBuySeckill 发售活动
type ProductBuySeckillModel struct {
	Id               uint64 `gorm:"column:id" json:"id"`
	ProductId        uint64 `gorm:"product_id" json:"product_id"`
	NftSizeId        uint64 `gorm:"nft_size_id" json:"nft_size_id"`
	Picture          string `gorm:"column:picture" json:"picture"`
	SecondSaleTime   uint64 `gorm:"second_sale_time " json:"second_sale_time "`
	NewEffectiveDays uint64 `gorm:"new_effective_days" json:"new_effective_days"`
	IsDelete         uint64 `gorm:"is_delete" json:"is_delete"`
	NewFeeRate       uint32 `gorm:"new_fee_rate" json:"new_fee_rate"`
	StartTime        int64  `gorm:"column:start_time" json:"start_time"`
	EndTime          int64  `gorm:"column:end_time" json:"end_time"`
	ResultsTime      uint64 `gorm:"column:results_time" json:"results_time"`
	BuyTime          uint64 `gorm:"column:buy_time" json:"buy_time"`
	SalePrice        string `gorm:"column:sale_price" json:"sale_price"`
	MarketPrice      string `gorm:"column:market_price" json:"market_price"`
	Type             string `gorm:"column:type" json:"type"`
	SaleType         string `gorm:"column:sale_type" json:"sale_type"`
	Count            int64  `gorm:"column:count" json:"count"`
	VirSaleNum       uint64 `gorm:"column:vir_sale_num" json:"vir_sale_num"`
	BoxPropId        uint32 `gorm:"column:box_prop_id" json:"box_prop_id"`
	IsScoreLottery   uint32 `gorm:"column:is_score_lottery" json:"is_score_lottery"`
	LotteryScore     uint32 `gorm:"column:lottery_score" json:"lottery_score"`
	LotteryCodeInfo  string `gorm:"column:lottery_code_info" json:"lottery_code_info"`
	ShelfTime        uint64 `gorm:"column:shelf_time" json:"shelf_time"`
	Status           string `gorm:"column:status" json:"status"`
	MinVersion       uint32 `gorm:"column:min_version" json:"min_version"`
	Unit             string `gorm:"column:unit" json:"unit"`
	Hot              uint32 `gorm:"column:hot" json:"hot"`
	BoxOpenTime      uint64 `gorm:"column:box_open_time" json:"box_open_time"`
	OrderLimit       string `gorm:"column:order_limit" json:"order_limit"`
	Weight           uint32 `gorm:"column:weight" json:"weight"`
	Currency         string `gorm:"column:currency" json:"currency"`
	DestroyProductId uint64 `gorm:"column:destroy_product_id" json:"destroy_product_id"`
	DestroySizeId    uint64 `gorm:"column:destroy_size_id" json:"destroy_size_id"`
	DestroySizeNum   uint32 `gorm:"column:destroy_size_num" json:"destroy_size_num"`
	ApplyNumLimit    uint32 `gorm:"column:apply_num_limit" json:"apply_num_limit"`
	DestroyType      uint32 `gorm:"column:destroy_type" json:"destroy_type"`
}

func (ProductBuySeckillModel) TableName() string {
	return "product_buy_seckill"
}
