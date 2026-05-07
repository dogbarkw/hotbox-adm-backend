package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"gorm.io/datatypes"
)

type AiMatchProductNftPropModel struct {
	PropId                   int64          `json:"prop_id" gorm:"prop_id"`
	PropType                 string         `json:"prop_type" gorm:"prop_type"`       // 道具类型：photo图片/video视频
	BenefitType              string         `json:"benefit_type" gorm:"benefit_type"` // 权益类型：first_buy优先购 discount手续减免
	Discount                 float64        `json:"discount" gorm:"discount"`         // 折扣
	ActiveId                 int64          `json:"active_id" gorm:"active_id"`       // 优先购活动id
	ProductTitle             string         `json:"product_title" gorm:"product_title"`
	SubTitle                 string         `json:"sub_title" gorm:"sub_title"`
	HeadPicture              string         `json:"head_picture" gorm:"head_picture"` // 头图
	FileAddress              string         `json:"file_address" gorm:"file_address"` // 视频地址
	PropConfig               datatypes.JSON `json:"prop_config" gorm:"prop_config"`
	PropDetailPictures       datatypes.JSON `json:"prop_detail_pictures" gorm:"prop_detail_pictures"`
	PropDesc                 string         `json:"prop_desc" gorm:"prop_desc"`
	ArtistId                 int64          `json:"artist_id" gorm:"artist_id"`
	OnSaleStatus             int64          `json:"on_sale_status" gorm:"on_sale_status"`
	Extend                   string         `json:"extend" gorm:"extend"`               // 标签
	PriorityHour             int64          `json:"priority_hour" gorm:"priority_hour"` // 优先时间
	StartTime                int64          `json:"start_time" gorm:"start_time"`
	EndTime                  int64          `json:"end_time" gorm:"end_time"`
	CreateTime               time.Time      `json:"create_time" gorm:"create_time"`
	UpdateTime               time.Time      `json:"update_time" gorm:"update_time"`
	IsDelete                 int64          `json:"is_delete" gorm:"is_delete"`
	ListTitle                string         `json:"list_title" gorm:"list_title"`
	ListPicture              string         `json:"list_picture" gorm:"list_picture"`
	AvailableDay             int64          `json:"available_day" gorm:"available_day"`                               // 可用时长(天)
	StockCount               int64          `json:"stock_count" gorm:"stock_count"`                                   // 剩余库存
	TotalCount               int64          `json:"total_count" gorm:"total_count"`                                   // 总库存
	NeedPointNum             int64          `json:"need_point_num" gorm:"need_point_num"`                             // 消耗积分数
	IsConsume                int8           `json:"is_consume" gorm:"is_consume"`                                     // 是否消耗藏品 0:否 1:是
	Level                    int64          `json:"level" gorm:"level"`                                               // 版权盲盒等级
	RandomEndTime            int64          `json:"random_end_time" gorm:"random_end_time"`                           // 抽奖券截止时间
	RandomLink               string         `json:"random_link" gorm:"random_link"`                                   // 抽奖券跳转链接
	RandomPropId             int64          `json:"random_prop_id" gorm:"random_prop_id"`                             // 大转盘活动id
	SweepWheelTitle          string         `json:"sweep_wheel_title" gorm:"sweep_wheel_title"`                       // 页面导航标题
	SweepWheelRuleLink       string         `json:"sweep_wheel_rule_link" gorm:"sweep_wheel_rule_link"`               // 活动规则\n链接
	SweepWheelRuleText       string         `json:"sweep_wheel_rule_text" gorm:"sweep_wheel_rule_text"`               // 活动规则\n文案
	SweepWheelNftProductId   int64          `json:"sweep_wheel_nft_product_id" gorm:"sweep_wheel_nft_product_id"`     // 消耗的藏品id
	SweepWheelNftSizeId      int64          `json:"sweep_wheel_nft_size_id" gorm:"sweep_wheel_nft_size_id"`           // 消耗的藏品size id
	SweepWheelNftProductName string         `json:"sweep_wheel_nft_product_name" gorm:"sweep_wheel_nft_product_name"` // 藏品名称
	SweepWheelType           int64          `json:"sweep_wheel_type" gorm:"sweep_wheel_type"`                         // 抽奖类型0消耗抽奖券,1 消耗藏品
	BoxStartTs               int64          `json:"box_start_ts" gorm:"box_start_ts"`                                 // 百宝箱开启时间
	IsMergeTime              int8           `json:"is_merge_time" gorm:"is_merge_time"`                               // 是否合并道具时间: 0不合并1合并
	SweepWheelNum            int64          `json:"sweep_wheel_num" gorm:"sweep_wheel_num"`                           // 消耗数量
	TwistedNeedType          int64          `json:"twisted_need_type" gorm:"twisted_need_type"`                       // 扭蛋消耗类型0：积分1：道具2：藏品
	TwistedNeedId            int64          `json:"twisted_need_id" gorm:"twisted_need_id"`                           // 扭蛋消耗对应id
	TwistedNeedSubId         int64          `json:"twisted_need_sub_id" gorm:"twisted_need_sub_id"`                   // 扭蛋消耗对应副id
	BubbleText               string         `json:"bubble_text" gorm:"bubble_text"`                                   // 气泡文案
	BonusId                  int64          `json:"bonus_id" gorm:"bonus_id"`                                         // 十连奖励的奖品id
	BonusCount               int64          `json:"bonus_count" gorm:"bonus_count"`                                   // 十连奖励的奖品数量
	IsCanTransfer            int8           `json:"is_can_transfer" gorm:"is_can_transfer"`                           // 是否可以转赠
	TransferInfo             datatypes.JSON `json:"transfer_info" gorm:"transfer_info"`                               // 转赠条件
	AvatarValidityHours      int64          `json:"avatar_validity_hours" gorm:"avatar_validity_hours"`               // 头像权益有效期(小时)
	AvatarConfig             datatypes.JSON `json:"avatar_config" gorm:"avatar_config"`                               // 寄售头像框
	IsMinProtect             int8           `json:"is_min_protect" gorm:"is_min_protect"`                             // 是否开启保底模式: 1不开启2开启
}

func (AiMatchProductNftPropModel) TableName() string {
	return "ai_match_product_nft_prop"
}

type AiMatchProductNftProp struct{}

func NewAiMatchProductNftProp() *AiMatchProductNftProp {
	return &AiMatchProductNftProp{}
}

func (a *AiMatchProductNftProp) GetSaleCalendarProductSizeByParams(ctx context.Context, where map[string]any, order []string, limit *int) (list []AiMatchProductNftPropModel, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftPropModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}

func (*AiMatchProductNftProp) One(ctx context.Context, propId uint64, fields string) (res AiMatchProductNftPropModel, err error) {
	if fields == "" {
		fields = "*"
	}
	err = cli.HotDogGormDB.WithContext(ctx).
		Select(fields).
		Where("prop_id = ?", propId).
		Where("is_delete = 0").
		First(&res).Error
	return res, err
}
