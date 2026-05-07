package models

import (
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type SaleCalendarProductModel struct {
	Id                    uint64         `json:"id"`
	ArtNo                 string         `json:"art_no" gorm:"art_no"`
	ProductName           string         `json:"product_name" gorm:"product_name"`
	ProductAbbreviation   string         `json:"product_abbreviation" gorm:"product_abbreviation"`
	Price                 string         `json:"price" gorm:"price"`
	Currency              string         `json:"currency" gorm:"currency"`
	SaleTime              int64          `json:"sale_time" gorm:"sale_time"`
	Pictures              string         `json:"pictures" gorm:"pictures"`
	Links                 string         `json:"links" gorm:"links"`
	Status                int64          `json:"status" gorm:"status"`
	Operator              int64          `json:"operator" gorm:"operator"`
	Extend                string         `json:"extend" gorm:"extend"`
	CreateTime            time.Time      `json:"create_time" gorm:"create_time"`
	UpdateTime            time.Time      `json:"update_time" gorm:"update_time"`
	IsDelete              int64          `json:"is_delete" gorm:"is_delete"`
	SaleDate              time.Time      `json:"sale_date" gorm:"sale_date"`
	NewPictures           string         `json:"new_pictures" gorm:"new_pictures"`
	IconUrl               string         `json:"icon_url" gorm:"icon_url"` // 图标
	ClickCount            int64          `json:"click_count" gorm:"click_count"`
	DistributorCount      int64          `json:"distributor_count" gorm:"distributor_count"`
	Type                  string         `json:"type" gorm:"type"`
	ShowType              string         `json:"show_type" gorm:"show_type"`
	BrandName             string         `json:"brand_name" gorm:"brand_name"`
	Style                 string         `json:"style" gorm:"style"`
	InformationCount      int64          `json:"information_count" gorm:"information_count"`
	HeatDecayCount        int64          `json:"heat_decay_count" gorm:"heat_decay_count"`
	Classification        int64          `json:"classification" gorm:"classification"`
	EstimateCount         int64          `json:"estimate_count" gorm:"estimate_count"`
	DetailDesc            string         `json:"detail_desc" gorm:"detail_desc"`                         // 商品详情
	DressingPictures      string         `json:"dressing_pictures" gorm:"dressing_pictures"`             // 穿搭图片
	IsBusiness            int64          `json:"is_business" gorm:"is_business"`                         // 是否商家商品
	BrandId               int64          `json:"brand_id" gorm:"brand_id"`                               // 品牌id
	IsMerchantProvide     int64          `json:"is_merchant_provide" gorm:"is_merchant_provide"`         // 是否由合作商家提供
	ShipmentTime          string         `json:"shipment_time" gorm:"shipment_time"`                     // 商品发货时间：小时
	AfterSales            string         `json:"after_sales" gorm:"after_sales"`                         // 售后内容
	SoldCount             int64          `json:"sold_count" gorm:"sold_count"`                           // 已售出数量
	StylePicture          string         `json:"style_picture" gorm:"style_picture"`                     // 版本信息
	SizeComparisonPicture string         `json:"size_comparison_picture" gorm:"size_comparison_picture"` // 尺码对照表
	Weight                int64          `json:"weight" gorm:"weight"`                                   // 权重
	ThemeId               int64          `json:"theme_id" gorm:"theme_id"`                               // 主题id
	OnSaleStatus          int64          `json:"on_sale_status" gorm:"on_sale_status"`                   // 上下架   1上架 - 0未上架
	FreightAmount         float64        `json:"freight_amount" gorm:"freight_amount"`                   // 邮费
	Supplier              string         `json:"supplier" gorm:"supplier"`                               // 供应商
	IsHighSpider          int64          `json:"is_high_spider" gorm:"is_high_spider"`                   // 是否高频爬取
	InscriptionConfig     datatypes.JSON `gorm:"inscription_config" json:"inscription_config"`
	CouponIds             datatypes.JSON `json:"coupon_ids"`
	DetailHeadDisplayPic  datatypes.JSON `gorm:"column:detail_head_display_pic" json:"detail_head_display_pic"`
	CopyrightTopFile      datatypes.JSON `gorm:"column:copyright_top_file" json:"copyright_top_file"`
	ConsumeNftConfig      datatypes.JSON `gorm:"column:consume_nft_config" json:"consume_nft_config"`

	Specification           string  `json:"specification" gorm:"specification"`
	NftType                 string  `json:"nft_type" gorm:"nft_type"`               // nft类型
	AuthorName              string  `json:"author_name" gorm:"author_name"`         // 创作者
	AuthorAvatar            string  `json:"author_avatar" gorm:"author_avatar"`     // 创作者头像
	StartSaleTime           string  `json:"start_sale_time" gorm:"start_sale_time"` // 开售时间
	DiscountPrice           float64 `json:"discount_price" gorm:"discount_price"`   // 折扣价
	OriginalPrice           float64 `json:"original_price" gorm:"original_price"`   // 划线价
	MarketType              string  `json:"market_type" gorm:"market_type"`
	BoughtVideo             string  `json:"bought_video" gorm:"bought_video"` // pfp购买后视频
	OnSaleTime              int64   `json:"on_sale_time" gorm:"on_sale_time"`
	NotSale                 int64   `json:"not_sale" gorm:"not_sale"`
	SaleMethod              string  `json:"sale_method" gorm:"sale_method"` // 发售形式: active活动发售/combine合成发售
	IsTest                  int64   `json:"is_test" gorm:"is_test"`
	OperateSellOut          int64   `json:"operate_sell_out" gorm:"operate_sell_out"`
	ThingId                 int64   `json:"thing_id" gorm:"thing_id"`                                     // 版权实物の商品id
	IsCanRetrieve           int64   `json:"is_can_retrieve" gorm:"is_can_retrieve"`                       // 是否可以取回
	FreightDiscount         int64   `json:"freight_discount" gorm:"freight_discount"`                     // 邮费折扣（非卖品使用）
	CopyrightThingPic       string  `json:"copyright_thing_pic" gorm:"copyright_thing_pic"`               // 版权实物封面图
	IsOneCanOpen            int64   `json:"is_one_can_open" gorm:"is_one_can_open"`                       // 一级发售/空投后是否可开启: 1可以/0不可以
	Point                   int64   `json:"point" gorm:"point"`                                           // 最低积分
	IsFar                   int64   `json:"is_far" gorm:"is_far"`                                         // 是否偏远地区:0不是/1是
	NftCoverPic             string  `json:"nft_cover_pic" gorm:"nft_cover_pic"`                           // 藏品封面图（只用来展示盲盒内容）
	Scene                   int8    `json:"scene" gorm:"scene"`                                           // 0默认1土地
	CategoryId              int64   `json:"category_id" gorm:"category_id"`                               // 分类id
	BlindboxOpenAmount      int64   `json:"blindbox_open_amount" gorm:"blindbox_open_amount"`             // 盲盒开启费用
	IsOpenNfts              int8    `json:"is_open_nfts" gorm:"is_open_nfts"`                             // 是否支持多开
	OpenNftNum              int8    `json:"open_nft_num" gorm:"open_nft_num"`                             // 多开数量
	DisplayOption           int8    `json:"display_option" gorm:"display_option"`                         // nft详情页:1=剩余份数 2=流通份数 3=剩余份数，流通分数都展示
	IsDisplayArtistList     int8    `json:"is_display_artist_list" gorm:"is_display_artist_list"`         // 是否在艺术家列表里面展示
	IsOpenBlindboxAnimation int8    `json:"is_open_blindbox_animation" gorm:"is_open_blindbox_animation"` // 是否开启盲盒动画
	BlindboxAnimationGif    string  `json:"blindbox_animation_gif" gorm:"blindbox_animation_gif"`         // 盲盒gif动画
	NftProductId            int64   `json:"nft_product_id" gorm:"nft_product_id"`                         // 藏品ID
	NftProductSizeId        int64   `json:"nft_product_size_id" gorm:"nft_product_size_id"`               // NFT商品款id
	IsSelfPickUp            int8    `json:"is_self_pick_up" gorm:"is_self_pick_up"`                       // 是否自提
	NftProductName          string  `json:"nft_product_name" gorm:"nft_product_name"`                     // 藏品名称
	PartnerId               int64   `json:"partner_id" gorm:"partner_id"`                                 // 业务伙伴id
	IsSelfdomBox            int8    `json:"is_selfdom_box" gorm:"is_selfdom_box"`                         // 是否是个性盲盒
	QuestionNum             int8    `json:"question_num" gorm:"question_num"`                             // 答题数量
	RegularTs               int64   `json:"regular_ts" gorm:"regular_ts"`                                 // 定时时间戳
	RegularStatus           int64   `json:"regular_status" gorm:"regular_status"`                         // 定时状态
	PwdTradeSwitch          int8    `json:"pwd_trade_switch" gorm:"pwd_trade_switch"`                     // 口令交易开关 0=关闭 1=打开
	RunRestStockStatus      int8    `json:"run_rest_stock_status" gorm:"run_rest_stock_status"`           // 跑库存状态 1可以跑库存 2正在跑库存 3已跑满库存
}

func (s *SaleCalendarProductModel) TableName() string {
	return "sale_calendar_product"
}

type SaleCalendarProduct struct {
	Ctx *gin.Context
}

type ListByConditionReq struct {
	Page         int
	PageSize     int
	IsBusiness   int32
	Fields       string
	MarketType   []string
	ProductId    uint64
	ProductName  string
	OnSaleStatus uint32
	IsTest       uint32
	ProductIds   []uint64
}

func (s SaleCalendarProduct) ListByCondition(req ListByConditionReq) (res []*SaleCalendarProductModel, err error) {
	if req.Fields == "" {
		req.Fields = "*"
	}
	query := cli.HotDogGormDB.WithContext(s.Ctx).Model(&SaleCalendarProductModel{}).
		Select(req.Fields)
	if req.IsBusiness != -1 {
		query.Where("is_business = ?", req.IsBusiness)
	}
	if len(req.MarketType) > 0 {
		query.Where("market_type IN ?", req.MarketType)
	}
	if req.ProductId != 0 {
		query.Where("id = ?", req.ProductId)
	} else if len(req.ProductIds) > 0 {
		query.Where("id in ?", req.ProductIds)
	}
	if req.ProductName != "" {
		query.Where("product_name LIKE ?", "%"+req.ProductName+"%")
	}
	if req.OnSaleStatus > 0 {
		query.Where("on_sale_status = ?", req.OnSaleStatus)
	}
	if req.IsTest > 0 {
		query.Where("is_test = ?", req.IsTest)
	}
	if len(req.ProductIds) == 0 {
		query.Order("click_count desc,weight desc")
	} else {
		query.Order("weight desc")
	}
	err = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).
		Find(&res).Error
	if err != nil {
		klog.Errorf("ecommerce.product.ListByCondition, error: %v", err)
		return nil, err
	}

	return res, nil
}

