package form

import (
	"hotbox-adm-backend/dto"
)

type GetPendingActivityScoreListReq struct {
	PagingReq
	ActivityType int `json:"activity_type" binding:"oneof=1 3 4"`
}

type GetEndedActivityScoreListReq struct {
	PagingReq
	ActivityType int `json:"activity_type" binding:"oneof=1 3 4"`
	OrderType    int `json:"order_type" binding:"oneof=0 1 2"` // 排序，默认0,1 按开始时间降序, 2按报警次数降序
}

type PendingActivityScoreListResp struct {
	ID                               uint                       `gorm:"primarykey;type:int(0)" json:"id" redis:"id"`
	ActivityId                       int                        `gorm:"activity_id; uniqueIndex:uniq_idx_activity_id_type" json:"activity_id"`
	ActivityType                     int                        `gorm:"activity_type; uniqueIndex:uniq_idx_activity_id_type" json:"activity_type"`
	ActivityTitle                    string                     `gorm:"activity_title" json:"activity_title"`
	ActivityStartTs                  int64                      `gorm:"activity_start_ts" json:"activity_start"`
	Artist                           []dto.Artist               `gorm:"artist" json:"artist"`
	Materials                        dto.ActivityScoreMaterials `gorm:"materials" json:"material"`                                                      // 材料
	OutPutProduct                    []dto.OutPutProduct        `gorm:"out_put_product" json:"out_put_product"`                                         // 产物
	TotalCost                        float64                    `gorm:"total_cost" json:"total_cost"`                                                   // 总成本
	AntiFrictionLine                 float64                    `gorm:"anti_friction_line" json:"anti_friction_line"`                                   // 反撸线
	ExpectedProductCirculationAmount float64                    `gorm:"expected_product_circulation_amount" json:"expected_product_circulation_amount"` // 产物预期流通份数
	ExpectedMarketProductValue       float64                    `gorm:"expected_market_product_value" json:"expected_market_product_value"`             // 产物预期流通市值
	Score                            float64                    `json:"score"`
	RealProductCirculationAmount     float64                    `json:"real_product_circulation_amount"` // 产物实际流通份数
	RealMarketProductValue           float64                    `json:"real_market_product_value"`       // 产物实际流通市值
	Duration                         int64                      `json:"duration"`                        // 活动持续时长
	SaleMinPrice                     float64                    `json:"sale_min_price"`                  // 产物当前最低挂售价
	SellPriceMinLimit                float64                    `json:"sell_price_min_limit"`            // 产物当前限价
	ProductAvgCost                   float64                    `json:"product_avg_cost"`                // 产物当日平均成交价
	WarnTime                         int                        `gorm:"warn_time" json:"warn_time"`      // 警告次数
	RealScore                        float64                    `gorm:"real_score" json:"real_score"`    // 实时评分
	NotchNumber                      int64                      `json:"notch_number"`                    // 缺口数
}

type UpdatePendingActivityScoreReq struct {
	Id        int64    `json:"id" binding:"required"`
	TotalCost *float64 `json:"total_cost"`
}

type UpdateArtistRecommendScoreReq struct {
	Id    int64  `json:"id" binding:"required"`
	Score string `json:"score"`
}
