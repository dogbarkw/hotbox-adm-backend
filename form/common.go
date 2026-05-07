package form

type PagingReq struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}
