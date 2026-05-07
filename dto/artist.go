package dto

type ArtistNftNum struct {
	ArtistId   int64  `json:"artist_id"`
	ArtistName string `json:"artist_name"`
	NftNum     int    `json:"nft_num"`
}

type CombinationMaterial struct {
	Name                        string `json:"name"`
	Type                        string `json:"type"`
	Serial                      int64  `json:"serial"`
	Picture                     string `json:"picture"`
	PropID                      int64  `json:"prop_id"`
	NeedNum                     int64  `json:"need_num"`
	SubProductID                int64  `json:"sub_product_id"`
	NftProductSizeID            int64  `json:"nft_product_size_id"`
	InscriptionProductID        int64  `json:"inscription_product_id"`
	InscriptionProductNftSizeID int64  `json:"inscription_product_nft_size_id"`
}

type UpgradeMaterial struct {
	Serial        int64 `json:"serial"`
	ProductID     int64 `json:"product_id"`
	ProductSizeID int64 `json:"product_size_id"`
}
