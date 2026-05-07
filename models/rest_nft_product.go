package models

import (
	"hotbox-adm-backend/cli"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type AiMatchRestNftProductModel struct {
	ID                int64          `gorm:"column:id" json:"id"`                                   //  主键id
	ProductId         int64          `gorm:"column:product_id" json:"product_id"`                   //  商品id
	ProductTitle      string         `gorm:"column:product_title" json:"product_title"`             //  商品名称
	NftProductSizeId  int64          `gorm:"column:nft_product_size_id" json:"nft_product_size_id"` //  nft商品款id
	ActivePicture     string         `gorm:"column:active_picture" json:"active_picture"`           //  活动封面图
	ReceiverProvince  string         `gorm:"column:receiver_province" json:"receiver_province"`     //  nft六位随机编号
	ReceiverCity      string         `gorm:"column:receiver_city" json:"receiver_city"`             //  nft编号
	ReceiverRegion    string         `gorm:"column:receiver_region" json:"receiver_region"`         //  nft库存总数
	CouponIds         datatypes.JSON `gorm:"column:coupon_ids" json:"coupon_ids"`                   //  优惠券ids
	IosSourceFile     string         `gorm:"column:ios_source_file" json:"ios_source_file"`
	AndroidSourceFile string         `gorm:"column:android_source_file" json:"android_source_file"`
	IsRelease         int64          `gorm:"column:is_release" json:"is_release"`     //  是否发放
	ReleaseTime       TimeWrapper    `gorm:"column:release_time" json:"release_time"` //  发放时间
	ToUser            int64          `gorm:"column:to_user" json:"to_user"`           //  发放用户
	OperatorId        int64          `gorm:"column:operator_id" json:"operator_id"`   //  操作人
	SizeId            int64          `gorm:"column:size_id" json:"size_id"`           //  尺码id
	CreateTime        TimeWrapper    `gorm:"column:create_time" json:"create_time"`
	UpdateTime        TimeWrapper    `gorm:"column:update_time" json:"update_time"`
	IsDelete          int64          `gorm:"column:is_delete" json:"is_delete"`
	BoxContent        string         `gorm:"column:box_content" json:"box_content"`
	CombineInfo       string         `gorm:"column:combine_info" json:"combine_info"`
	CombineActiveId   int            `gorm:"column:combine_active_id" json:"combine_active_id"`
	Note              string         `gorm:"column:note" json:"note"`
	PlayMethod        string         `gorm:"column:play_method" json:"play_method"` //  藏品玩法 normal普通 combin合成 upgrade升级
}

func (AiMatchRestNftProductModel) TableName() string {
	return "ai_match_rest_nft_product"
}

type AiMatchRestNftProduct struct {
	Ctx *gin.Context
}

func (n AiMatchRestNftProduct) GetRestNftProductCount(where map[string]any) (count int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchRestNftProductModel{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	err = query.Count(&count).Error
	return
}

func (n AiMatchRestNftProduct) GetRestNftProductWithParams(where map[string]any) (result []AiMatchRestNftProductModel, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchRestNftProductModel{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	err = query.Scan(&result).Error
	return
}

func (n AiMatchRestNftProduct) AddRestNftProduct(m AiMatchRestNftProductModel) error {
	return cli.HotDogGormDB.WithContext(n.Ctx).Create(&m).Error
}

func (n AiMatchRestNftProduct) GetRestNftProductList(where map[string]any, order []string, pageNum, pageSize int) (list []*AiMatchRestNftProductModel, count int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchRestNftProductModel{})
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		return
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return
}

func (n AiMatchRestNftProduct) RestNftProductStock(productId, productSizId uint64) (num int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchRestNftProductModel{}).
		Where("is_delete = 0").
		Where("is_release = 0")

	if productId > 0 {
		query.Where("product_id = ?", productId)
	}
	if productSizId > 0 {
		query.Where("nft_product_size_id = ?", productSizId)
	}
	err = query.Count(&num).Error
	if err != nil {
		return 0, err
	}

	return num, nil
}

type RestNftProductAvailableStockReq struct {
	ProductIds    []uint64
	ProductSizIds []uint64
}

type RestNftProductAvailableStockResp struct {
	ProductId      uint64 `gorm:"column:product_id"`
	ProductSizId   uint64 `gorm:"column:nft_product_size_id"`
	AvailableStock uint32 `gorm:"column:available_stock"`
}

func (n AiMatchRestNftProduct) RestNftProductAvailableStock(req RestNftProductAvailableStockReq) (res []*RestNftProductAvailableStockResp, err error) {
	// 查询可用库存
	res = make([]*RestNftProductAvailableStockResp, 0)
	// HOTDOG PASS(1019327, 1321) 数量太大，不查询库存

	queryProductIds := make([]uint64, 0)
	for _, pid := range req.ProductIds {
		if pid != 1019327 {
			queryProductIds = append(queryProductIds, pid)
		}
	}
	querySizeIds := make([]uint64, 0)
	for _, sid := range req.ProductSizIds {
		if sid != 1321 {
			querySizeIds = append(querySizeIds, sid)
		}
	}
	// 如果只查询pass卡返回, 空数组
	if len(req.ProductIds) == 1 && req.ProductIds[0] == 1019327 {
		return res, nil
	}
	if len(req.ProductSizIds) == 1 && req.ProductSizIds[0] == 1321 {
		return res, nil
	}
	query := cli.HotDogGormDB.WithContext(n.Ctx).Model(&AiMatchRestNftProductModel{}).
		Select(`product_id, nft_product_size_id, count(*) as available_stock`).
		Where("is_delete = 0").
		Where("is_release = 0")

	if len(queryProductIds) > 0 {
		query.Where("product_id IN ?", queryProductIds)
	}

	if len(querySizeIds) > 0 {
		query.Where("nft_product_size_id IN ?", querySizeIds)
	}
	err = query.Group(`product_id, nft_product_size_id`).Find(&res).Error
	if err != nil {
		klog.Errorf("ecommerce.ProductNft.RestNftProductAvailableStock, error: %v", err)
		return nil, err
	}

	return res, nil
}
