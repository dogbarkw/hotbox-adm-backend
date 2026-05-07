package form

type AddUserAirdropTaskReq struct {
	Count            uint64 `json:"count" binding:"required"`
	Title            string `json:"title,omitempty"`
	PropId           uint64 `json:"prop_id,omitempty"`
	PropName         string `json:"prop_name,omitempty"`
	BenefitType      string `json:"benefit_type,omitempty"`
	Source           uint64 `json:"source,omitempty"` // 1:手动输入名单 2:选择快照数据 3:导入名单文件
	SnapshotDataId   uint64 `json:"snapshot_data_id,omitempty"`
	SnapshotDataName string `json:"snapshot_data_name,omitempty"`
	SendNum          uint64 `json:"send_num,omitempty"`
	DropTime         string `json:"drop_time,omitempty"`
	FileUrl          string `json:"file_url,omitempty"`
	AirdropType      string `binding:"oneof=NFT PROP" json:"airdrop_type"` // NFT-藏品, PROP-道具
	ProductSizeId    uint64 `json:"product_size_id,omitempty"`
}

type AirdropTreasureReq struct {
	Mobile string `json:"mobile" binding:"required"`
	PropId uint64 `json:"prop_id" binding:"required"`
	Number int    `json:"number" binding:"required"`
}
