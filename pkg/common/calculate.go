package common

import (
	"context"
	"sync"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/models/hd_adb_models"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type CalculateActivityScoreResp struct {
	TotalCost, // 总成本
	AntiFrictionLine, // 反撸线
	ExpectedProductCirculationAmount, // 产物预期流通份数
	ExpectedMarketProductValue, // 产物预期流通市值
	Score,
	RealProductCirculationAmount, // 产物实际流通份数
	RealMarketProductValue, // 产物实际流通市值
	SaleMinPrice, // 产物当前最低挂售价
	SellPriceMaxLimit, // 产物当前限价
	ProductAvgCost float64 // 产物平均成交价
}

func CalculateActivityScore(
	ctx *gin.Context,
	activityType, activityId int,
	materials dto.ActivityScoreMaterials,
	outPutProduct []dto.OutPutProduct,
	artist []dto.Artist,
	notchNumber int64, // 缺口数
	expectedProductAmount int64, // 预期流通份数
	praiseScored float64,
	activityEnd bool,
) (result CalculateActivityScoreResp, err error) {
	// 总成本
	activityScoreData, err := models.ActivityScoreDal.GetByUniqKey(ctx, activityId, activityType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	}
	// 总成本
	result.TotalCost = activityScoreData.TotalCost
	if result.TotalCost == 0 {
		result.TotalCost, err = CalculateTotalCost(materials)
		if err != nil {
			return result, err
		}
	}
	result.TotalCost = util.Decimal(result.TotalCost)
	// 反撸线
	result.AntiFrictionLine = util.Decimal(result.TotalCost * 1.1)
	// 艺术家口碑分
	artistReputationScore := float64(0)
	artistIds := lo.Map[dto.OutPutProduct, int](outPutProduct, func(item dto.OutPutProduct, index int) int {
		return item.ArtistData.ArtistId
	})
	artistList, err := models.ArtistRecommendScoreConfigDal.GetByParams(ctx, map[string]any{
		"artist_id": artistIds,
	})
	if err != nil {
		return result, err
	}

	if activityEnd {
		// 活动结束，口碑分不变
		artistReputationScore = praiseScored
	} else {
		if len(artistList) == 0 { // 艺术家为空，取材料的艺术家口碑分
			artistScore, err := getMaterialArtistScore(ctx, materials)
			if err != nil {
				return result, err
			}
			artistReputationScore = artistScore
		} else {
			artistRecommendScoreList := lo.Map[models.ArtistRecommendScoreConfig, float64](artistList, func(item models.ArtistRecommendScoreConfig, index int) float64 {
				return cast.ToFloat64(item.Score)
			})
			artistReputationScore = lo.Min(artistRecommendScoreList)
		}
	}

	// 产物预期流通份数(取最大值),缺口数*产物系数
	expectedProductCirculationAmountArr := []float64{}
	// 产物实际流通份数(取最大值)
	realProductCirculationAmountArr := []float64{}
	// 产物当前最低挂售价(取最大值)
	saleMinPriceArr := []float64{}
	// 产物当前限价(取最大值)
	sellPriceMaxLimitArr := []float64{}
	// 产物平均成交价 (取最大值)
	productAvgCostArr := []float64{}

	if expectedProductAmount == 0 {
		for _, v := range outPutProduct {
			// 道具跳过
			if v.ProductId == 0 {
				continue
			}
			expectedProductCirculationAmountArr = append(expectedProductCirculationAmountArr, float64(v.Num*int(notchNumber)))
		}
		result.ExpectedProductCirculationAmount = util.Decimal(lo.Max(expectedProductCirculationAmountArr))
	} else {
		result.ExpectedProductCirculationAmount = decimal.NewFromInt(expectedProductAmount).InexactFloat64()
	}

	for _, output := range outPutProduct {
		// 道具跳过
		if output.ProductId == 0 {
			continue
		}

		if activityEnd {
			remainCount, saleMinPrice, sellPriceMaxLimit, productAvgCost, err := getOutPutRealInfo(ctx, int64(output.ProductId), int64(output.NftProductSizeId))
			if err != nil {
				return result, errors.Wrapf(err, "CalculateActivityScore 活动结束获取产物实时信息错误，getOutPutRealInfo，productId=%v,nftProductSizeId=%v", output.ProductId, output.NftProductSizeId)
			}
			realProductCirculationAmountArr = append(realProductCirculationAmountArr, decimal.NewFromInt(remainCount).InexactFloat64())
			saleMinPriceArr = append(saleMinPriceArr, decimal.NewFromFloat(saleMinPrice).InexactFloat64())
			sellPriceMaxLimitArr = append(sellPriceMaxLimitArr, decimal.NewFromFloat(sellPriceMaxLimit).InexactFloat64())
			productAvgCostArr = append(productAvgCostArr, decimal.NewFromFloat(productAvgCost).InexactFloat64())

		}
	}

	// 产物预期流通市值 缺口数*反撸线
	result.ExpectedMarketProductValue = util.Decimal(decimal.NewFromInt(notchNumber).Mul(decimal.NewFromFloat(result.AntiFrictionLine)).InexactFloat64())
	// 产物限价
	result.SellPriceMaxLimit = util.Decimal(lo.Max(sellPriceMaxLimitArr))
	// 产物实际流通份数
	result.RealProductCirculationAmount = util.Decimal(lo.Max(realProductCirculationAmountArr))
	// 产物实际流通市值
	result.RealMarketProductValue = util.Decimal(float64(result.RealProductCirculationAmount * result.SellPriceMaxLimit))
	// 产物当前最低挂售价
	result.SaleMinPrice = util.Decimal(lo.Max(saleMinPriceArr))

	// 产品当日平均成交价
	result.ProductAvgCost = util.Decimal(lo.Max(productAvgCostArr))
	/*
		产物预期流通市值设为x（万），如产物预期流通市值=12300，则x=1.23万（易得x肯定是大于0的）
		（1）x＜1，y=50x+5
		（2）x∈[1,3.6)，y=−25𝑥²+85𝑥+20
		（3）x≥3.6，y=−7.8125x+78.125
	*/
	expectedMarketProductValueDivideWan := result.ExpectedMarketProductValue / 10000
	y := float64(0)
	if expectedMarketProductValueDivideWan < 1 {
		y = 50*expectedMarketProductValueDivideWan + 5
	} else if expectedMarketProductValueDivideWan >= 1 && expectedMarketProductValueDivideWan < 3.6 {
		y = -25*expectedMarketProductValueDivideWan*expectedMarketProductValueDivideWan + 85*expectedMarketProductValueDivideWan + 20
	} else {
		y = -7.8125*expectedMarketProductValueDivideWan + 78.125
	}
	result.Score = util.Decimal(0.3*artistReputationScore + 0.7*y)

	return result, nil
}

