package form

type (
	YopTestUserListReq struct {
		Mobile    string `json:"mobile"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		PagingReq
	}
	YopTestUserCheckReq struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	YopTestUserAddReq struct {
		Mobile     string `json:"mobile" binding:"required"`
		UserType   int    `json:"user_type" binding:"oneof=1 2"` // 1 实名 2测试
		Rate       int    `json:"rate" binding:"gte=0,lte=100"`  // 分成比例
		FreezeRate int    `json:"freeze_rate" binding:"gte=0,lte=100"` // 到账冻结比例
		Remark     string `json:"remark"`                        // 备注
	}
	YopTestUserUpdateReq struct {
		Id         int64  `json:"id" binding:"required"`
		Rate       int    `json:"rate" binding:"gte=0,lte=100"` // 分成比例
		FreezeRate int    `json:"freeze_rate" binding:"gte=0,lte=100"` // 到账冻结比例
		Remark     string `json:"remark"`                       // 备注
	}
	YopTestUserDelReq struct {
		Ids []int64 `json:"ids"` // 备注
	}

	YopTestUserIdReq struct {
		Id int64 `json:"id" binding:"required"`
	}
)

type YopTestUserStatDetailReq struct {
	Mouth string `json:"mouth"`
}
