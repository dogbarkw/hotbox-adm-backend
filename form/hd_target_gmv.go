package form

type (
	HdTargetGmvUpdateReq struct {
		MainId    int64 `json:"main_id"`
		ChildId   int64 `json:"child_id"`
		TargetGmv int64 `json:"target_gmv"`
	}
	HdTargetGmvSwitchReq struct {
		TargetId int64 `json:"target_id" binding:"required"`
		Status   int   `json:"status" binding:"oneof=1 2"` // 1恢复(分账中) 2暂停(正常比例)
	}
	HdTargetGmvQuantRatioUpdateReq struct {
		MainId     int64   `json:"main_id"`
		ChildId    int64   `json:"child_id"`
		QuantRatio float64 `json:"quant_ratio" binding:"gte=0"` // 量化配比倍数，保留两位小数，默认0表示不限制
	}
)
