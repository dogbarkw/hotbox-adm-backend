package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/pkg/constant"

	"github.com/samber/lo"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type ProductNftSize struct {
	Id                   uint64  `gorm:"column:id" json:"id"`
	ProductId            uint64  `gorm:"column:product_id" json:"product_id"`
	ProductTitle         string  `gorm:"column:product_title" json:"product_title"`
	Weight               uint32  `gorm:"column:weight" json:"weight"`
	LabelType            string  `gorm:"column:label_type" json:"label_type"`
	DetailHref           string  `gorm:"column:detail_href" json:"detail_href"`
	DisplayHref          string  `gorm:"column:display_href" json:"display_href"`
	IosSourceFile        string  `gorm:"column:ios_source_file" json:"ios_source_file"`
	AndroidSourceFile    string  `gorm:"column:android_source_file" json:"android_source_file"`
	IsDelete             uint32  `gorm:"column:is_delete" json:"is_delete"`
	CreateTime           string  `gorm:"column:create_time" json:"create_time"`
	UpdateTime           string  `gorm:"column:update_time" json:"update_time"`
	StockCount           uint32  `gorm:"column:stock_count" json:"stock_count"`
	TotalCount           uint32  `gorm:"column:total_count" json:"total_count"`
	SourceType           string  `gorm:"column:source_type" json:"source_type"`
	ReleaseRate          string  `gorm:"column:release_rate" json:"release_rate"`
	Score                uint32  `gorm:"column:score" json:"score"`
	ScoreRate            uint32  `gorm:"column:score_rate" json:"score_rate"`
	CycleDays            uint32  `gorm:"column:cycle_days" json:"cycle_days"`
	CouponId             uint64  `gorm:"column:coupon_id" json:"coupon_id"`
	ArVersion            uint32  `gorm:"column:ar_version" json:"ar_version"`
	ProductAbbreviation  uint32  `gorm:"column:product_abbreviation" json:"product_abbreviation"`
	ArType               string  `gorm:"column:ar_type" json:"ar_type"`
	UnityIosFile         string  `gorm:"column:unity_ios_file" json:"unity_ios_file"`
	SecondPicture        string  `gorm:"column:second_picture" json:"second_picture"`
	UnityAndroidFile     string  `gorm:"column:unity_android_file" json:"unity_android_file"`
	AfterSalePicture     string  `gorm:"column:after_sale_picture" json:"after_sale_picture"`
	ArDisplayPicture     string  `gorm:"column:ar_display_picture" json:"ar_display_picture"`
	ArDisplayVideo       string  `gorm:"column:ar_display_video" json:"ar_display_video"`
	DisplayPicture       string  `gorm:"column:display_picture" json:"display_picture"`
	QrCodePic            string  `gorm:"column:qr_code_pic" json:"qr_code_pic"`
	NftCoverPic          string  `gorm:"column:nft_cover_pic" json:"nft_cover_pic"`
	DetailHeadDisplayPic string  `gorm:"column:detail_head_display_pic" json:"detail_head_display_pic"`
	SaleMinPrice         float32 `gorm:"column:sale_min_price" json:"sale_min_price"`
	SaleMaxPrice         float32 `gorm:"column:sale_max_price" json:"sale_max_price"`
	CategoryId           uint32  `gorm:"column:category_id" json:"category_id"`
	Hot                  uint64  `gorm:"column:hot;default:1000000000" json:"hot"`
	NftProductSizeId     uint64  `gorm:"column:nft_product_size_id" json:"nft_product_size_id"`
	BoxId                uint64  `gorm:"column:box_id" json:"box_id"`
}

type Nft struct {
	Ctx *gin.Context
}

