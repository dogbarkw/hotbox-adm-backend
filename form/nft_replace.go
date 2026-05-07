package form

type UpdateNftReplaceDisplaceCountReq struct {
	EditRemainCount  *uint `json:"edit_remain_count"`
	AlterRemainCount *int  `json:"alter_remain_count"`
}

type UpdateNftReplaceReserveNumReq struct {
	ReserveNum *int `json:"reserve_num" binding:"required"`
}

type GetDisplaceListReq struct {
	PagingReq
	Id                uint64 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Name              string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Status            uint32 `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
	ActiveType        uint32 `protobuf:"varint,6,opt,name=active_type,json=activeType,proto3" json:"active_type,omitempty"`
	IsSecondaryActive uint32 `protobuf:"varint,7,opt,name=is_secondary_active,json=isSecondaryActive,proto3" json:"is_secondary_active,omitempty"` // 是否是辅活动
	IsUnused          uint32 `protobuf:"varint,8,opt,name=is_unused,json=isUnused,proto3" json:"is_unused,omitempty"`                              // 是否未使用，0表示全部
}
