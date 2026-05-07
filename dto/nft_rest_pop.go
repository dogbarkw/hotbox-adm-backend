package dto

type NftRestPopSendReq struct {
	ProductId     int `json:"product_id"`
	ProductSizeId int `json:"product_size_id"`
	Limit         int `json:"limit"`
}

type NftRestPopSendResp struct {
	CommonResp
	Data []int `json:"data"`
}
