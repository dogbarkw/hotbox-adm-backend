package activity_score_nft_second_price

import (
	"fmt"
	"time"

	"hotbox-adm-backend/models/dg_models"
	"hotbox-adm-backend/models/hd_adb_models"
	"hotbox-adm-backend/util"

	"github.com/samber/lo"

	"hotbox-adm-backend/api"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/pkg/common"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"

	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
)

type ActivityScoreInitCorn struct{}

func (p *ActivityScoreInitCorn) Run() {
	ctx := &gin.Context{}
	lockKey := "hd:activity_score_init_job:lock"
	lock := cli.HotDogRedis.SetNX(ctx, lockKey, "lock", 5*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, lockKey)

	group := errgroup.Group{}
	//group.Go(func() error {
	//	// 初始化8.1之后已结束的活动数据
	//	return startActivityEndData(ctx)
	//})
	group.Go(func() error {
		// 活动前初始化
		return startActivityStartData(ctx)
	})
	group.Go(func() error {
		// 活动已结束，保存活动时长，口碑分，推荐分
		return checkActivityIfEnd(ctx)
	})
	group.Go(func() error {
		// 监听产物是否报警
		return startCheckNftSecondPrice(ctx)
	})
	if err := group.Wait(); err != nil {
		logrus.Errorf("ActivityScoreInitCorn err:%+v", err)
		httpReq.FeiShuDebugRootBot(fmt.Sprintf("ActivityScoreInitCorn err:%+v", err))
	}
}

// 提前初始化活动
func startActivityStartData(ctx *gin.Context) error {
	timeNow := time.Now()
	// 非突袭活动
	// startTime := timeNow.UnixMilli() - 5*60*1000 // 从活动开始前5分钟开始检查
	// startTime := 1722441600000 // 只取8.1开始的活动
	list, err := dg_models.Activity{Ctx: ctx}.GetActivityByParams(map[string][]any{
		"activity_start_ts >=? and activity_start_ts <?": {timeNow.UnixMilli(), timeNow.UnixMilli() + 5*60*1000},
		"activity_type in (?)":                           {[]int{constant.ACTIVITY_TYPE_COMBINATION, constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE}},
		"is_delete":                                      {0},
		"raid":                                           {0},
	}, []string{"id desc"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "获取非突袭活动列表错误")
	}

	// 突袭活动
	raidList, err := dg_models.Activity{Ctx: ctx}.GetActivityByParams(map[string][]any{
		"on_shelf_ts >=? and on_shelf_ts <?": {timeNow.UnixMilli(), timeNow.UnixMilli() + 5*60*1000},
		"activity_type in (?)":               {[]int{constant.ACTIVITY_TYPE_COMBINATION, constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE}},
		"is_delete":                          {0},
		"raid":                               {1},
	}, []string{"id desc"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "获取突袭活动列表错误")
	}
	for _, actItem := range list {
		notchNumber, totalCost, expectedProductCirculationAmount, err := activityInitCalculate(ctx, actItem)
		if err != nil {
			return err
		}
		if _, err = models.ActivityScoreDal.ActivityInitAndUpdate(ctx, models.ActivityScore{
			ActivityId:            actItem.DeActivityId,
			ActivityType:          actItem.ActivityType,
			ActivityTitle:         actItem.ActivityTitle,
			ActivityStartTs:       actItem.ActivityStartTs,
			NotchNumber:           notchNumber,
			ExpectedProductAmount: int64(expectedProductCirculationAmount),
			TotalCost:             totalCost,
		}); err != nil {
			return err
		}
	}

	for _, actItem := range raidList {
		notchNumber, totalCost, expectedProductCirculationAmount, err := activityInitCalculate(ctx, actItem)
		if err != nil {
			return err
		}
		if _, err = models.ActivityScoreDal.ActivityInitAndUpdate(ctx, models.ActivityScore{
			ActivityId:            actItem.DeActivityId,
			ActivityType:          actItem.ActivityType,
			ActivityTitle:         actItem.ActivityTitle,
			ActivityStartTs:       actItem.OnShelfTs, // 突袭活动使用上架时间
			NotchNumber:           notchNumber,
			ExpectedProductAmount: int64(expectedProductCirculationAmount),
			TotalCost:             totalCost,
		}); err != nil {
			return err
		}
	}
	return nil
}

