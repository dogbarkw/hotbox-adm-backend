package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SaleProductNftSizeModel struct {
	ID                   int64     `gorm:"column:id" json:"id"`                                     //  主键id
	ProductId            int64     `gorm:"column:product_id" json:"product_id"`                     //  商品id
	NftProductSizeId     int64     `gorm:"column:nft_product_size_id" json:"nft_product_size_id"`   //  藏品款式id
	ProductTitle         string    `gorm:"column:product_title" json:"product_title"`               //  商品名称
	ProductAbbreviation  int64     `gorm:"column:product_abbreviation" json:"product_abbreviation"` //  产品编号
	StockCount           int64     `gorm:"column:stock_count" json:"stock_count"`                   //  数量
	TotalCount           int64     `gorm:"column:total_count" json:"total_count"`                   //  总数量
	AfterSalePicture     string    `gorm:"column:after_sale_picture" json:"after_sale_picture"`     //  购买后图片
	IosSourceFile        string    `gorm:"column:ios_source_file" json:"ios_source_file"`           //  ios源文件
	AndroidSourceFile    string    `gorm:"column:android_source_file" json:"android_source_file"`   //  android源文件
	SourceType           string    `gorm:"column:source_type" json:"source_type"`                   //  源文件类型
	Weight               int64     `gorm:"column:weight" json:"weight"`                             //  权重
	CreateTime           time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime           time.Time `gorm:"column:update_time" json:"update_time"`
	IsDelete             int64     `gorm:"column:is_delete" json:"is_delete"`
	DisplayPicture       string    `gorm:"column:display_picture" json:"display_picture"`
	SaleMinPrice         float64   `gorm:"column:sale_min_price" json:"sale_min_price"` //  最低售价
	SaleMaxPrice         float64   `gorm:"column:sale_max_price" json:"sale_max_price"` //  最高售价
	SecondPicture        string    `gorm:"column:second_picture" json:"second_picture"` //  二级市场展示图
	MtOrderSn            string    `gorm:"column:mt_order_sn" json:"mt_order_sn"`
	DdcId                int64     `gorm:"column:ddc_id" json:"ddc_id"`
	TransactionHash      string    `gorm:"column:transaction_hash" json:"transaction_hash"`
	Version              int64     `gorm:"column:version" json:"version"`
	QrCodePic            string    `gorm:"column:qr_code_pic" json:"qr_code_pic"`   //  ar图片扫码
	CouponId             int64     `gorm:"column:coupon_id" json:"coupon_id"`       //  优惠券id
	DetailHref           string    `gorm:"column:detail_href" json:"detail_href"`   //  详情页链接
	DisplayHref          string    `gorm:"column:display_href" json:"display_href"` //  展示页链接
	OwnerAddress         string    `gorm:"column:owner_address" json:"owner_address"`
	OwnerPrivateKey      string    `gorm:"column:owner_private_key" json:"owner_private_key"`
	Hot                  int64     `gorm:"column:hot" json:"hot"`
	HotdogTokenid        int64     `gorm:"column:hotdog_tokenid" json:"hotdog_tokenid"` //  hotdog_chain_tokenid
	TrxHash              string    `gorm:"column:trx_hash" json:"trx_hash"`             //  hotdog_chain_交易hash
	HotdogAddr           string    `gorm:"column:hotdog_addr" json:"hotdog_addr"`       //  hotdog_chain_地址
	HotdogKey            string    `gorm:"column:hotdog_key" json:"hotdog_key"`         //  hotdog_chain_私钥
	Grade                string    `gorm:"column:grade" json:"grade"`
	Score                int64     `gorm:"column:score" json:"score"`
	CycleDays            int64     `gorm:"column:cycle_days" json:"cycle_days"`                           //  持有天数
	ScoreRate            int64     `gorm:"column:score_rate" json:"score_rate"`                           //  消费发积分比例；1就是1% -> 0.01
	LabelType            string    `gorm:"column:label_type" json:"label_type"`                           //  标签类别
	ReleaseRate          float64   `gorm:"column:release_rate" json:"release_rate"`                       //  按比例分配盲盒 的 比例
	ArVersion            int64     `gorm:"column:ar_version" json:"ar_version"`                           //  ar模型走哪个版本 - 1url版本/2unity版本
	ArType               string    `gorm:"column:ar_type" json:"ar_type"`                                 //  ar识别图类型 - plane平面/photo图片
	UnityIosFile         string    `gorm:"column:unity_ios_file" json:"unity_ios_file"`                   //  unity_ios源文件
	UnityAndroidFile     string    `gorm:"column:unity_android_file" json:"unity_android_file"`           //  unity_ios源文件
	ArDisplayVideo       string    `gorm:"column:ar_display_video" json:"ar_display_video"`               //  ar展示视频
	ArDisplayPicture     string    `gorm:"column:ar_display_picture" json:"ar_display_picture"`           //  ar展示图片
	AppearForm           string    `gorm:"column:appear_form" json:"appear_form"`                         //  出现形式
	NftCoverPic          string    `gorm:"column:nft_cover_pic" json:"nft_cover_pic"`                     //  藏品封面图
	DetailHeadDisplayPic string    `gorm:"column:detail_head_display_pic" json:"detail_head_display_pic"` //  详情页头部展示图
	MarketType           string    `gorm:"column:market_type" json:"market_type"`                         //  市场类型
	ConsignStatus        int64     `gorm:"column:consign_status" json:"consign_status"`                   //  挂售状态 1=寄售中
	ConsignOrderId       int64     `gorm:"column:consign_order_id" json:"consign_order_id"`               //  挂售订单id
	CategoryId           int64     `gorm:"column:category_id" json:"category_id"`                         //  分类id
	BoxId                int64     `gorm:"column:box_id" json:"box_id"`                                   //  盲盒商品id(主盲盒id；nft_product_size_id=0时使用)
	PwdTradeSwitch       int64     `gorm:"column:pwd_trade_switch" json:"pwd_trade_switch"`               //  口令交易开关 0=关闭 1=打开
}

