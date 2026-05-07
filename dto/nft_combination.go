package dto

type AiMatchProductNftCombinationRes struct {
	ID                   int      `json:"id"`
	NftMapTitle          string   `json:"nft_map_title"`
	NftMapTitleSub       string   `json:"nft_map_title_sub"`
	NftMapTitlePic       string   `json:"nft_map_title_pic"`
	NftMapTitleSubPic    []string `json:"nft_map_title_sub_pic"`
	ProductID            int      `json:"product_id"`
	NftSizeID            int      `json:"nft_size_id"`
	SubProductInfo       []any    `json:"sub_product_info"`
	StartTime            int64    `json:"start_time"`
	EndTime              int64    `json:"end_time"`
	SecondSaleTime       int      `json:"second_sale_time"`
	Weight               int      `json:"weight"`
	CreateTime           string   `json:"create_time"`
	UpdateTime           string   `json:"update_time"`
	IsDelete             int      `json:"is_delete"`
	OnSaleStatus         int      `json:"on_sale_status"`
	CombinationPicture   string   `json:"combination_picture"`
	CombineType          string   `json:"combine_type"`
	Version              int      `json:"version"`
	SubPropInfo          []any    `json:"sub_prop_info"`
	SendScore            int      `json:"send_score"`
	NewFeeRate           int      `json:"new_fee_rate"`
	NewEffectiveDays     int      `json:"new_effective_days"`
	BoxOpenTime          int      `json:"box_open_time"`
	Hot                  int      `json:"hot"`
	IsLimitCombine       int      `json:"is_limit_combine"`
	LimitNum             int      `json:"limit_num"`
	CombineCode          []any    `json:"combine_code"`
	GenerationType       string   `json:"generation_type"`
	PropID               int      `json:"prop_id"`
	SendPropID           int      `json:"send_prop_id"`
	CombineMoreThanNum   int      `json:"combine_more_than_num"`
	CombineMoreThanRate  int      `json:"combine_more_than_rate"`
	IsDisplayCount       int      `json:"is_display_count"`
	IsDisplayTime        int      `json:"is_display_time"`
	TotalTime            int      `json:"total_time"`
	TotalTimeValue       int      `json:"total_time_value"`
	UserMaxTime          int      `json:"user_max_time"`
	UserMaxTimeValue     int      `json:"user_max_time_value"`
	OriginTotalTimeValue int      `json:"origin_total_time_value"`
	AppType              int      `json:"app_type"`
	RunRestStockStatus   int      `json:"run_rest_stock_status"`
	AdvanceReservation   int      `json:"advance_reservation"`
	ActivityType         int      `json:"activity_type"`
	PublishTime          string   `json:"publish_time"`
	PublishFlag          int      `json:"publish_flag"`
	DrawEndTime          string   `json:"draw_end_time"`
	RecomPrice           int      `json:"recom_price"`
	OncePrice            int      `json:"once_price"`
	IsNewMode            int      `json:"is_new_mode"`
	TemporaryReservation int      `json:"temporary_reservation"`
	RemainNum            int      `json:"remain_num"`
	SubscribeShow        int      `json:"subscribe_show"`
	HashID               string   `json:"hash_id"`
	PropBenefitType      string   `json:"prop_benefit_type"`
	CombineID            int      `json:"combine_id"`
	StockCount           int      `json:"stock_count"`
	TotalCount           int      `json:"total_count"`
}

type GetNftCombinationMaterialReserveTaskResp struct {
	CombinationMaterialReserveInfo
}

type CombinationMaterialReserveInfo struct {
	CombinationId      int64                                     `json:"combination_id"`       // 预留材料任务id
	ExecTime           int64                                     `json:"exec_time"`            // 定时修改时间
	Remark             string                                    `json:"remark"`               // 备注
	ActivityReserveNum int64                                     `json:"activity_reserve_num"` // 活动预留份数
	MaterialGroups     []CombinationMaterialReserveMaterialGroup `json:"material_groups"`
}

type CombinationMaterialReserveMaterialGroup struct {
	MaterialUuid  string                                   `json:"material_uuid"` // 材料组id
	MaterialInfos []CombinationMaterialReserveMaterialInfo `json:"material_infos"`
}

type CombinationMaterialReserveMaterialInfo struct {
	MaterialUuid        string                            `json:"material_uuid"`   // 材料组id
	MaterialName        string                            `json:"material_name"`   // 材料名称
	AppRemainNum        int64                             `json:"app_remain_num"`  // app内剩余份数
	UserRemainNum       int64                             `json:"user_remain_num"` // 普通用户剩余份数
	MaterialType        string                            `json:"material_type"`   // 材料类型
	ProductId           uint64                            `json:"product_id"`      // 藏品id
	ProductSizeId       uint64                            `json:"product_size_id"`
	PropId              uint64                            `json:"prop_id"`               // 道具id
	ReserveNum          int64                             `json:"reserve_num"`           // 预留份数
	ExpectReserveNum    int64                             `json:"expect_reserve_num"`    // 预期预留份数
	AvailableReserveNum int64                             `json:"available_reserve_num"` // 除去本次活动可扣减的剩余份数
	TaskList            []MaterialReserveMaterialTaskInfo `json:"task_list"`             // 其他的任务的信息
}

type MaterialReserveMaterialTaskInfo struct {
	ActivityType int   `json:"activity_type"`
	IsCountReset bool  `json:"is_count_reset"` // 是否定时更新库存
	ActivityId   int64 `json:"activity_id"`
	ExecTime     int64 `json:"exec_time"`
	ReserveNum   int64 `json:"reserve_num"` // 预留份数
}

type UpdateNftCombinationMaterialReserveTaskResp struct {
	FailMaterialNameList []string `json:"fail_material_name_list"`
}

type CalculateNftCombinationMaterialReserveTaskResp struct {
	CombinationMaterialReserveInfo
}
