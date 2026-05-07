package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var AiMatchProductNftReplaceDal = &AiMatchProductNftReplace{}

type AiMatchProductNftReplaceModel struct {
	ReplaceId          int64     `gorm:"column:replace_id" db:"replace_id" json:"replace_id" form:"replace_id"`
	ReplaceName        string    `gorm:"column:replace_name" db:"replace_name" json:"replace_name" form:"replace_name"`
	ReplacePic         string    `gorm:"column:replace_pic" db:"replace_pic" json:"replace_pic" form:"replace_pic"`
	HeadPicture        string    `gorm:"column:head_picture" db:"head_picture" json:"head_picture" form:"head_picture"`
	DetailPictures     string    `gorm:"column:detail_pictures" db:"detail_pictures" json:"detail_pictures" form:"detail_pictures"`
	TotalCount         int64     `gorm:"column:total_count" db:"total_count" json:"total_count" form:"total_count"`
	ReplaceCount       int64     `gorm:"column:replace_count" db:"replace_count" json:"replace_count" form:"replace_count"`
	StartTime          int64     `gorm:"column:start_time" db:"start_time" json:"start_time" form:"start_time"`
	EndTime            int64     `gorm:"column:end_time" db:"end_time" json:"end_time" form:"end_time"`
	SaleTime           int64     `gorm:"column:sale_time" db:"sale_time" json:"sale_time" form:"sale_time"`
	Weight             int64     `gorm:"column:weight" db:"weight" json:"weight" form:"weight"`
	MinVersion         int64     `gorm:"column:min_version" db:"min_version" json:"min_version" form:"min_version"`
	OnSaleStatus       int64     `gorm:"column:on_sale_status" db:"on_sale_status" json:"on_sale_status" form:"on_sale_status"`
	CreateTime         time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
	IsDelete           int64     `gorm:"column:is_delete" db:"is_delete" json:"is_delete" form:"is_delete"`
	ToProductTitle     string    `gorm:"column:to_product_title" db:"to_product_title" json:"to_product_title" form:"to_product_title"`
	AppType            int64     `gorm:"column:app_type" db:"app_type" json:"app_type" form:"app_type"`                                                 //  活动指向 1app 2游戏
	SupportBatch       int64     `gorm:"column:support_batch" db:"support_batch" json:"support_batch" form:"support_batch"`                             //  是否支持批量 0=否 1=是
	IsDisplayCount     int64     `gorm:"column:is_display_count" db:"is_display_count" json:"is_display_count" form:"is_display_count"`                 //  是否展示活动数量
	IsDisplayTime      int64     `gorm:"column:is_display_time" db:"is_display_time" json:"is_display_time" form:"is_display_time"`                     //  是否展示活动时间
	ActiveType         int64     `gorm:"column:active_type" db:"active_type" json:"active_type" form:"active_type"`                                     //  活动类型：0置换/1分解 2=抽签置换
	UserMaxTime        int64     `gorm:"column:user_max_time" db:"user_max_time" json:"user_max_time" form:"user_max_time"`                             //  每个用户最多可操作的次数: 1无限制 / 2上限
	UserMaxTimeValue   int64     `gorm:"column:user_max_time_value" db:"user_max_time_value" json:"user_max_time_value" form:"user_max_time_value"`     //  每个用户最多可操作的次数值
	DrawEndTime        time.Time `gorm:"column:draw_end_time" db:"draw_end_time" json:"draw_end_time" form:"draw_end_time"`                             //  抽签截止时间
	PublishTime        time.Time `gorm:"column:publish_time" db:"publish_time" json:"publish_time" form:"publish_time"`                                 //  中签公布时间
	PublishFlag        int64     `gorm:"column:publish_flag" db:"publish_flag" json:"publish_flag" form:"publish_flag"`                                 //  中签是否已公布 1=否 2=是
	RecomPrice         int64     `gorm:"column:recom_price" db:"recom_price" json:"recom_price" form:"recom_price"`                                     //  推荐gas
	OncePrice          int64     `gorm:"column:once_price" db:"once_price" json:"once_price" form:"once_price"`                                         //  每次加多少gas
	DestroyProductId   int64     `gorm:"column:destroy_product_id" db:"destroy_product_id" json:"destroy_product_id" form:"destroy_product_id"`         //  销毁藏品productid
	DestroySizeId      int64     `gorm:"column:destroy_size_id" db:"destroy_size_id" json:"destroy_size_id" form:"destroy_size_id"`                     //  销毁藏品sizeid
	DestroySizeNum     int64     `gorm:"column:destroy_size_num" db:"destroy_size_num" json:"destroy_size_num" form:"destroy_size_num"`                 //  销毁藏品数量
	SubscribeShow      int64     `gorm:"column:subscribe_show" db:"subscribe_show" json:"subscribe_show" form:"subscribe_show"`                         //  是否显示订阅按钮
	IsSecondaryActive  int64     `gorm:"column:is_secondary_active" db:"is_secondary_active" json:"is_secondary_active" form:"is_secondary_active"`     //  是否是辅活动
	UseByMainId        int64     `gorm:"column:use_by_main_id" db:"use_by_main_id" json:"use_by_main_id" form:"use_by_main_id"`                         //  关联的主活动id
	ReferenceActiveId  string    `gorm:"column:reference_active_id" db:"reference_active_id" json:"reference_active_id" form:"reference_active_id"`     //  关联的辅活动id
	QueueSwitch        int64     `gorm:"column:queue_switch" db:"queue_switch" json:"queue_switch" form:"queue_switch"`                                 //  排队开关 0=关闭 1=打开
	QueueDestroyConfig string    `gorm:"column:queue_destroy_config" db:"queue_destroy_config" json:"queue_destroy_config" form:"queue_destroy_config"` //  排队销毁材料配置
	DisplaceType       int64     `gorm:"column:displace_type" db:"displace_type" json:"displace_type" form:"displace_type"`                             //  显示类型：1立即置换/分解 2常驻置换/分解
	IsTreasure         int64     `gorm:"column:is_treasure" db:"is_treasure" json:"is_treasure" form:"is_treasure"`                                     //  是否带有百宝箱
	TreasureDetail     string    `gorm:"column:treasure_detail" db:"treasure_detail" json:"treasure_detail" form:"treasure_detail"`                     //  百宝箱配置
	ReserveNum         int64     `gorm:"column:reserve_num" db:"reserve_num" json:"reserve_num" form:"reserve_num"`                                     //  预留数量
	UserMaxTimeType    int8      `json:"user_max_time_type" gorm:"user_max_time_type"`                                                                  // 0 无 1 拥有配置藏品/道具可置换N次 2 每个用户最多可置换
	IsStrike           int8      `json:"is_strike" gorm:"is_strike"`                                                                                    // 是否是突袭合成(1不是2是)
	UserMaxDoorNum     int64     `json:"user_max_door_num" gorm:"user_max_door_num"`                                                                    // 拥有指定藏品/道具可合成次数(且)
	QueueDestroyType   int8      `json:"queue_destroy_type" gorm:"queue_destroy_type"`                                                                  // 排队消耗类型  1:或，2:且
	Channel            int64     `json:"channel" gorm:"channel"`                                                                                        // 通道顺序 【需求：运营活动列表展示自动排序】
	ReserveNumY        int64     `json:"reserve_num_y" gorm:"reserve_num_y"`                                                                            // Y池预留量
	ReplaceCountY      int64     `json:"replace_count_y" gorm:"replace_count_y"`                                                                        // Y池置换分解数量
}

