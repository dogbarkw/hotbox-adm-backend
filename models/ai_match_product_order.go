package models

import (
	"context"
	"time"

	"hotbox-adm-backend/dto"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

var AiMatchProductOrderDal = &AiMatchProductOrder{}

type AiMatchProductOrderModel struct {
	ID                    int64          `gorm:"column:id" json:"id"`                                           //  订单id
	OrderSn               *string        `gorm:"column:order_sn" json:"order_sn"`                               //  订单编号
	ProductId             int64          `gorm:"column:product_id" json:"product_id"`                           //  商品id
	ProductTitle          string         `gorm:"column:product_title" json:"product_title"`                     //  商品名称
	ProductPicture        string         `gorm:"column:product_picture" json:"product_picture"`                 //  商品封面图
	SizeId                int64          `gorm:"column:size_id" json:"size_id"`                                 //  尺码id
	Size                  string         `gorm:"column:size" json:"size"`                                       //  尺寸信息
	UserId                int64          `gorm:"column:user_id" json:"user_id"`                                 //  用户id
	TotalAmount           *float64       `gorm:"column:total_amount" json:"total_amount"`                       //  订单总金额
	PayAmount             float64        `gorm:"column:pay_amount" json:"pay_amount"`                           //  应付金额（实际支付金额）
	FreightAmount         *float64       `gorm:"column:freight_amount" json:"freight_amount"`                   //  运费金额
	Status                int64          `gorm:"column:status" json:"status"`                                   //  订单状态：待付款；待发货；已发货；已完成；已关闭；无效订单
	PayType               int64          `gorm:"column:pay_type" json:"pay_type"`                               //  支付方式：0->未支付；1->支付宝；2->微信
	DeliveryStatus        int64          `gorm:"column:delivery_status" json:"delivery_status"`                 //  发货状态：未发货  已发货
	ConfirmStatus         *int           `gorm:"column:confirm_status" json:"confirm_status"`                   //  确认收货状态：0->未确认；1->已确认
	OutTradeNo            *string        `gorm:"column:out_trade_no" json:"out_trade_no"`                       //  支付单号
	DeliverySn            string         `gorm:"column:delivery_sn" json:"delivery_sn"`                         //  物流单号
	ReceiverName          string         `gorm:"column:receiver_name" json:"receiver_name"`                     //  收货人姓名
	ReceiverPhone         *string        `gorm:"column:receiver_phone" json:"receiver_phone"`                   //  收货人电话
	ReceiverProvince      string         `gorm:"column:receiver_province" json:"receiver_province"`             //  省份/直辖市
	ReceiverCity          string         `gorm:"column:receiver_city" json:"receiver_city"`                     //  城市
	ReceiverRegion        string         `gorm:"column:receiver_region" json:"receiver_region"`                 //  区
	ReceiverDetailAddress string         `gorm:"column:receiver_detail_address" json:"receiver_detail_address"` //  详细地址
	Note                  string         `gorm:"column:note" json:"note"`                                       //  备注
	PaymentTime           *time.Time     `gorm:"column:payment_time" json:"payment_time"`                       //  支付时间
	DeliveryTime          *time.Time     `gorm:"column:delivery_time" json:"delivery_time"`                     //  发货时间
	ArrivalTime           *time.Time     `gorm:"column:arrival_time" json:"arrival_time"`                       //  确认收货时间
	ReceiveTime           time.Time      `gorm:"column:receive_time" json:"receive_time"`                       //  确认收货时间
	CreateTime            time.Time      `gorm:"column:create_time" json:"create_time"`                         //  提交时间
	UpdateTime            time.Time      `gorm:"column:update_time" json:"update_time"`                         //  修改时间
	IsDelete              int64          `gorm:"column:is_delete" json:"is_delete"`                             //  是否删除
	OrderNo               string         `gorm:"column:order_no" json:"order_no"`                               //  订单编号
	Version               int64          `gorm:"column:version" json:"version"`
	BuyCount              int64          `gorm:"column:buy_count" json:"buy_count"`
	EstimatedDeliveryTime *string        `gorm:"column:estimated_delivery_time" json:"estimated_delivery_time"`
	ProductType           string         `gorm:"column:product_type" json:"product_type"`
	LogisticsInfo         string         `gorm:"column:logistics_info" json:"logistics_info"`
	Supplier              *string        `gorm:"column:supplier" json:"supplier"` //  供应商
	SkuId                 string         `gorm:"column:sku_id" json:"sku_id"`
	AfterSales            string         `gorm:"column:after_sales" json:"after_sales"` //  售后内容
	ApplyRefundStatus     int64          `gorm:"column:apply_refund_status" json:"apply_refund_status"`
	CouponId              *int           `gorm:"column:coupon_id" json:"coupon_id"`                     //  优惠券id
	CouponAmount          *float64       `gorm:"column:coupon_amount" json:"coupon_amount"`             //  优惠券抵扣金额
	Invoice               datatypes.JSON `gorm:"column:invoice" json:"invoice"`                         //  开票内容
	InvoiceTime           *time.Time     `gorm:"column:invoice_time" json:"invoice_time"`               //  申请开票时间
	UserCouponId          *int           `gorm:"column:user_coupon_id" json:"user_coupon_id"`           //  用户优惠券id
	NewFlashId            int            `gorm:"column:new_flash_id" json:"new_flash_id"`               //  快讯动态id
	AndroidNftSource      string         `gorm:"column:android_nft_source" json:"android_nft_source"`   //  安卓nft源文件
	NftProductSizeId      int64          `gorm:"column:nft_product_size_id" json:"nft_product_size_id"` //  nft商品款id
	SecondId              int64          `gorm:"column:second_id" json:"second_id"`
	MiddleUserId          int64          `gorm:"column:middle_user_id" json:"middle_user_id"`
	PlayMethod            string         `gorm:"column:play_method" json:"play_method"` //  藏品玩法
	PayTypeThird          int64          `gorm:"column:pay_type_third" json:"pay_type_third"`
	BoxContent            string         `gorm:"column:box_content" json:"box_content"`           //  盲盒拆前内容
	RedeemCouponId        int64          `gorm:"column:redeem_coupon_id" json:"redeem_coupon_id"` //  待兑换优惠券
	BoxCode               string         `gorm:"column:box_code" json:"box_code"`                 //  盲盒拆前编号
	PropUserUuid          string         `gorm:"column:prop_user_uuid" json:"prop_user_uuid"`
	PresaleId             int64          `gorm:"column:presale_id" json:"presale_id"` //  预售活动id
	SourceType            string         `gorm:"column:source_type" json:"source_type"`
	ExaUserId             int64          `gorm:"column:exa_user_id" json:"exa_user_id"`
	ExaMiddleUserId       int64          `gorm:"column:exa_middle_user_id" json:"exa_middle_user_id"`
	MarketType            string         `gorm:"column:market_type" json:"market_type"` //  市场类型
	Point                 int64          `gorm:"column:point" json:"point"`             //  消耗积分数量
	IsChild               int64          `gorm:"column:is_child" json:"is_child"`       //  子订单 0:否 1:是
	Scene                 int64          `gorm:"column:scene" json:"scene"`             //  场景: 0默认常规/1土地
	IsRent                int64          `gorm:"column:is_rent" json:"is_rent"`         //  出租状态 1=出租中
	// PrevAgentOrderId      int64          `gorm:"column:prev_agent_order_id" json:"prev_agent_order_id"` //  上级代理订单id
	// AgentStartId int64  `gorm:"column:agent_start_id" json:"agent_start_id"` //  一级代理订单id
	// GroupUuid    string `gorm:"column:group_uuid" json:"group_uuid"`         //  批量订单组id
}

func (AiMatchProductOrderModel) TableName() string {
	return "ai_match_product_order"
}

type AiMatchProductOrder struct {
	Ctx *gin.Context
}

func (n AiMatchProductOrder) GetProductOrderList(where map[string]any, count *int) (list []AiMatchProductOrderModel, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{}).Where(where)
	if count != nil {
		query.Limit(*count)
	}
	err = query.Scan(&list).Error
	return
}

