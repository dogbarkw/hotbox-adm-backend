package form

import "hotbox-adm-backend/dto"

type GetNftCombinationListReq struct {
	PagingReq
	SearchName string `json:"search_name"`
}
type (
	GetNftCombinationMaterialReserveTaskReq struct {
		CombinationId int64 `json:"combination_id" form:"combination_id"  binding:"required"`
	}
	UpdateNftCombinationMaterialReserveTaskReq struct {
		CombinationId int64 `json:"combination_id" form:"combination_id"  binding:"required"`
		dto.CombinationMaterialReserveInfo
	}
	DeleteNftCombinationMaterialReserveTaskReq struct {
		CombinationId int64 `json:"combination_id" form:"combination_id"  binding:"required"`
	}
	CalculateNftCombinationMaterialReserveTaskReq struct {
		CombinationId int64 `json:"combination_id" form:"combination_id"  binding:"required"`
		dto.CombinationMaterialReserveInfo
	}
)

type (
	GetNftDisplaceMaterialReserveTaskReq struct {
		DisplaceId   int64 `json:"displace_id" form:"displace_id"  binding:"required"`
		ActivityType int64 `json:"activity_type" form:"activity_type" binding:"oneof=3 4"`
	}
	UpdateNftDisplaceMaterialReserveTaskReq struct {
		DisplaceId   int64 `json:"displace_id" form:"displace_id"  binding:"required"`
		ActivityType int64 `json:"activity_type" form:"activity_type" binding:"oneof=3 4"`
		dto.DisplaceMaterialReserveInfo
	}
	DeleteNftDisplaceMaterialReserveTaskReq struct {
		DisplaceId   int64 `json:"displace_id" form:"displace_id"  binding:"required"`
		ActivityType int64 `json:"activity_type" form:"activity_type" binding:"oneof=3 4"`
	}
	CalculateNftDisplaceMaterialReserveTaskReq struct {
		DisplaceId   int64 `json:"displace_id" form:"displace_id"  binding:"required"`
		ActivityType int64 `json:"activity_type" form:"activity_type" binding:"oneof=3 4"`
		dto.DisplaceMaterialReserveInfo
	}
)
