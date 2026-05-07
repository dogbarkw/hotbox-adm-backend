package models

import "time"

// AiMatchProductNftArtist undefined
type AiMatchProductNftArtist struct {
	ArtistId     int64     `json:"artist_id" gorm:"artist_id"`
	ArtistName   string    `json:"artist_name" gorm:"artist_name"`
	ArtistBrief  string    `json:"artist_brief" gorm:"artist_brief"`
	AvatarPic    string    `json:"avatar_pic" gorm:"avatar_pic"`
	TitlePic     string    `json:"title_pic" gorm:"title_pic"`
	ArtistFanNum int64     `json:"artist_fan_num" gorm:"artist_fan_num"`
	Weight       int64     `json:"weight" gorm:"weight"`
	OnSaleStatus int64     `json:"on_sale_status" gorm:"on_sale_status"`
	IsDelete     int64     `json:"is_delete" gorm:"is_delete"`
	CreateTime   time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime   time.Time `json:"update_time" gorm:"update_time"`
}

// TableName 表名称
func (*AiMatchProductNftArtist) TableName() string {
	return "ai_match_product_nft_artist"
}
