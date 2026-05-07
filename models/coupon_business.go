package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var CouponBusinessDal = &CouponBusiness{}

// CouponBusiness undefined
type CouponBusiness struct {
	ID              int64     `json:"id" gorm:"id"` // 主键id
	Name            string    `json:"name" gorm:"name"`
	Desc            string    `json:"desc" gorm:"desc"` // 优惠券描述
	Type            string    `json:"type" gorm:"type"` // 优惠券类型
	Amount          float64   `json:"amount" gorm:"amount"`
	UseRestriction  int64     `json:"use_restriction" gorm:"use_restriction"`   // 使用门槛：0无门槛；1满多少可用
	AvailableAmount float64   `json:"available_amount" gorm:"available_amount"` // 满减金额
	Quantity        int64     `json:"quantity" gorm:"quantity"`
	ReleaseQuantity int64     `json:"release_quantity" gorm:"release_quantity"`
	ExpiryDay       int64     `json:"expiry_day" gorm:"expiry_day"`
	AvailableDay    int64     `json:"available_day" gorm:"available_day"` // 可用时长(天)
	BeginTime       time.Time `json:"begin_time" gorm:"begin_time"`
	EndTime         time.Time `json:"end_time" gorm:"end_time"`
	GetLimit        int64     `json:"get_limit" gorm:"get_limit"`           // 领取限制：0限领；1无限制
	LimitNum        int64     `json:"limit_num" gorm:"limit_num"`           // 限领几张
	UseLimit        string    `json:"use_limit" gorm:"use_limit"`           // 使用限制：all所有商品可用；category品类券；brand品牌券；specify指定商品
	ReleaseMethod   string    `json:"release_method" gorm:"release_method"` // 发放方式：new新人券；operator运营发放；command口令
	Operator        int64     `json:"operator" gorm:"operator"`             // 操作人
	Online          int64     `json:"online" gorm:"online"`                 // 是否启用优惠券: 0:禁用, 1:启用
	CreateTime      time.Time `json:"create_time" gorm:"create_time"`
	UpdateTime      time.Time `json:"update_time" gorm:"update_time"`
	IsDelete        int64     `json:"is_delete" gorm:"is_delete"`
	UsedQuantity    int64     `json:"used_quantity" gorm:"used_quantity"` // 累计使用数量
	Discount        float64   `json:"discount" gorm:"discount"`           // 折扣
	TypeCode        string    `json:"type_code" gorm:"type_code"`         // 优惠券code
	Createtime      time.Time `json:"CreateTime" gorm:"CreateTime"`
	Version         int64     `json:"version" gorm:"version"`
	CouponPicture   string    `json:"coupon_picture" gorm:"coupon_picture"`
	Remark          string    `json:"remark" gorm:"remark"`
}

// TableName 表名称
func (*CouponBusiness) TableName() string {
	return "coupon_business"
}

func (a *CouponBusiness) GetByParams(ctx context.Context, where map[string]any) (list []CouponBusiness, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&CouponBusiness{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	query = query.Scan(&list)
	err = query.Error
	return
}