func (n Nft) NftProductPriceList(req form.SecondPriceListReq) (res []*dto.ProductNftSecondPriceJoinSaleCalendar, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(n.Ctx).Table("ai_match_product_nft_second_price").Select("ai_match_product_nft_second_price.*, sale_calendar_product.author_name, sale_calendar_product.nft_type, sale_calendar_product.is_test, sale_calendar_product.brand_id as artist_id,t.nft_count, t.remain_count as total_remain_count").
		Joins("join sale_calendar_product on ai_match_product_nft_second_price.product_id = sale_calendar_product.id").
		Joins("left join business_nft_market_warehouse_total_count as t on t.product_id = ai_match_product_nft_second_price.product_id and ai_match_product_nft_second_price.nft_product_size_id = t.product_size_id").
		Where("ai_match_product_nft_second_price.is_delete = 0 and sale_calendar_product.is_delete = 0")

	if req.ProductId > 0 {
		query.Where("ai_match_product_nft_second_price.product_id", req.ProductId)
	}
	if req.NftProductSizeId > 0 {
		query.Where("ai_match_product_nft_second_price.nft_product_size_id", req.NftProductSizeId)
	}
	if req.OnSaleStatus != -1 {
		query.Where("sale_calendar_product.on_sale_status", req.OnSaleStatus)
	}

	if req.ProductTitle != "" {
		query.Where("ai_match_product_nft_second_price.product_title like ? ", "%"+req.ProductTitle+"%")
	}

	err = query.Count(&total).Error
	if err != nil {
		klog.Errorf("NftProductPriceList query count error: %v", err)
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	if req.PageSize != 0 && req.PageNumber != 0 {
		query.Offset((int(req.PageNumber) - 1) * int(req.PageSize)).Limit(int(req.PageSize))
	}
	err = query.Order("ai_match_product_nft_second_price.create_time desc").Scan(&res).Error
	if err != nil {
		klog.Errorf("NftProductPriceList query list error: %v", err)
		return nil, 0, err
	}
	productIds := make([]uint64, 0)
	nftProductSizeIds := make([]uint64, 0)
	nftProductSizeData := []dto.SaleProductNFTSizeModel{}
	nftProductSizeDataMap := make(map[int]int)
	saleCalendarProductSizeData := []dto.SaleCalendarProductSizeModel{}
	saleCalendarProductSizeDataMap := make(map[int]int)
	for _, v := range res {
		t, err := time.Parse("2006-01-02T15:04:05-07:00", v.CountResetTime)
		if err != nil {
			v.CountResetTime = "0"
		} else {
			if t.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
				v.CountResetTime = "0"
			} else {
				v.CountResetTime = fmt.Sprintf("%d", t.UnixMilli())
			}
		}
		productIds = append(productIds, v.ProductId)
		nftProductSizeIds = append(nftProductSizeIds, v.NftProductSizeId)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("sale_product_nft_size").
			Where("nft_product_size_id IN (?)", nftProductSizeIds).
			Where("is_delete", 0).
			Scan(&nftProductSizeData).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("sale_calendar_product_size").
			Where("product_id IN (?)", productIds).
			Where("is_delete", 0).
			Scan(&saleCalendarProductSizeData).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		klog.Errorf("NftProductPriceList error group error: %v", err)
		return nil, 0, err
	}
	for _, v := range nftProductSizeData {
		nftProductSizeDataMap[v.NftProductSizeId] = v.TotalCount
	}
	for _, v := range saleCalendarProductSizeData {
		saleCalendarProductSizeDataMap[v.ProductId] = v.TotalCount
	}
	for _, v := range res {
		if v.NftProductSizeId == 0 {
			v.TotalCount = saleCalendarProductSizeDataMap[int(v.ProductId)]
		} else {
			v.TotalCount = nftProductSizeDataMap[int(v.NftProductSizeId)]
		}
	}
	return res, total, nil
}

func GetProductOrderCount(ctx context.Context, specialUids []string, result *int64, productId, nftProductSizeId int) error {
	query := cli.HotDogGormDB.WithContext(ctx).
		Table("ai_match_product_order").
		Where("nft_product_size_id = ? ", nftProductSizeId).
		Where("product_id = ? ", productId).
		Where("product_type = ? ", "NFT").
		Where("status = ? ", 2).
		Where("is_delete = ? ", 0)
	if specialUids != nil {
		query.Where("receiver_name NOT IN (?)", specialUids)
	}
	query = query.Where("note NOT IN (?) OR (note = ? AND create_time != receive_time) ", []string{"TEST_CONCURRENCY", "TEST_CONGISN"}, "TEST_CONGISN")

	err := query.Count(result).Error
	return err
}