// 初始化8.1之后结束的活动，初始化完之后可以移除
func startActivityEndData(ctx *gin.Context) error {
	timeNow := time.Now()
	// 非突袭活动
	startTime := 1722441600000 // 只取8.1开始的活动
	list, err := dg_models.Activity{Ctx: ctx}.GetActivityByParams(map[string][]any{
		"activity_start_ts >=? and activity_start_ts <?": {startTime, timeNow.UnixMilli()},
		"actual_activity_end_ts > ?":                     {0},
		"activity_type in (?)":                           {[]int{constant.ACTIVITY_TYPE_COMBINATION, constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE}},
		"is_delete":                                      {0},
		"raid":                                           {0},
	}, []string{"id desc"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "获取非突袭活动列表错误")
	}

	// 突袭活动
	raidList, err := dg_models.Activity{Ctx: ctx}.GetActivityByParams(map[string][]any{
		"on_shelf_ts >=? and on_shelf_ts <?": {startTime, timeNow.UnixMilli()},
		"actual_activity_end_ts > ?":         {0},
		"activity_type in (?)":               {[]int{constant.ACTIVITY_TYPE_COMBINATION, constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE}},
		"is_delete":                          {0},
		"raid":                               {1},
	}, []string{"id desc"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "获取突袭活动列表错误")
	}
	for _, actItem := range list {
		notchNumber, totalCost, expectedProductCirculationAmount, err := activityInitCalculate(ctx, actItem)
		if err != nil {
			return err
		}
		if _, err = models.ActivityScoreDal.ActivityInitAndUpdate(ctx, models.ActivityScore{
			ActivityId:            actItem.DeActivityId,
			ActivityType:          actItem.ActivityType,
			ActivityTitle:         actItem.ActivityTitle,
			ActivityStartTs:       actItem.ActivityStartTs,
			NotchNumber:           notchNumber,
			ExpectedProductAmount: int64(expectedProductCirculationAmount),
			TotalCost:             totalCost,
		}); err != nil {
			return err
		}
	}

	for _, actItem := range raidList {
		notchNumber, totalCost, expectedProductCirculationAmount, err := activityInitCalculate(ctx, actItem)
		if err != nil {
			return err
		}
		if _, err = models.ActivityScoreDal.ActivityInitAndUpdate(ctx, models.ActivityScore{
			ActivityId:            actItem.DeActivityId,
			ActivityType:          actItem.ActivityType,
			ActivityTitle:         actItem.ActivityTitle,
			ActivityStartTs:       actItem.OnShelfTs, // 突袭活动使用上架时间
			NotchNumber:           notchNumber,
			ExpectedProductAmount: int64(expectedProductCirculationAmount),
			TotalCost:             totalCost,
		}); err != nil {
			return err
		}
	}
	return nil
}

func activityInitCalculate(ctx *gin.Context, actItem dg_models.ActivityModel) (int64, float64, float64, error) {
	var (
		notchNumber   int64 // 活动缺口数
		materials     dto.ActivityScoreMaterials
		outPutProduct []dto.OutPutProduct
		err           error
	)

	// 获取材料
	switch actItem.ActivityType {
	case constant.ACTIVITY_TYPE_COMBINATION:
		activityData, err := models.NftCombinationDal.One(ctx, actItem.DeActivityId, "")
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, 0, errors.Wrap(err, "获取合成活动错误")
		}
		materials, outPutProduct, _, notchNumber, err = api.HandleCombinationPendingActivityScoreList(activityData, 0, actItem.ActivityType, actItem.ActualActivityEndTs > 0)
		if err != nil {
			return 0, 0, 0, err
		}
	case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
		activityData, err := models.NewAiMatchProductNftReplace().GetByReplaceId(ctx, int(actItem.DeActivityId))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, 0, errors.Wrap(err, "获取置换分解活动错误")
		}
		materials, outPutProduct, _, notchNumber, err = api.HandleDisplacePendingActivityScoreList(activityData, 0, actItem.ActivityType, actItem.ActualActivityEndTs > 0)
		if err != nil {
			return 0, 0, 0, err
		}
	default: // 忽略其他活动
		return 0, 0, 0, nil
	}

	totalCost, err := common.CalculateTotalCost(materials)
	if err != nil {
		return 0, 0, 0, err
	}
	var expectedProductCirculationAmountArr []float64
	for _, v := range outPutProduct {
		// 道具跳过
		if v.ProductId == 0 {
			continue
		}
		expectedProductCirculationAmountArr = append(expectedProductCirculationAmountArr, float64(v.Num*int(notchNumber)))
	}

	return notchNumber, totalCost, lo.Max(expectedProductCirculationAmountArr), err
}