// 获取材料艺术家分数
func getMaterialArtistScore(ctx *gin.Context, materials dto.ActivityScoreMaterials) (float64, error) {
	materialIds := make([]int, 0)
	for _, v := range materials.MaterialsData.Materials {
		for _, item := range v {
			if item.ProductId == 0 { // 过滤道具
				continue
			}
			materialIds = append(materialIds, item.ProductId)
		}
	}
	if len(materialIds) == 0 {
		return 0, nil
	}
	artistData, err := models.SaleCalendarProduct{Ctx: ctx}.GetListJoinWithNftArtist(map[string]any{
		"sale_calendar_product.id": materialIds,
	})
	if err != nil {
		return 0, err
	}
	artistIds := []int64{}
	for _, v := range artistData {
		artistIds = append(artistIds, v.ArtistId)
	}
	if len(artistIds) == 0 {
		return 0, nil
	}
	artistList, err := models.ArtistRecommendScoreConfigDal.GetByParams(ctx, map[string]any{
		"artist_id": artistIds,
	})
	if err != nil {
		return 0, err
	}

	artistScore := make([]float64, 0)
	for _, v := range artistList {
		artistScore = append(artistScore, util.Decimal(cast.ToFloat64(v.Score)))
	}
	return lo.Min(artistScore), nil
}