func getTotalCount(ctx context.Context, totalCount, productId, nftProductSizeId int, combineAppType /* 合成渠道，hotdog还是游戏 */, isCombineProduct /* 是否是合成来的*/, isUpgrade /* 是否是升级来的*/, isLimitCombine int /*是否是指定编号合成*/, sp *dto.SaleProductNFTSizeModel) (c int, err error) {
	if totalCount == 0 && nftProductSizeId != 0 {
		c = 0
	} else {
		if nftProductSizeId == 0 { // 盲盒
			size := struct {
				TotalCount int
				ProductId  int
				StockCount int
			}{}
			err = cli.HotDogGormDB.WithContext(ctx).
				Table("sale_calendar_product_size").
				Where("product_id", productId).
				First(&size).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, err
			}

			if isUpgrade > 0 { // 升级
				c = size.TotalCount - size.StockCount
			} else {
				c = size.TotalCount
			}
		} else { // 某个款式
			if isUpgrade > 0 {
				c = sp.TotalCount - sp.StockCount
			} else { // 发售来的
				c = sp.TotalCount
			}
		}
	}
	return
}

func (n Nft) NftSecondPriceFlushSurplus(req form.SecondPriceFlushSurplusReq) (code int, total int, err error) {
	var realUserSurplus int64
	var restSurplusCount int64
	var remainCount int64
	switch req.FlushType {
	case 1:
		err = GetProductOrderCount(n.Ctx, cli.SpecialUserIds, &realUserSurplus, req.ProductId, req.NftProductSizeId)
		if err != nil {
			return 0, 0, err
		}
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("ai_match_product_nft_second_price").
			Where("product_id = ? ", req.ProductId).
			Where("nft_product_size_id = ? ", req.NftProductSizeId).
			Update("real_user_surplus", realUserSurplus).Error
		if err != nil {
			return 0, 0, err
		}
		return 0, cast.ToInt(realUserSurplus), nil
	case 2:
		err = GetProductOrderCount(n.Ctx, nil, &realUserSurplus, req.ProductId, req.NftProductSizeId)
		if err != nil {
			return 0, 0, err
		}
		err = cli.HotDogGormDB.WithContext(n.Ctx).
			Table("ai_match_rest_nft_product").
			Where("product_id = ? ", req.ProductId).
			Where("nft_product_size_id = ? ", req.NftProductSizeId).
			Where("is_release = ? ", 0).
			Where("is_delete = ? ", 0).
			Count(&restSurplusCount).Error
		if err != nil {
			return 0, 0, err
		}
		allUserSurplus := restSurplusCount + realUserSurplus
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("ai_match_product_nft_second_price").
			Where("product_id = ? ", req.ProductId).
			Where("nft_product_size_id = ? ", req.NftProductSizeId).
			Update("all_user_surplus", allUserSurplus).Error
		if err != nil {
			return 0, 0, err
		}
		return
	case 3:
		err = GetProductOrderCount(n.Ctx, nil, &remainCount, req.ProductId, req.NftProductSizeId)
		if err != nil {
			return 0, 0, err
		}
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("ai_match_product_nft_second_price").
			Where("product_id = ? ", req.ProductId).
			Where("nft_product_size_id = ? ", req.NftProductSizeId).
			Update("remain_count", remainCount).Error
		if err != nil {
			return 0, 0, err
		}
		return
	case 4:
		sp := dto.SaleProductNFTSizeModel{}
		err = cli.HotDogGormDB.WithContext(n.Ctx).Table("sale_product_nft_size").
			Where("nft_product_size_id", req.NftProductSizeId).
			Where("is_delete", 0).
			First(&sp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 10018, 0, nil
			}
			return 0, 0, err
		}
		mpc := dto.AiMatchProductNftCombination{}
		err = cli.HotDogGormDB.WithContext(n.Ctx).
			Table("ai_match_product_nft_combination").
			Where("product_id", req.ProductId).
			Where("is_delete", 0).
			Where("on_sale_status", 1).
			First(&mpc).Error
		isCombineProduct := 0
		combineAppType := 1
		isLimitCombine := 0
		isUpgrade := 0
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, err
		} else {
			if err == nil {
				isCombineProduct = 1
				combineAppType = mpc.AppType
				isLimitCombine = mpc.IsLimitCombine
			}
		}
		err = cli.HotDogGormDB.WithContext(n.Ctx).
			Table("activity_upgrade").
			Where("product_id", req.ProductId).
			Where("status", 2).
			First(&dto.ActivityUpgradeModel{}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, err
		} else {
			if err == nil {
				isUpgrade = 1
			}
		}

		c, err := getTotalCount(n.Ctx,
			sp.TotalCount, req.ProductId,
			req.NftProductSizeId,
			combineAppType,
			isCombineProduct,
			isUpgrade, isLimitCombine, &sp)
		if err != nil {
			return 0, 0, err
		}
		return 0, c, nil
	default:
		return
	}
}