func (AiMatchProductNftReplaceModel) TableName() string {
	return "ai_match_product_nft_replace"
}

type AiMatchProductNftReplace struct{}

func NewAiMatchProductNftReplace() *AiMatchProductNftReplace {
	return &AiMatchProductNftReplace{}
}

func (a *AiMatchProductNftReplace) UpdateAiMatchProductNftReplaceByReplaceId(ctx context.Context, replaceId int, where, payload map[string]any) (rowAffected int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftReplaceModel{}).Where("replace_id", replaceId)
	for k, v := range where {
		query.Where(k, v)
	}
	query.Updates(payload)
	err = query.Error
	rowAffected = query.RowsAffected
	return
}

func (a *AiMatchProductNftReplace) GetByReplaceId(ctx context.Context, replaceId int) (result AiMatchProductNftReplaceModel, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&result).Where("replace_id", replaceId).First(&result).Error
	return
}

func (a *AiMatchProductNftReplace) GetAiMatchProductNftReplaceList(ctx context.Context, where map[string]any, order []string, limit, offset *int) (list []AiMatchProductNftReplaceModel, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftReplaceModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	if limit != nil {
		query.Limit(*limit)
	}
	if offset != nil {
		query.Offset(*offset)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}

func (a *AiMatchProductNftReplace) GetAiMatchProductNftReplaceByParams(ctx context.Context, where map[string]any, order []string, limit, offset int) (list []AiMatchProductNftReplaceModel, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftReplaceModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	if err := query.Count(&total).Error; err != nil {
		return list, 0, err
	}
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}

func (a *AiMatchProductNftReplace) GetNftReplaceListByParams(ctx context.Context, where map[string][]any) (list []AiMatchProductNftReplace, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftReplaceModel{})
	for k, v := range where {
		query.Where(k, v...)
	}
	err = query.Scan(&list).Error
	return
}
