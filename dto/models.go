package dto

type Artist struct {
	ArtistId   int    `json:"artist_id"`
	ArtistName string `json:"artist_name"`
}

type Materials struct {
	MaterialId       int    `json:"material_id"`
	MaterialUuid     string `json:"material_uuid"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ProductId        int    `json:"product_id"`
	NftProductSizeId int    `json:"nft_product_size_id"`
	PropId           int    `json:"prop_id"`
	Num              int    `json:"num"`           // 材料所需要数量
	MaterialNum      uint32 `json:"material_num"`  // 合成一次活动所需要的材料数,只有MaterialType为0时才有
	MaterialType     uint8  `json:"material_type"` // 材料类型；0：且，1：或
}

type OutPutProduct struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ProductId        int    `json:"product_id"`
	NftProductSizeId int    `json:"nft_product_size_id"`
	PropId           int    `json:"prop_id"`
	Num              int    `json:"num"`
	ArtistData       Artist `json:"artist_data"`
}

type (
	MaterialsArr           []Materials
	ActivityScoreMaterials struct {
		ActivityId    int `json:"activity_id"`
		ActivityType  int `json:"activity_type"`
		MaterialsData struct {
			Materials []MaterialsArr `json:"materials"`
		} `json:"materials_data"`
	}
)
