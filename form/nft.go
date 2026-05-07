package form

type SecondPriceListReq struct {
	PagingReq
	ProductId        int    `json:"product_id"`
	ProductTitle     string `json:"product_title"`
	OnSaleStatus     int    `json:"on_sale_status"`
	NftProductSizeId int    `json:"nft_product_size_id"`
}

type GetSecondPriceFlushSurplusReq struct {
	ProductId        int  `json:"product_id" form:"product_id"  binding:"required"`
	NftProductSizeId *int `json:"nft_product_size_id" form:"nft_product_size_id"  binding:"required"`
}

type SecondPriceFlushSurplusReq struct {
	ProductId        int `json:"product_id,omitempty" binding:"required"`
	NftProductSizeId int `json:"nft_product_size_id"`
	FlushType        int `json:"flush_type" binding:"oneof=1 2 3 4" `
}

type NftRecyclingByCountReq struct {
	Mobile           string `json:"mobile" binding:"required"`
	ProductID        int    `json:"product_id" binding:"required"`
	NftProductSizeID int    `json:"nft_product_size_id"`
	Count            int    `json:"count" binding:"required"`
}
type NftAirdropByCountReq struct {
	Mobile           string `json:"mobile" binding:"required"`
	ProductID        int    `json:"product_id" binding:"required"`
	NftProductSizeID int    `json:"nft_product_size_id"`
	Count            int    `json:"count" binding:"required"`
}

type NftSimpleListReq struct {
	PagingReq
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	IsBusiness uint32 `json:"is_business"`
}

type NftRestListReq struct {
	PagingReq
	ProductId        int `json:"product_id" binding:"required"`
	NftProductSizeId int `json:"nft_product_size_id"`
}

type NftDestroyByCountReq struct {
	Mobile           string `json:"mobile" binding:"required"`
	ProductID        int    `json:"product_id" binding:"required"`
	NftProductSizeID int    `json:"nft_product_size_id" binding:"required"`
	Count            int    `json:"count" binding:"required"`
	Remark           string `json:"remark" binding:"required"`
}

type NftSecondPriceFlushUserPercentageReq struct {
	ProductId        int `json:"product_id,omitempty" binding:"required"`
	NftProductSizeId int `json:"nft_product_size_id"`
}