func (n AiMatchProductOrder) GetRestNftProductCount(where map[string]any) (count int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	err = query.Count(&count).Error
	return
}

func (n AiMatchProductOrder) GetProductOrderCountByParams(where map[string]any, notIn map[string]any) (count int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	for k, v := range notIn {
		query.Not(k, v)
	}
	err = query.Count(&count).Error
	return
}

func (n AiMatchProductOrder) GetProductOrderCountByGroupParams(where map[string]any, fields, group string) (result []OpCount, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{}).Select(fields)
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Group(group).Scan(&result).Error
	return
}

func (n AiMatchProductOrder) GetUserCountInProductOrderGroupByReceiverName(where map[string]any, fields, group string) (result []int, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{}).Select(fields)
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Group(group).Scan(&result).Error
	return
}

func (n AiMatchProductOrder) UpdateProductOrderByParams(where, data map[string]any) (err error) {
	return cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{}).Where(where).Updates(data).Error
}

func (n AiMatchProductOrder) GetDailyGMV(ctx context.Context, startTime string, endTime string) (res []dto.GmvData, err error) {
	sql := `
WITH t1 AS (
	SELECT
		sum( pay_amount ) AS gmv,
		DATE ( payment_time ) AS dt,
		count( DISTINCT user_id ) AS user_cnt 
	FROM
		ai_match_product_order 
	WHERE
		is_delete = 0 
		AND STATUS IN ( 2, 15, 76, 86, 99, 36 ) 
		AND product_type = 'NFT_SE' 
		AND product_title not like '%测试%'
		AND middle_user_id > 0 
		AND payment_time >= ? 
		AND payment_time <= ?
	GROUP BY
		DATE ( payment_time ) 
	ORDER BY
		DATE ( payment_time ) DESC 
	),
	t2 AS (
	SELECT
		sum( pay_amount ) AS r_gmv,
		DATE ( payment_time ) AS dt,
		count( DISTINCT user_id ) AS r_user_cnt 
	FROM
		ai_match_product_order 
	WHERE
		is_delete = 0 
		AND STATUS IN ( 2, 15, 76, 86, 99, 36 ) 
		AND product_type = 'NFT_SE' 
		AND product_title not like '%测试%'
		AND middle_user_id > 0 
		AND user_id NOT IN ( SELECT user_id FROM test_user )
		AND payment_time >= ? 
		AND payment_time <= ?
	GROUP BY
		DATE ( payment_time ) 
	ORDER BY
		DATE ( payment_time ) DESC 
	),
	res AS(
 SELECT t1.dt,t1.gmv,t1.user_cnt,IFNULL(t2.r_gmv,0) r_gmv,IFNULL(t2.r_user_cnt,0) r_user_cnt FROM t1 LEFT JOIN t2 ON t1.dt=t2.dt
)SELECT * FROM res
`
	err = cli.HotDogGormDB.WithContext(ctx).Raw(sql, startTime, endTime, startTime, endTime).Find(&res).Error
	return
}

// GetTestUserOrderLogs 统计用户进账额度
func (n AiMatchProductOrder) GetTestUserOrderLogs(userId int64, startTime, endTime time.Time, fields []string) (result []AiMatchProductOrderModel, err error) {
	err = cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchProductOrderModel{}).
		Where("status in (2, 86, 99, 15, 76, 67, 36) and product_type = 'NFT_SE' and is_delete = 0 and middle_user_id = ?", userId).
		Where("payment_time >= ?", startTime).
		Where("payment_time < ?", endTime).
		Where("receiver_name  NOT IN (?)", cli.SpecialUserIds).
		Select(fields).
		Scan(&result).Error
	return
}