func checkActivityIfEnd(ctx *gin.Context) error {
	// 活动是否结束
	activityScores, err := models.ActivityScoreDal.GetByParams(ctx, map[string]any{"status": 0}, []string{"id desc"})
	if err != nil {
		return errors.Wrap(err, "获取活动打分列表错误")
	}
	for _, score := range activityScores {
		result, err := dg_models.Activity{Ctx: ctx}.GetByActivityIdAndType(score.ActivityId, int64(score.ActivityType))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		logrus.Infof("活动id:%d,活动类型:%d,活动实际结束实际:%d", score.ActivityId, score.ActivityType, result.ActualActivityEndTs)
		if result.ActualActivityEndTs == 0 {
			continue
		}

		// 活动已结束
		score.Status = 1
		score.ActivityDuration = result.ActualActivityEndTs - score.ActivityStartTs
		// 获取艺术家，多个艺术家取最低口碑分
		var artist []dto.Artist
		switch score.ActivityType {
		case constant.ACTIVITY_TYPE_COMBINATION:
			act, err := models.NftCombinationDal.One(ctx, score.ActivityId, "")
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrapf(err, "获取合成活动错误,id=%d", score.ActivityId)
			}
			materials, outPutProduct, arts, notchNumber, err := api.HandleCombinationPendingActivityScoreList(act, score.NotchNumber, score.ActivityType, false)
			if err != nil {
				return err
			}
			calculateActivityScoreResp, err := common.CalculateActivityScore(ctx, score.ActivityType, int(score.ActivityId), materials, outPutProduct, arts, notchNumber, score.ExpectedProductAmount, 0, false)
			if err != nil {
				return err
			}
			score.RecommendScoreBeforeEnd = calculateActivityScoreResp.Score
			artist = arts
		case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
			act, err := models.NewAiMatchProductNftReplace().GetByReplaceId(ctx, int(score.ActivityId))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "获取置换分解活动错误")
			}
			materials, outPutProduct, arts, notchNumber, err := api.HandleDisplacePendingActivityScoreList(act, score.NotchNumber, score.ActivityType, false)
			if err != nil {
				return err
			}
			calculateActivityScoreResp, err := common.CalculateActivityScore(ctx, score.ActivityType, int(score.ActivityId), materials, outPutProduct, arts, notchNumber, score.ExpectedProductAmount, 0, false)
			if err != nil {
				return err
			}
			score.RecommendScoreBeforeEnd = calculateActivityScoreResp.Score
			artist = arts
		}

		artistIds := make([]int, 0, len(artist))
		for _, artistItem := range artist {
			artistIds = append(artistIds, artistItem.ArtistId)
		}
		artistIds = lo.Uniq(artistIds)
		if len(artistIds) > 1 { // 多艺术家
			logrus.Infof("活动ID:%v,活动类型:%v,多艺术家，艺术家id=%v", score.ActivityId, score.ActivityType, artistIds)
		}

		if len(artistIds) > 0 {
			artistRecommendScoreConfigs, err := models.ArtistRecommendScoreConfigDal.GetByParams(ctx, map[string]any{"artist_id in (?)": artistIds})
			if err != nil {
				return err
			}
			if len(artistRecommendScoreConfigs) == 0 {
				logrus.Warnf("活动ID:%v,活动类型:%v,艺术家推荐分设置列表找不到艺术家，艺术家ID=%v", score.ActivityId, score.ActivityType, artistIds)
				// httpReq.FeiShuDebugRootBot(fmt.Sprintf("活动ID:%v,活动类型:%v,艺术家推荐分设置列表找不到艺术家，艺术家ID=%v", score.ActivityId, score.ActivityType, artistIds))
				continue
			}

			minScore := artistRecommendScoreConfigs[0].Score
			for _, config := range artistRecommendScoreConfigs {
				if cast.ToFloat64(minScore) > cast.ToFloat64(config.Score) {
					minScore = config.Score
				}
			}
			scoreDecimal, err := decimal.NewFromString(minScore)
			if err != nil {
				return errors.Wrapf(err, "活动ID:%v,活动类型:%v,NewFromStringErr,minScore=%v", score.ActivityId, score.ActivityType, minScore)
			}

			score.PraiseScore = scoreDecimal.InexactFloat64()
		}
		if err = models.ActivityScoreDal.Save(ctx, score); err != nil {
			return err
		}
	}

	return nil
}

