package form

type GmvListReq struct {
	PagingReq
	YMD       string `json:"ymd"`       // 可选，查询日期,具体到某天
	DateRange string `json:"ymd_range"` // 可选，查询日期范围,起始日期和结束日期用逗号分开
}