func (SaleProductNftSizeModel) TableName() string {
	return "sale_product_nft_size"
}

type SaleProductNftSize struct{}

func NewSaleProductNftSize() *SaleProductNftSize {
	return &SaleProductNftSize{}
}

func (s *SaleProductNftSize) GetOneByParams(ctx context.Context, where map[string]any) (r SaleProductNftSizeModel, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Where(where).First(&r).Error
	return
}

type SaleProductNftSizeListByConditionReq struct {
	Page          int
	PageSize      int
	Id            uint64
	Ids           []uint64
	ProductId     uint64
	ProductIds    []uint64
	Fields        string
	ProductTitle  string
	Status        uint32
	MarketType    string
	SearchTitle   string
	ConsignStatus uint8
	Orderby       uint8
}

func (s *SaleProductNftSize) NftProductListByCondition(ctx context.Context, req SaleProductNftSizeListByConditionReq) (res []*SaleProductNftSizeModel, err error) {
	fields := req.Fields
	if fields == "" {
		fields = "*"
	}
	query := cli.HotDogGormDB.WithContext(ctx).Model(&SaleProductNftSizeModel{}).Select(fields)
	if req.Id > 0 {
		query.Where("nft_product_size_id = ?", req.Id)
	}
	if len(req.Ids) > 0 {
		query.Where("nft_product_size_id IN ?", req.Ids)
	}
	if req.ProductId > 0 {
		query.Where("product_id = ?", req.ProductId)
	}
	if len(req.ProductIds) > 0 {
		query.Where("product_id IN ?", req.ProductIds)
	}
	switch req.Status { // 0=所有 1=正常
	case 1:
		query.Where("is_delete = 0")
	}
	if req.ProductTitle != "" {
		query.Where("product_title LIKE ?", "%"+req.ProductTitle+"%")
	}
	if len(req.SearchTitle) > 0 {
		query.Where("product_title = ?", req.SearchTitle)
	}
	if req.Page > 0 && req.PageSize > 0 {
		query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	}
	if req.MarketType != "" {
		query.Where("market_type = ?", req.MarketType)
	}
	if req.ConsignStatus != 0 {
		query.Where("consign_status = 1")
	}
	switch req.Orderby {
	case 1: // 价格升序
		query.Order("sale_min_price ASC")
	case 2: // 价格降序
		query.Order("sale_min_price DESC")
	}
	err = query.Find(&res).Error
	if err != nil {
		klog.Errorf("ecommerce.ProductNftSize.NftProductListByCondition error: %v", err)
		return nil, err
	}

	return res, nil
}

func (a *SaleProductNftSize) GetSaleProductNftSizeByParams(ctx context.Context, where map[string]any, order []string, limit *int) (list []SaleProductNftSizeModel, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&SaleProductNftSizeModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}