func startCheckNftSecondPrice(ctx *gin.Context) error {
	activityScores, err := models.ActivityScoreDal.GetByParams(ctx, map[string]any{"status": 1}, []string{"id desc"})
	// activityScores, err := models.ActivityScoreDal.GetByParams(ctx, map[string]any{"status": 1, "activity_id": 2183}, []string{"id desc"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(activityScores) == 0 {
		return nil
	}

	timeNow := time.Now()
	for _, score := range activityScores {
		var activityData any
		switch score.ActivityType {
		case constant.ACTIVITY_TYPE_COMBINATION:
			act, err := models.NftCombinationDal.One(ctx, score.ActivityId, "")
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "获取合成活动错误")
			}
			activityData = act
		case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
			act, err := models.NewAiMatchProductNftReplace().GetByReplaceId(ctx, int(score.ActivityId))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "获取置换分解活动错误")
			}
			activityData = act
		}
		outPutProduct, _, err := common.GetOutPutProductAndArtist(score.ActivityType, activityData)
		if err != nil {
			return errors.Wrap(err, "获取产物和艺术家错误")
		}

		for _, outPut := range outPutProduct {
			// 检查产物，是否报警
			if err = checkOutPutNftWarn(ctx, timeNow, score, outPut); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkOutPutNftWarn(ctx *gin.Context, timeNow time.Time, score models.ActivityScore, outPutProduct dto.OutPutProduct) error {
	queryParam := map[string][]any{
		"on_sale_status":                   {1}, // 上架
		"second_sale_time between ? and ?": {timeNow.UnixMilli() - 4*1000*3600, timeNow.UnixMilli()},
		"is_delete":                        {0},
		"product_id":                       {outPutProduct.ProductId},
		"nft_product_size_id":              {outPutProduct.NftProductSizeId},
	}
	nftInfo, err := models.AiMatchProductNftSecondPrice{}.GetOneByParams(ctx, queryParam)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(err, "获取寄售活动列表错误， GetListByParams，params=%v", queryParam)
	}
	if nftInfo.ID == 0 {
		return nil
	}

	subTime := timeNow.Unix() - nftInfo.SecondSaleTime/1000
	if (subTime/60)%10 == 0 { // 每10分钟开始快照
		productAvgCost, err := hd_adb_models.AiMatchProductOrder{Ctx: &gin.Context{}}.GetProductAvgCostByDate(outPutProduct.ProductId, outPutProduct.NftProductSizeId, time.Now().Format("2006-01-02 00:00:00"), time.Now().Format("2006-01-02 23:59:59"))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(err, "获取产品当日平均成交价错误，GetProductAvgCostByDate，productId=%v,nftProductSizeId=%v", outPutProduct.ProductId, outPutProduct.NftProductSizeId)
		}
		// min{产物当前最低挂售价,产物当前限价,产物平均成交价}＜反撸线*70%
		// 产物当前限价＜反撸线*85%
		compare := []float64{}
		for _, v := range []float64{nftInfo.SaleMinPrice, nftInfo.SellPriceMaxLimit, productAvgCost.AvgPayAmount} {
			if v == 0 { // 过滤掉0
				continue
			}
			compare = append(compare, v)
		}
		minPrice := lo.Min(compare)
		minPrice = util.Decimal(minPrice)
		if outPutProduct.Num == 0 {
			return errors.Wrapf(err, "获取产物数量错误，FirstOrUpdate，outPutProduct=%v", outPutProduct)
		}
		antiFrictionLine := util.Decimal(score.TotalCost * 1.1 / float64(outPutProduct.Num) * 0.7)
		antiFrictionLine1 := util.Decimal(score.TotalCost * 1.1 / float64(outPutProduct.Num) * 0.85)

		isWarn := false
		if minPrice < antiFrictionLine || antiFrictionLine1 > nftInfo.SellPriceMaxLimit {
			// 报警次数+1
			score.WarnTime += 1
			isWarn = true
			if err = models.ActivityScoreDal.Save(ctx, score); err != nil {
				return errors.Wrapf(err, "更新活动报警次数错误，FirstOrUpdate，score=%v", score)
			}
		}

		httpReq.FeiShuDebugRootBot(fmt.Sprintf("活动ID：%v，活动名称：%v，开始检查产物报警，产物ID：%v, 产物名称：%v，min{最低挂售价,当前限价,平均成交价}=%v,（反撸线/产物系数）*70=%v,当前限价%v,（反撸线/产物系数）*85=%v, 触发报警：%v，当前时间：%s",
			score.ActivityId, score.ActivityTitle, outPutProduct.ProductId, outPutProduct.Name, minPrice, antiFrictionLine, nftInfo.SellPriceMaxLimit, antiFrictionLine1, isWarn, timeNow.Format("2006-01-02 15:04:05")))
	}
	return nil
}