func getOutPutRealInfo(ctx *gin.Context, productId, nftProductSizeId int64) (remainCount int64, saleMinPrice, sellPriceMaxLimit, productAvgCost float64, err error) {
	// 产物剩余流通份数
	var (
		realUserSurplus int64
		data            models.AiMatchProductNftSecondPrice
		avgCost         hd_adb_models.ProductAvgCost
		group           = errgroup.Group{}
	)
	group.Go(func() error {
		return models.GetProductOrderCount(ctx, cli.SpecialUserIds, &realUserSurplus, int(productId), int(nftProductSizeId))
	})

	group.Go(func() error {
		data, err = models.NewAiMatchProductNftSecondPrice(ctx).GetByProductIdAndNftProductSizeId(int(productId), int(nftProductSizeId))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location()).Format("2006-01-02 15:04:05")
		avgCost, err = hd_adb_models.AiMatchProductOrder{Ctx: ctx}.GetProductAvgCostByDate(int(productId), int(nftProductSizeId), startOfDay, endOfDay)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})
	if err = group.Wait(); err != nil {
		return 0, 0, 0, 0, err
	}
	remainCount = realUserSurplus
	saleMinPrice = data.SaleMinPrice
	sellPriceMaxLimit = data.SellPriceMaxLimit
	productAvgCost = avgCost.AvgPayAmount
	return
}

// 总成本，为其中一组材料的成本之和
// （成本：优先取当前最低挂售价，如取不到当前最低挂售价则取当日平均成交价，如取不到当日平均成交价，则取限价）。总成本=材料1系数*材料1成本+材料2系数*材料2成本……。当存在多组材料时，取其中总成本最低的一组。
func CalculateTotalCost(materials dto.ActivityScoreMaterials) (totalCost float64, err error) {
	totalCostArr := []float64{}
	for _, v := range materials.MaterialsData.Materials {
		var (
			productTotalCostMp = map[dto.ProductNftSizeId]float64{}
			group              = errgroup.Group{}
			mx                 = sync.Mutex{}
		)

		for _, vv := range v {
			if vv.ProductId == 0 {
				continue
			}
			item := vv
			group.Go(func() error {
				// 如果是0 的话，材料系数就取MaterialNum，否则就按材料内的配置
				coefficient := item.Num
				// 如果是合成活动，需要区分活动类型
				if materials.ActivityType == constant.ACTIVITY_TYPE_COMBINATION && item.MaterialType == 0 {
					coefficient = int(item.MaterialNum)
				}
				cost := float64(0)
				secondPriceData, err := models.NewAiMatchProductNftSecondPrice(context.Background()).GetByProductIdAndNftProductSizeId(item.ProductId, item.NftProductSizeId)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if secondPriceData.SaleMinPrice != 0 {
					cost = secondPriceData.SaleMinPrice
				} else {
					now := time.Now()
					startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
					endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location()).Format("2006-01-02 15:04:05")
					productAvgCost, err := hd_adb_models.AiMatchProductOrder{Ctx: &gin.Context{}}.GetProductAvgCostByDate(item.ProductId, item.NftProductSizeId, startOfDay, endOfDay)
					if err != nil {
						return err
					}
					cost = productAvgCost.AvgPayAmount
					if productAvgCost.AvgPayAmount == 0 {
						cost = secondPriceData.SellPriceMaxLimit
					}
				}
				cost = cost * float64(coefficient)

				mx.Lock()
				productTotalCostMp[dto.ProductNftSizeId{
					ProductId:        item.ProductId,
					NftProductSizeId: item.NftProductSizeId,
				}] += cost
				mx.Unlock()
				return nil
			})

		}
		tt := []float64{}
		for _, vv := range productTotalCostMp {
			tt = append(tt, vv)
		}
		totalCostArr = append(totalCostArr, lo.Min(tt))
	}
	return lo.Sum(totalCostArr), nil
}