func (s SaleCalendarProduct) GetByParams(where map[string]any) (list []SaleCalendarProductModel, err error) {
	query := cli.HotDogGormDB.WithContext(s.Ctx).Model(&SaleCalendarProductModel{})
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (s SaleCalendarProduct) One(id int64) (data SaleCalendarProductModel, err error) {
	err = cli.HotDogGormDB.WithContext(s.Ctx).Model(&SaleCalendarProductModel{}).Where("id", id).First(&data).Error
	return
}

func (s SaleCalendarProduct) UpdateByParams(where, payload map[string]any) (err error) {
	query := cli.HotDogGormDB.WithContext(s.Ctx).Model(&SaleCalendarProductModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	return query.Updates(payload).Error
}

func (s SaleCalendarProduct) CountNftNum() (result []dto.ArtistNftNum, err error) {
	query := cli.HotDogGormDB.WithContext(s.Ctx).Model(&SaleCalendarProductModel{})
	query.Select(`ai_match_product_nft_artist.artist_id,
	ai_match_product_nft_artist.artist_name,
	count(1) nft_num`)
	query.Joins("JOIN ai_match_product_nft_artist ON sale_calendar_product.brand_id = ai_match_product_nft_artist.artist_id")
	query.Where("ai_match_product_nft_artist.is_delete", 0)
	query.Where("sale_calendar_product.is_delete", 0)
	err = query.Group("sale_calendar_product.brand_id").Scan(&result).Error
	return
}

type SaleCalendarProductWithArtist struct {
	SaleCalendarProductModel
	ArtistId   int64  `json:"artist_id" gorm:"artist_id"`
	ArtistName string `json:"artist_name" gorm:"artist_name"`
}

func (s SaleCalendarProduct) GetListJoinWithNftArtist(where map[string]any) (data []SaleCalendarProductWithArtist, err error) {
	fields := []string{
		"sale_calendar_product.*",
		"ai_match_product_nft_artist.artist_name", "ai_match_product_nft_artist.artist_id",
	}
	query := cli.HotDogGormDB.WithContext(s.Ctx).Table("sale_calendar_product").Select(fields).
		Joins("JOIN ai_match_product_nft_artist ON sale_calendar_product.brand_id = ai_match_product_nft_artist.artist_id")
	for k, v := range where {
		query = query.Where(k, v)
	}
	err = query.Scan(&data).Error
	return
}
