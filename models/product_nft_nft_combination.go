package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/datatypes"
)

type AiMatchProductNftCombinationModel struct {
	ID                int64          `json:"id" gorm:"id"` // 主键id
	NftMapTitle       string         `json:"nft_map_title" gorm:"nft_map_title"`
	NftMapTitleSub    string         `json:"nft_map_title_sub" gorm:"nft_map_title_sub"`
	NftMapTitlePic    string         `json:"nft_map_title_pic" gorm:"nft_map_title_pic"`
	NftMapTitleSubPic datatypes.JSON `gorm:"column:nft_map_title_sub_pic" db:"nft_map_title_sub_pic" json:"nft_map_title_sub_pic" form:"nft_map_title_sub_pic"`

	ProductId      int64          `json:"product_id" gorm:"product_id"`
	NftSizeId      int64          `json:"nft_size_id" gorm:"nft_size_id"`
	SubProductInfo datatypes.JSON `gorm:"column:sub_product_info" db:"sub_product_info" json:"sub_product_info" form:"sub_product_info"`

	StartTime          int64          `json:"start_time" gorm:"start_time"`
	EndTime            int64          `json:"end_time" gorm:"end_time"`
	SecondSaleTime     int64          `json:"second_sale_time" gorm:"second_sale_time"`
	OnSaleStatus       int64          `json:"on_sale_status" gorm:"on_sale_status"`
	Weight             int64          `json:"weight" gorm:"weight"`
	IsDelete           int64          `json:"is_delete" gorm:"is_delete"`
	CreateTime         time.Time      `json:"create_time" gorm:"create_time"`
	UpdateTime         time.Time      `json:"update_time" gorm:"update_time"`
	CombinationPicture string         `json:"combination_picture" gorm:"combination_picture"` // 合成活动封面图
	CombineType        string         `json:"combine_type" gorm:"combine_type"`
	Version            int64          `json:"version" gorm:"version"`
	SubPropInfo        datatypes.JSON `gorm:"column:sub_prop_info" db:"sub_prop_info" json:"sub_prop_info" form:"sub_prop_info"`

	SendScore            int64          `json:"send_score" gorm:"send_score"` // 合成赠送积分
	NewFeeRate           int64          `json:"new_fee_rate" gorm:"new_fee_rate"`
	NewEffectiveDays     int64          `json:"new_effective_days" gorm:"new_effective_days"`
	ProductType          string         `json:"product_type" gorm:"product_type"`   // 合成产物类型
	PropId               int64          `json:"prop_id" gorm:"prop_id"`             // 道具id
	BoxOpenTime          int64          `json:"box_open_time" gorm:"box_open_time"` // 盲盒开启时间
	Hot                  int64          `json:"hot" gorm:"hot"`
	IsLimitCombine       bool           `json:"is_limit_combine" gorm:"is_limit_combine"`                                      // 是否限制合成数量
	LimitNum             int64          `json:"limit_num" gorm:"limit_num"`                                                    // 限制单次合成数量
	CombineCode          datatypes.JSON `gorm:"column:combine_code" db:"combine_code" json:"combine_code" form:"combine_code"` //  限制单次合成数量
	GenerationType       string         `json:"generation_type" gorm:"generation_type"`                                        // 产物类型
	SendPropId           int64          `json:"send_prop_id" gorm:"send_prop_id"`                                              // 赠送道具id
	SendProductId        int64          `json:"send_product_id" gorm:"send_product_id"`                                        // 奖励藏品 ID
	SendProductSizeId    int64          `json:"send_product_size_id" gorm:"send_product_size_id"`                              // 奖励藏品源文件 ID
	SendNum              int16          `json:"send_num" gorm:"send_num"`                                                      // 奖励数量
	CombineMoreThanNum   int64          `json:"combine_more_than_num" gorm:"combine_more_than_num"`                            // 合成大于等于\n次数
	CombineMoreThanRate  int64          `json:"combine_more_than_rate" gorm:"combine_more_than_rate"`                          // 合成大于等于次数后的概率(0-100)
	AppType              int64          `json:"app_type" gorm:"app_type"`                                                      // 1 app 2 game
	IsDisplayCount       int16          `json:"is_display_count" gorm:"is_display_count"`                                      // 是否展示活动数量
	IsDisplayTime        int16          `json:"is_display_time" gorm:"is_display_time"`                                        // 是否展示活动时间
	TotalTime            int8           `json:"total_time" gorm:"total_time"`                                                  // 合成总次数上限: 1无限制 / 2上限
	TotalTimeValue       int64          `json:"total_time_value" gorm:"total_time_value"`                                      // 合成总次数上限值
	UserMaxTime          int8           `json:"user_max_time" gorm:"user_max_time"`                                            // 每个用户最多可合成次数: 1无限制 / 2上限
	UserMaxTimeValue     int64          `json:"user_max_time_value" gorm:"user_max_time_value"`                                // 每个用户最多可合成次数值
	OriginTotalTimeValue int64          `json:"origin_total_time_value" gorm:"origin_total_time_value"`                        // 合成总次数上限值原始值
	ActivityType         int8           `json:"activity_type" gorm:"activity_type"`                                            // 活动类型 1=普通合成 2=抽签合成
	DrawEndTime          time.Time      `json:"draw_end_time" gorm:"draw_end_time"`                                            // 抽签截止时间
	PublishTime          time.Time      `json:"publish_time" gorm:"publish_time"`                                              // 中签公布时间
	PublishFlag          int8           `json:"publish_flag" gorm:"publish_flag"`                                              // 中签是否已公布 1=否 2=是
	AdvanceReservation   int64          `json:"advance_reservation" gorm:"advance_reservation"`                                // 提前预留量
	RunRestStockStatus   int8           `json:"run_rest_stock_status" gorm:"run_rest_stock_status"`                            // 跑库存状态 1可以跑库存 2正在跑库存 3已跑满库存
	RecomPrice           int64          `json:"recom_price" gorm:"recom_price"`                                                // 推荐gas
	OncePrice            int64          `json:"once_price" gorm:"once_price"`                                                  // 每次加多少gas
	RemainNum            int64          `json:"remain_num" gorm:"remain_num"`                                                  // 剩余合成次数
	RemainNumOri         int64          `json:"remain_num_ori" gorm:"remain_num_ori"`                                          // 剩余次数初始值
	IsNewMode            int8           `json:"is_new_mode" gorm:"is_new_mode"`                                                // 合成模式: 0原生合成 1Flutter合成
	TemporaryReservation int64          `json:"temporary_reservation" gorm:"temporary_reservation"`
	CombineNumLimitType  int8           `json:"combine_num_limit_type" gorm:"combine_num_limit_type"`                                                          // 合成次数上限配置
	SubscribeShow        int64          `json:"subscribe_show" gorm:"subscribe_show"`                                                                          // 是否显示订阅按钮
	ShowType             int8           `json:"show_type" gorm:"show_type"`                                                                                    // 展示方式 1=每日合成 2=常驻合成
	QueueSwitch          int8           `json:"queue_switch" gorm:"queue_switch"`                                                                              // 排队开关 0=关闭 1=打开
	MasterId             int64          `json:"master_id" gorm:"master_id"`                                                                                    // 主活动Id
	QueueDestroyConfig   datatypes.JSON `gorm:"column:queue_destroy_config" db:"queue_destroy_config" json:"queue_destroy_config" form:"queue_destroy_config"` //  排队销毁材料配置
	Cnt                  int64          `json:"cnt" gorm:"cnt"`                                                                                                // 合成产物数量
	IsRandomProduct      int8           `json:"is_random_product" gorm:"is_random_product"`                                                                    // 是否随机合成产物 1 否 2 是
	IsBatchComb          int8           `json:"is_batch_comb" gorm:"is_batch_comb"`
	SlaveId              datatypes.JSON `gorm:"column:slave_id;type:json;comment:子活动Id;NOT NULL" json:"slave_id"`

	IsStrike         int8  `json:"is_strike" gorm:"is_strike"`                   // 是否是突袭合成(1不是2是)
	UserMaxDoorNum   int64 `json:"user_max_door_num" gorm:"user_max_door_num"`   // 拥有指定藏品/道具可合成次数(且)
	QueueDestroyType int8  `json:"queue_destroy_type" gorm:"queue_destroy_type"` // 排队消耗类型  1:或，2:且
	Channel          int64 `json:"channel" gorm:"channel"`                       // 通道顺序 【需求：运营活动列表展示自动排序】
	RemainNumQ       int64 `json:"remain_num_q" gorm:"remain_num_q"`             // Q池剩余合成数量
}