func (n Nft) GetNftSecondPriceFlushSurplus(req form.GetSecondPriceFlushSurplusReq) (total int64, err error) {
	var realUserSurplus int64
	err = GetProductOrderCount(n.Ctx, cli.SpecialUserIds, &realUserSurplus, req.ProductId, *req.NftProductSizeId)
	if err != nil {
		return 0, err
	}
	return realUserSurplus, nil
}

type OpCount struct {
	ProductTitle string `json:"product_title"`
	ProductId    int    `json:"product_id"`
	Count        int    `json:"count"`
}

func (n Nft) NftSecondPriceFlushUserPercentage(req form.NftSecondPriceFlushUserPercentageReq) (realUserSurplus int64, nftCount int64, userPercentage float64, err error) {
	// 更新普通用户剩余份数
	err = GetProductOrderCount(n.Ctx, cli.SpecialUserIds, &realUserSurplus, req.ProductId, req.NftProductSizeId)
	if err != nil {
		return
	}
	nftCountData, err := NewBusinessNftMarketWarehouseTotalCount(context.Background()).GetByProductIdAndSizeId(int64(req.ProductId), int64(req.NftProductSizeId))
	if err != nil {
		return
	}
	nftCount = nftCountData.NftCount
	if nftCountData.NftCount == 0 {
		return
	}
	// 用户百分比（计算规则是普通用户剩余份数/app内展示剩余份数）
	userPercentage = (float64(realUserSurplus) / float64(nftCountData.NftCount)) * 100
	_, err = NewAiMatchProductNftSecondPrice(context.Background()).UpdateByParams(map[string]any{
		"product_id":          req.ProductId,
		"nft_product_size_id": req.NftProductSizeId,
	}, map[string]any{
		"user_percentage":   userPercentage,
		"real_user_surplus": realUserSurplus,
	})
	if err != nil {
		return
	}
	return
}

// 获取藏品的APP内展示剩余份数和用户拥有的总份数
func (n Nft) GetNftUserOwnNumAndAppShowNum(productId, nftProductSizeId int) (realUserSurplus int64, nftCount int64, err error) {
	// pass卡默认返回固定值
	if lo.IndexOf(constant.IGNORE_RESERVE_COLLECTION, uint64(productId)) > -1 {
		return 0, 100000, nil
	}
	err = GetProductOrderCount(n.Ctx, cli.SpecialUserIds, &realUserSurplus, productId, nftProductSizeId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	nftCountData, err := NewBusinessNftMarketWarehouseTotalCount(context.Background()).GetByProductIdAndSizeId(int64(productId), int64(nftProductSizeId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return realUserSurplus, 0, nil
	}
	nftCount = nftCountData.NftCount
	return
}
