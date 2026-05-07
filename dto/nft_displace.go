package dto

type DisplaceMaterialReserveInfo struct {
	DisplaceId         int64                                 `json:"displace_id"`
	ActivityType       int                                   `json:"activity_type"`        // 预留材料任务id
	ExecTime           int64                                 `json:"exec_time"`            // 定时修改时间
	Remark             string                                `json:"remark"`               // 备注
	ActivityReserveNum int64                                 `json:"activity_reserve_num"` // 活动预留份数
	Materials          []DisplaceMaterialReserveMaterialInfo `json:"materials"`
}

type GetNftDisplaceMaterialReserveTaskResp struct {
	DisplaceMaterialReserveInfo
}

type DisplaceMaterialReserveMaterialInfo struct {
	MaterialId          int                               `json:"material_id"`     // 材料组id
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
	TaskList            []MaterialReserveMaterialTaskInfo `json:"task_list"`             // 其他的任务的信息}
}