func (*AiMatchProductNftCombinationModel) TableName() string {
	return "ai_match_product_nft_combination"
}

type AiMatchProductNftCombination struct{}

var NftCombinationDal = NewAiMatchProductNftCombination()

var NftCombinationModel = AiMatchProductNftCombinationModel{}

func NewAiMatchProductNftCombination() *AiMatchProductNftCombination {
	return &AiMatchProductNftCombination{}
}

func (a *AiMatchProductNftCombination) GetAiMatchProductNftCombinationList(ctx context.Context, where map[string]any, order []string, limit, offset *int) (list []AiMatchProductNftCombinationModel, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftCombinationModel{})
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

func (a *AiMatchProductNftCombination) One(ctx context.Context, id int64, fields string) (resp AiMatchProductNftCombinationModel, err error) {
	if fields == "" {
		fields = "*"
	}
	err = cli.HotDogGormDB.WithContext(ctx).Model(AiMatchProductNftCombinationModel{}).Select(fields).Where("id", id).First(&resp).Error
	return
}

func (a *AiMatchProductNftCombination) GetAiMatchProductNftCombinationByParams(ctx context.Context, where map[string]any, order []string, limit, offset int) (list []AiMatchProductNftCombinationModel, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftCombinationModel{})
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

func (a *AiMatchProductNftCombination) GetNftCombinationListByParams(ctx context.Context, where map[string][]any) (list []AiMatchProductNftCombinationModel, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftCombinationModel{})
	for k, v := range where {
		query.Where(k, v...)
	}
	err = query.Scan(&list).Error
	return
}

func (a *AiMatchProductNftCombination) GetByParamsLimit(ctx context.Context, where map[string]any, limit int) (list []AiMatchProductNftCombinationModel, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftCombinationModel{})
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Limit(limit).Scan(&list).Error
	return
}
