package dto

type NftReplaceListData struct {
	Id                 uint64   `json:"id"` // 活动ID
	Name               string   `json:"name"`
	CoverImg           string   `json:"cover_img"` // 封面图
	Images             []string `json:"images"`    // 详情页图片
	TotalNum           uint32   `json:"total_num"`
	DisplaceNum        uint32   `json:"displace_num"`
	MinVersion         uint32   `json:"min_version"`
	Status             uint32   `json:"status"`
	ActiveType         uint32   `json:"active_type"`
	DrawEndTime        string   `json:"draw_end_time"`
	PublishTime        string   `json:"publish_time"`
	Eid                string   `json:"eid"`                  // 加密 ID
	IsSecondaryActive  uint32   `json:"is_secondary_active"`  // 是否是辅活动
	UsedMainActive     uint64   `json:"used_main_active"`     // 被使用的主活动
	ReferenceActiveIds []uint64 `json:"reference_active_ids"` // 被关联的活动
	ReserveNum         int64    `json:"reserve_num"`
}
