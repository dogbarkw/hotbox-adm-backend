package form

type RecycleRecordListReq struct {
	PagingReq
	Type *int `json:"type"`
}
