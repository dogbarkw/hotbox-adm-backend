package dto

type ProductNftSecondPriceJoinSaleCalendar struct {
	Id               uint64 `json:"id"`
	ProductId        uint64 `json:"product_id"`
	NftProductSizeId uint64 `json:"nft_product_size_id"`
	ProductTitle     string `json:"product_title"`

	SaleMinPrice      float32 `json:"sale_min_price"`
	SaleMaxPrice      float32 `json:"sale_max_price"`
	SellPriceMaxLimit float32 `json:"sell_price_max_limit"`
	SellPriceMinLimit float32 `json:"sell_price_min_limit"`
	SecondSaleTime    uint64  `json:"second_sale_time"`
	ActivePicture     string  `json:"active_picture"`

	RealUserSurplus   float32 `gorm:"column:real_user_surplus" json:"real_user_surplus"`
	AllUserSurplus    float32 `json:"all_user_surplus"`
	TheoreticalValues string  `json:"theoretical_values"`
	CountResetTime    string  `gorm:"column:count_reset_time;type:datetime;not null;default:'0001-01-01 00:00:00'" json:"count_reset_time,omitempty"`
	CountResetValue   int     `json:"count_reset_value"`
	ArtistID          int     `json:"artist_id"`
	AuthorName        string  `json:"author_name"`
	NftType           string  `json:"nft_type"`
	IsTest            int     `json:"is_test"`
	NftCount          int     `json:"nft_count"`
	TotalRemainCount  int     `json:"total_remain_count"`
	RemainCount       int     `json:"remain_count"`
	TotalCount        int     `json:"total_count"`
	UserPercentage    float64 `json:"user_percentage"`
}

type AiMatchProductNftCombination struct {
	AppType        int
	IsLimitCombine int
}

type SaleProductNFTSizeModel struct {
	TotalCount       int
	NftProductSizeId int
	IsDelete         int
	StockCount       int
}

type SaleCalendarProductSizeModel struct {
	ProductId  int
	IsDelete   int
	TotalCount int
}

type ActivityUpgradeModel struct {
	ProductId int
	Status    int
}

type ProductNftSecondPriceJoinSaleCalendarReserveTask struct {
	ExecTime     int64  `json:"exec_time"`     // 执行时间
	ReserveNum   int64  `json:"reserve_num"`   // 预留份数
	TaskSource   int    `json:"task_source"`   // 任务来源 1=手动 2=活动
	ActivityType int    `json:"activity_type"` // 活动类型 1=合成 3=置换 4=分解（仅活动任务返回）
	ActivityId   int64  `json:"activity_id"`   // 活动ID（仅活动任务返回）
	OperatorName string `json:"operator_name"` // 操作人昵称（仅手动任务返回）
}

type NftSecondPriceListRespItem struct {
	ProductNftSecondPriceJoinSaleCalendar
	TaskList []ProductNftSecondPriceJoinSaleCalendarReserveTask `json:"task_list"`
}
