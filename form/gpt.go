package form

import (
	"fmt"

	"hotbox-adm-backend/dto"

	"github.com/shopspring/decimal"
)

type GptCollectionArr struct {
	ProductTitle     string `json:"product_title" binding:"required"`
	ProductId        int    `json:"product_id" binding:"required"`
	NftProductSizeId int    `json:"nft_product_size_id" binding:"required"`
	IsMain           bool   `json:"is_main"`           // 是否主材料
	Quantity         *int   `json:"quantity"`          // 指定份数
	ConsumedQuantity *int   `json:"consumed_quantity"` // 消耗份数
	Lp               *int64 `json:"lp" `               // 限价
	Ac               *int64 `json:"ac"`                // 剩余份数
	Pbac             *int64 `json:"pbac"`              // 公池图份数
	Prac             *int64 `json:"prac"`              // 私池图份数
	Ct               *uint  `json:"ct"`                // 成本
}

type GetGptEventPlanningSchemaReq struct {
	CollectionPayload            []GptCollectionArr `json:"collection_payload"  binding:"required,dive"`
	ProductCoefficient           *uint              `json:"product_coefficient" binding:"omitempty"`                          // 产物系数
	CostCeiling                  *uint              `json:"cost_ceiling"`                                                     // 成本上限
	CostAdvice                   *uint              `json:"cost_advice" `                                                     // 成本建议
	ActivityType                 *string            `json:"activity_type" binding:"omitempty,oneof=合成 置换 分解"`                 // 指定活动类型
	IsRaid                       bool               `json:"is_raid"`                                                          // 是否突袭
	IncreaseProfitMultiple       *[2]float64        `json:"increase_profit_multiple"`                                         // 增润倍数
	MinimumGuaranteeFund         *uint              `json:"minimum_guarantee_fund" binding:"omitempty,oneof=2 3 5 6 8 10 20"` // 保底资金
	TemplateType                 uint               `json:"template_type" binding:"required,oneof=1 2 3"`                     // 活动模板, 1:一般活动、2：提纯、3：合流并线
	ProductRecommendedPriceLimit *uint              `json:"product_recommended_price_limit" `                                 // 产物建议限价
	TotalGuaranteeFund           *uint              `json:"total_guarantee_fund"`                                             // 总兜底资金
	ProductLimitRange            *[2]uint           `json:"product_limit_range"`                                              // 产物限价范围
	NewProductTitle              *string            `json:"new_product_title"`                                                // 新产物名称
}

func ValidGetGptEventPlanningSchemaReq(req GetGptEventPlanningSchemaReq) error {
	// for _, v := range req.ActivityType {
	// 	if !lo.Contains(str, v) {
	// 		return fmt.Errorf("活动类型只能是 [合成, 置换, 分解]")
	// 	}
	// }
	if len(req.CollectionPayload) == 0 {
		return fmt.Errorf("CollectionPayload为必填字段")
	}
	hasMainMaterial := false
	for _, v := range req.CollectionPayload {
		if v.IsMain {
			hasMainMaterial = true
		}
	}
	if !hasMainMaterial {
		return fmt.Errorf("至少需要添加一个主材料")
	}

	if req.IncreaseProfitMultiple != nil {
		if req.IncreaseProfitMultiple[0] < 0.9 || req.IncreaseProfitMultiple[0] > 3 {
			return fmt.Errorf("增润倍数前框不得低于0.9,不得高于3.0")
		}
		a := req.IncreaseProfitMultiple[0]
		b := req.IncreaseProfitMultiple[1]
		if decimal.NewFromFloat(b).Sub(decimal.NewFromFloat(a)).LessThan(decimal.NewFromFloat(0.2)) {
			return fmt.Errorf("增润倍数后框-前框必须大于等于0.2")
		}
	}
	return nil
}

type SendGptMsgReq struct {
	Msg []dto.GptMsg `json:"msg" binding:"required"`
}

type GetGptEventPlanningCollectionInfoReq struct {
	ProductId        int `json:"product_id" binding:"required"`
	NftProductSizeId int `json:"nft_product_size_id" binding:"required"`
}

type GenAiArticleMsgReq struct {
	VerifyContent string `json:"verify_content" binding:"required"` // 待审核文本,活动内容主题
}
