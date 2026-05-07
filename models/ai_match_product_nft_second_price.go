package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

// AiMatchProductNftSecondPrice undefined
type AiMatchProductNftSecondPrice struct {
	ctx                        context.Context
	ID                         int64     `json:"id" gorm:"id"`                                   // 主键id
	ProductId                  int64     `json:"product_id" gorm:"product_id"`                   // 商品id
	NftProductSizeId           int64     `json:"nft_product_size_id" gorm:"nft_product_size_id"` // 款式id
	SaleMinPrice               float64   `json:"sale_min_price" gorm:"sale_min_price"`           // 最低售价
	SaleMaxPrice               float64   `json:"sale_max_price" gorm:"sale_max_price"`           // 最高售价
	CreateTime                 time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime                 time.Time `json:"update_time" gorm:"update_time"`
	IsDelete                   int64     `json:"is_delete" gorm:"is_delete"`
	ActiveTime                 int64     `json:"active_time" gorm:"active_time"`                   // 对应配置活动的活动开始时间
	BargainTime                int64     `json:"bargain_time" gorm:"bargain_time"`                 // 定时改价时间
	ActivePicture              string    `json:"active_picture" gorm:"active_picture"`             // 款式对应展示的豆腐块视图
	SecondSaleTime             int64     `json:"second_sale_time" gorm:"second_sale_time"`         // 对应配置活动的二级上架时间
	ProductTitle               string    `json:"product_title" gorm:"product_title"`               // 款式对应展示的显示名称
	SellPriceMaxLimit          float64   `json:"sell_price_max_limit" gorm:"sell_price_max_limit"` // 挂单出价最高限制
	AppearForm                 string    `json:"appear_form" gorm:"appear_form"`                   // 出现形式
	MarketType                 string    `json:"market_type" gorm:"market_type"`                   // 市场类型
	ConsignStatus              int8      `json:"consign_status" gorm:"consign_status"`             // 寄售状态 1=寄售中
	Hot                        int64     `json:"hot" gorm:"hot"`                                   // 热度
	OnSaleStatus               int16     `json:"on_sale_status" gorm:"on_sale_status"`             // 上下架：0下架1上架
	OnlineTime                 time.Time `json:"online_time" gorm:"online_time"`
	OfflineTime                time.Time `json:"offline_time" gorm:"offline_time"`
	SupportDeposit             int8      `json:"support_deposit" gorm:"support_deposit"`                       // 是否支持定金模式  0=否 1=是
	IsCanBatchBuy              int8      `json:"is_can_batch_buy" gorm:"is_can_batch_buy"`                     // 是否支持批量购买 0=否 1=是
	AboveMaxPriceSecondOff     int64     `json:"above_max_price_second_off" gorm:"above_max_price_second_off"` // 二级超过限价订单是否，自动下架1否，2是
	SellPriceMinLimit          float64   `json:"sell_price_min_limit" gorm:"sell_price_min_limit"`             // 挂单出价最低限制
	RealUserSurplus            float64   `json:"real_user_surplus" gorm:"real_user_surplus"`                   // 真实用户剩余份数
	AllUserSurplus             float64   `json:"all_user_surplus" gorm:"all_user_surplus"`                     // 用户剩余份数
	DisplayOption              int8      `json:"display_option" gorm:"display_option"`                         // nft详情页:1=剩余份数 2=流通份数 3=剩余份数，流通分数都展示 4=都不展示
	IsCanBatchSell             int8      `json:"is_can_batch_sell" gorm:"is_can_batch_sell"`                   // 支持批量出售:0=不能 1=能
	TheoreticalValues          string    `json:"theoretical_values" gorm:"theoretical_values"`
	RemainCount                int64     `json:"remain_count" gorm:"remain_count"`                                     // 流通数量
	CanBatchBuyTime            time.Time `json:"can_batch_buy_time" gorm:"can_batch_buy_time"`                         // 定时开启批量时间
	CanNoBatchBuyTime          time.Time `json:"can_no_batch_buy_time" gorm:"can_no_batch_buy_time"`                   // 定时关闭批量时间
	YestClosePrice             int64     `json:"yest_close_price" gorm:"yest_close_price"`                             // 昨日收盘价
	CountResetTime             time.Time `json:"count_reset_time" gorm:"count_reset_time"`                             // 定时修改库存时间
	CountResetValue            int64     `json:"count_reset_value" gorm:"count_reset_value"`                           // 定时修改库存值
	IsShowHoldDistributionRate int16     `json:"is_show_hold_distribution_rate" gorm:"is_show_hold_distribution_rate"` // 展示藏品持有分布,1=展示,2=不展示
	HoldDistributionRateMidSet int64     `json:"hold_distribution_rate_mid_set" gorm:"hold_distribution_rate_mid_set"` // 设定藏品持有分布的中间值,0-10000,使用时转为100%格式 = 值/100 %
	UserPercentage             float64   `json:"user_percentage" gorm:"user_percentage"`
}

// TableName 表名称
func (*AiMatchProductNftSecondPrice) TableName() string {
	return "ai_match_product_nft_second_price"
}

func NewAiMatchProductNftSecondPrice(ctx context.Context) *AiMatchProductNftSecondPrice {
	return &AiMatchProductNftSecondPrice{ctx: ctx}
}

func (a *AiMatchProductNftSecondPrice) UpdateByParams(where, payload map[string]any) (affectRow int64, err error) {
	query := cli.HotDogGormDB.WithContext(a.ctx).Model(&AiMatchProductNftSecondPrice{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	query = query.Updates(payload)
	affectRow = query.RowsAffected
	err = query.Error
	return
}

func (a *AiMatchProductNftSecondPrice) GetByParams(where map[string]any) (list []AiMatchProductNftSecondPrice, err error) {
	query := cli.HotDogGormDB.WithContext(a.ctx).Model(&AiMatchProductNftSecondPrice{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	query = query.Scan(&list)
	err = query.Error
	return
}

func (a *AiMatchProductNftSecondPrice) GetByProductIdAndNftProductSizeId(productId, nftProductSizeId int) (data AiMatchProductNftSecondPrice, err error) {
	err = cli.HotDogGormDB.WithContext(a.ctx).Model(a).
		Where("product_id", productId).Where("nft_product_size_id", nftProductSizeId).First(&data).Error
	return
}

func (a AiMatchProductNftSecondPrice) GetListByParams(ctx context.Context, where map[string][]any) (list []AiMatchProductNftSecondPrice, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v...)
	}
	err = query.Scan(&list).Error
	return
}

func (a AiMatchProductNftSecondPrice) GetOneByParams(ctx context.Context, where map[string][]any) (result AiMatchProductNftSecondPrice, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v...)
	}
	err = query.First(&result).Error
	return
}
