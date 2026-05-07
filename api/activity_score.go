package api

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/dg_models"
	"hotbox-adm-backend/pkg/common"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

// @BasePath /hotbox/operation

// @Summary 获取待开始列表
// @Description 获取待开始列表
// @Tags 活动打分
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetPendingActivityScoreListReq true "查询参数"
// @Success 200 {object} any
// @Router /activity_score/pending/list [post]
func GetPendingActivityScoreList(c *gin.Context) {
	req := form.GetPendingActivityScoreListReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	result := []form.PendingActivityScoreListResp{}
	totalNum := 0
	offset := (req.PageNumber - 1) * req.PageSize
	// 先排除突袭
	activityListWhere := map[string][]any{
		"status":                 {1},
		"actual_activity_end_ts": {0},
		"activity_end_ts > ?":    {time.Now().UnixMilli()},
		"raid":                   {0},
		"activity_type":          {req.ActivityType},
	}
	activityList, err := dg_models.Activity{Ctx: c}.GetActivityByParamsV2(activityListWhere, []string{"activity_start_ts DESC"}, req.PageSize, offset)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if len(activityList) == 0 {
		response.ResponseSuccess(nil)
		return
	}
	activityDeActivityIds := lo.Map[dg_models.ActivityModel, int64](activityList, func(item dg_models.ActivityModel, index int) int64 {
		return item.DeActivityId
	})
	// 获取所有未结束的活动(按开始时间正序)
	switch req.ActivityType {
	// 合成
	case constant.ACTIVITY_TYPE_COMBINATION:
		// 获取所有未结束的活动(按开始时间正序)
		combinationList, total, err := models.NftCombinationDal.GetAiMatchProductNftCombinationByParams(c, map[string]any{
			"is_delete": 0,
			"id":        activityDeActivityIds,
		}, []string{"start_time DESC"}, req.PageSize, offset)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
		totalNum = int(total)
		g, _ := errgroup.WithContext(c)
		var mu sync.Mutex
		for _, v := range combinationList {
			activityData := v
			g.Go(func() error {
				// 初始化必要数据
				activityScoreData := models.ActivityScore{
					ActivityId:      activityData.ID,
					ActivityType:    req.ActivityType,
					ActivityTitle:   activityData.NftMapTitle,
					ActivityStartTs: activityData.StartTime,
				}
				affectRow, err := models.ActivityScoreDal.FirstOrCreate(c, activityScoreData)
				if err != nil {
					return err
				}
				if affectRow == 0 {
					activityScoreData, err = models.ActivityScoreDal.GetByUniqKey(c, int(activityData.ID), req.ActivityType)
					if err != nil {
						return err
					}
				}
				materials, outPutProduct, artist, notchNumber, err := HandleCombinationPendingActivityScoreList(activityData, activityScoreData.NotchNumber, req.ActivityType, false)
				if err != nil {
					return err
				}
				calculateActivityScoreResp, err := common.CalculateActivityScore(c, req.ActivityType, int(activityData.ID), materials, outPutProduct, artist, notchNumber, activityScoreData.ExpectedProductAmount, 0, false)
				if err != nil {
					return err
				}
				var buildActivityScoreListResp form.PendingActivityScoreListResp
				err = copier.Copy(&buildActivityScoreListResp, &activityScoreData)
				if err != nil {
					return err
				}
				buildActivityScoreListResp.Artist = artist
				buildActivityScoreListResp.Materials = materials
				buildActivityScoreListResp.OutPutProduct = outPutProduct
				buildActivityScoreListResp.TotalCost = calculateActivityScoreResp.TotalCost
				buildActivityScoreListResp.AntiFrictionLine = calculateActivityScoreResp.AntiFrictionLine
				buildActivityScoreListResp.ExpectedMarketProductValue = calculateActivityScoreResp.ExpectedMarketProductValue
				buildActivityScoreListResp.ExpectedProductCirculationAmount = calculateActivityScoreResp.ExpectedProductCirculationAmount
				buildActivityScoreListResp.Score = calculateActivityScoreResp.Score
				buildActivityScoreListResp.NotchNumber = notchNumber

				mu.Lock()
				result = append(result, buildActivityScoreListResp)
				mu.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			response.ResponseFail(err.Error())
			return
		}
	// 置换、分解
	case constant.ACTIVITY_TYPE_DECOMPOSE, constant.ACTIVITY_TYPE_REPLACE:
		// 获取所有未结束的活动(按开始时间正序)
		displaceList, total, err := models.NewAiMatchProductNftReplace().GetAiMatchProductNftReplaceByParams(c, map[string]any{
			"replace_id": activityDeActivityIds,
			"is_delete":  0,
		}, []string{"start_time DESC"}, req.PageSize, (int(req.PageNumber)-1)*req.PageSize)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
		totalNum = int(total)
		g, _ := errgroup.WithContext(c)
		var mu sync.Mutex
		for _, v := range displaceList {
			activityData := v
			g.Go(func() error {
				// 初始化必要数据
				activityScoreData := models.ActivityScore{
					ActivityId:      activityData.ReplaceId,
					ActivityType:    req.ActivityType,
					ActivityTitle:   activityData.ReplaceName,
					ActivityStartTs: activityData.StartTime,
				}
				affectRow, err := models.ActivityScoreDal.FirstOrCreate(c, activityScoreData)
				if err != nil {
					return err
				}
				if affectRow == 0 {
					activityScoreData, err = models.ActivityScoreDal.GetByUniqKey(c, int(activityData.ReplaceId), req.ActivityType)
					if err != nil {
						return err
					}
				}
				materials, outPutProduct, artist, notchNumber, err := HandleDisplacePendingActivityScoreList(activityData, activityScoreData.NotchNumber, req.ActivityType, false)
				if err != nil {
					return err
				}
				calculateActivityScoreResp, err := common.CalculateActivityScore(c, req.ActivityType, int(activityData.ReplaceId), materials, outPutProduct, artist, notchNumber, activityScoreData.ExpectedProductAmount, 0, false)
				if err != nil {
					return err
				}
				var buildActivityScoreListResp form.PendingActivityScoreListResp
				err = copier.Copy(&buildActivityScoreListResp, &activityScoreData)
				if err != nil {
					return err
				}
				buildActivityScoreListResp.Artist = artist
				buildActivityScoreListResp.Materials = materials
				buildActivityScoreListResp.OutPutProduct = outPutProduct
				buildActivityScoreListResp.TotalCost = calculateActivityScoreResp.TotalCost
				buildActivityScoreListResp.AntiFrictionLine = calculateActivityScoreResp.AntiFrictionLine
				buildActivityScoreListResp.ExpectedMarketProductValue = calculateActivityScoreResp.ExpectedMarketProductValue
				buildActivityScoreListResp.ExpectedProductCirculationAmount = calculateActivityScoreResp.ExpectedProductCirculationAmount
				buildActivityScoreListResp.Score = calculateActivityScoreResp.Score
				buildActivityScoreListResp.NotchNumber = notchNumber

				mu.Lock()
				result = append(result, buildActivityScoreListResp)
				mu.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ActivityStartTs > result[j].ActivityStartTs
	})
	response.ResponseSuccessWithList(result, totalNum)
}

// @Summary 获取已结束的活动列表
// @Description 获取已结束的活动列表
// @Tags 活动打分
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetEndedActivityScoreListReq true "查询参数"
// @Success 200 {object} any
// @Router /activity_score/ended/list [post]
func GetEndedActivityScoreList(c *gin.Context) {
	req := form.GetEndedActivityScoreListReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	result := []form.PendingActivityScoreListResp{}

	activityScores, total, err := models.ActivityScoreDal.GetEndActivityScorePage(c, req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	for _, v := range activityScores {
		act := v
		g, _ := errgroup.WithContext(c)
		var mu sync.Mutex
		g.Go(func() error {
			// 获取所有未结束的活动(按开始时间正序)
			switch req.ActivityType {
			// 合成
			case constant.ACTIVITY_TYPE_COMBINATION:
				// 获取所有未结束的活动(按开始时间正序)
				activityData, err := models.NftCombinationDal.One(c, act.ActivityId, "")
				if err != nil {
					return err
				}

				materials, outPutProduct, artist, notchNumber, err := HandleCombinationPendingActivityScoreList(activityData, act.NotchNumber, req.ActivityType, false)
				if err != nil {
					return err
				}
				calculateActivityScoreResp, err := common.CalculateActivityScore(c, req.ActivityType, int(activityData.ID), materials, outPutProduct, artist, notchNumber, v.ExpectedProductAmount, act.PraiseScore, true)
				if err != nil {
					return err
				}
				var buildActivityScoreListResp form.PendingActivityScoreListResp
				err = copier.Copy(&buildActivityScoreListResp, &act)
				if err != nil {
					return err
				}
				buildActivityScoreListResp.Artist = artist
				buildActivityScoreListResp.Materials = materials
				buildActivityScoreListResp.OutPutProduct = outPutProduct
				buildActivityScoreListResp.TotalCost = calculateActivityScoreResp.TotalCost
				buildActivityScoreListResp.AntiFrictionLine = calculateActivityScoreResp.AntiFrictionLine
				buildActivityScoreListResp.ExpectedMarketProductValue = calculateActivityScoreResp.ExpectedMarketProductValue
				buildActivityScoreListResp.ExpectedProductCirculationAmount = calculateActivityScoreResp.ExpectedProductCirculationAmount
				buildActivityScoreListResp.Score = calculateActivityScoreResp.Score

				buildActivityScoreListResp.RealProductCirculationAmount = calculateActivityScoreResp.RealProductCirculationAmount
				buildActivityScoreListResp.RealMarketProductValue = calculateActivityScoreResp.RealMarketProductValue
				buildActivityScoreListResp.Duration = act.ActivityDuration / 1000
				buildActivityScoreListResp.SaleMinPrice = calculateActivityScoreResp.SaleMinPrice
				buildActivityScoreListResp.SellPriceMinLimit = calculateActivityScoreResp.SellPriceMaxLimit
				buildActivityScoreListResp.ProductAvgCost = calculateActivityScoreResp.ProductAvgCost
				buildActivityScoreListResp.WarnTime = act.WarnTime
				buildActivityScoreListResp.RealScore = decimal.NewFromFloat(act.RecommendScoreBeforeEnd).Sub(decimal.NewFromInt(int64(act.WarnTime))).InexactFloat64()
				mu.Lock()
				result = append(result, buildActivityScoreListResp)
				mu.Unlock()
				return nil

				// 置换、分解
			case constant.ACTIVITY_TYPE_DECOMPOSE, constant.ACTIVITY_TYPE_REPLACE:
				// 获取所有未结束的活动(按开始时间正序)
				activityData, err := models.NewAiMatchProductNftReplace().GetByReplaceId(c, int(act.ActivityId))
				if err != nil {
					return err
				}
				// 初始化必要数据
				activityScoreData := models.ActivityScore{
					ActivityId:      activityData.ReplaceId,
					ActivityType:    req.ActivityType,
					ActivityTitle:   activityData.ReplaceName,
					ActivityStartTs: activityData.StartTime,
				}
				affectRaw, err := models.ActivityScoreDal.FirstOrCreate(c, activityScoreData)
				if err != nil {
					return err
				}
				if affectRaw == 0 {
					activityScoreData, err = models.ActivityScoreDal.GetByUniqKey(c, int(activityData.ReplaceId), req.ActivityType)
					if err != nil {
						return err
					}
				}
				materials, outPutProduct, artist, notchNumber, err := HandleDisplacePendingActivityScoreList(activityData, activityScoreData.NotchNumber, req.ActivityType, false)
				if err != nil {
					return err
				}
				calculateActivityScoreResp, err := common.CalculateActivityScore(c, req.ActivityType, int(activityData.ReplaceId), materials, outPutProduct, artist, notchNumber, v.ExpectedProductAmount, act.PraiseScore, true)
				if err != nil {
					return err
				}
				var buildActivityScoreListResp form.PendingActivityScoreListResp
				err = copier.Copy(&buildActivityScoreListResp, &activityScoreData)
				if err != nil {
					return err
				}
				buildActivityScoreListResp.Artist = artist
				buildActivityScoreListResp.Materials = materials
				buildActivityScoreListResp.OutPutProduct = outPutProduct
				buildActivityScoreListResp.TotalCost = calculateActivityScoreResp.TotalCost
				buildActivityScoreListResp.AntiFrictionLine = calculateActivityScoreResp.AntiFrictionLine
				buildActivityScoreListResp.ExpectedMarketProductValue = calculateActivityScoreResp.ExpectedMarketProductValue
				buildActivityScoreListResp.ExpectedProductCirculationAmount = calculateActivityScoreResp.ExpectedProductCirculationAmount
				buildActivityScoreListResp.Score = calculateActivityScoreResp.Score

				buildActivityScoreListResp.RealProductCirculationAmount = calculateActivityScoreResp.RealProductCirculationAmount
				buildActivityScoreListResp.RealMarketProductValue = calculateActivityScoreResp.RealMarketProductValue
				buildActivityScoreListResp.Duration = act.ActivityDuration / 1000
				buildActivityScoreListResp.SaleMinPrice = calculateActivityScoreResp.SaleMinPrice
				buildActivityScoreListResp.SellPriceMinLimit = calculateActivityScoreResp.SellPriceMaxLimit
				buildActivityScoreListResp.ProductAvgCost = calculateActivityScoreResp.ProductAvgCost
				buildActivityScoreListResp.WarnTime = act.WarnTime
				buildActivityScoreListResp.RealScore = decimal.NewFromFloat(act.RecommendScoreBeforeEnd).Sub(decimal.NewFromInt(int64(act.WarnTime))).InexactFloat64()
				mu.Lock()
				result = append(result, buildActivityScoreListResp)
				mu.Unlock()
				return nil
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}

	if req.OrderType == 2 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].WarnTime > result[j].WarnTime
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i].ActivityStartTs > result[j].ActivityStartTs
		})
	}
	response.ResponseSuccessWithList(result, int(total))
}

func HandleCombinationPendingActivityScoreList(activityData models.AiMatchProductNftCombinationModel, scoreNotchNumber int64, activityType int, initIfEnd bool) (materials dto.ActivityScoreMaterials, outPutProduct []dto.OutPutProduct, artist []dto.Artist, notchNumber int64, err error) {
	group := errgroup.Group{}
	group.Go(func() error {
		outPutProduct, artist, err = common.GetOutPutProductAndArtist(activityType, activityData)
		return err
	})
	group.Go(func() error {
		materials, err = common.GetMaterials(activityType, int(activityData.ID))
		return err
	})
	if err = group.Wait(); err != nil {
		return
	}
	// 缺口数
	if scoreNotchNumber > 0 {
		notchNumber = scoreNotchNumber
	} else {
		var (
			ctx                         = context.Background()
			slaveIds                    = []int64{}
			slaveHdTotalSuccessCountNum int64 // 子活动成功合成数
			maxNotchNumber              int64 // 最大缺口数
		)

		// 查询子活动
		if !lo.Contains([]string{"null", "0", "[]"}, activityData.SlaveId.String()) {
			if err = json.Unmarshal([]byte(activityData.SlaveId), &slaveIds); err != nil {
				return
			}
		}
		if len(slaveIds) > 0 {
			for _, v := range slaveIds {
				slaveCombinationModel, nerr := models.NftCombinationDal.One(ctx, v, "")
				if nerr != nil && !errors.Is(nerr, gorm.ErrRecordNotFound) {
					err = nerr
					return
				}
				if slaveCombinationModel.ID == 0 {
					continue
				}
				slaveHdSuccessCount, gerr := models.AiMatchProductNftCombinationLogDal.GetSuccessCountByCombineId(ctx, int(v))
				if gerr != nil {
					err = gerr
					return
				}
				slaveHdTotalSuccessCountNum += slaveHdSuccessCount * slaveCombinationModel.Cnt
				if maxNotchNumber < (slaveCombinationModel.RemainNum / slaveCombinationModel.Cnt) {
					maxNotchNumber = slaveCombinationModel.RemainNum / slaveCombinationModel.Cnt
				}
			}
		}

		if initIfEnd {
			// 已结束的合成活动产物缺口数的计算
			var hdSuccessCount int64
			hdSuccessCount, err = models.AiMatchProductNftCombinationLogDal.GetSuccessCountByCombineId(ctx, int(activityData.ID))
			if err != nil {
				return
			}
			hdSuccessCount += slaveHdTotalSuccessCountNum
			notchNumber = hdSuccessCount*activityData.Cnt + activityData.RemainNum/activityData.Cnt
		} else {
			if maxNotchNumber < (activityData.RemainNum / activityData.Cnt) {
				maxNotchNumber = activityData.RemainNum / activityData.Cnt
			}
			notchNumber = maxNotchNumber
		}
	}
	return
}

func HandleDisplacePendingActivityScoreList(activityData models.AiMatchProductNftReplaceModel, scoreNotchNumber int64, activityType int, initIfEnd bool) (materials dto.ActivityScoreMaterials, outPutProduct []dto.OutPutProduct, artist []dto.Artist, notchNumber int64, err error) {
	group := errgroup.Group{}
	group.Go(func() error {
		outPutProduct, artist, err = common.GetOutPutProductAndArtist(activityType, activityData)
		return err
	})
	group.Go(func() error {
		materials, err = common.GetMaterials(activityType, int(activityData.ReplaceId))
		return err
	})
	err = group.Wait()
	if err != nil {
		return
	}
	// 缺口数
	if scoreNotchNumber > 0 {
		notchNumber = scoreNotchNumber
	} else {
		if initIfEnd { // 初始化的时候活动是否结束
			result, ferr := models.AiMatchProductNftReplaceLogDal.GetByReplaceLogListByReplaceId(context.Background(), int(activityData.ReplaceId))
			if ferr != nil {
				logrus.Errorf("GetByReplaceLogListByReplaceId err, replaceId:%d", activityData.ReplaceId)
				err = errors.Wrapf(err, "GetByReplaceLogListByReplaceId err, replaceId:%d", activityData.ReplaceId)
				return
			}
			materialMap := map[string]int{}
			for _, item := range result {
				for _, res := range item.ResultInfo {
					if res.ResultType == "prop" {
						continue
					}
					materialMap[res.ProductName] += res.Num
				}
			}

			maxNum := 0
			for _, v := range materialMap {
				if maxNum < v {
					maxNum = v
				}
			}
			notchNumber = int64(maxNum) + (activityData.TotalCount - activityData.ReserveNum - activityData.ReplaceCount) // 已成功的加未置换/分解的
		} else {
			notchNumber = activityData.TotalCount - activityData.ReserveNum - activityData.ReplaceCount
		}
	}
	return
}

// @Summary 修改未开始活动
// @Description 修改未开始活动
// @Tags 活动打分
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.UpdatePendingActivityScoreReq true "查询参数"
// @Success 200 {object} any
// @Router /activity_score/pending/update [post]
func UpdatePendingActivityScore(c *gin.Context) {
	req := form.UpdatePendingActivityScoreReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if req.TotalCost != nil {
		err := models.ActivityScoreDal.Update(c, int(req.Id), map[string]any{
			"total_cost": req.TotalCost,
		})
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}
	response.ResponseSuccess(nil)
}

// @Summary 艺术家推荐分列表
// @Description 艺术家推荐分列表
// @Tags 活动打分
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.PagingReq true "查询参数"
// @Success 200 {object} any
// @Router /activity_score/artist_recommend_score/list [post]
func GetArtistRecommendScoreList(c *gin.Context) {
	req := form.PagingReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	list, count, err := models.ArtistRecommendScoreConfigDal.GetArtistRecommendScoreConfigList(c, nil, []string{"nft_num desc"}, req.PageNumber, req.PageSize)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccessWithList(list, int(count))
}

// @Summary 修改艺术家口碑分
// @Description 修改艺术家口碑分
// @Tags 活动打分
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.UpdateArtistRecommendScoreReq true "查询参数"
// @Success 200 {object} any
// @Router /activity_score/artist_recommend_score/update [post]
func UpdateArtistRecommendScore(c *gin.Context) {
	req := form.UpdateArtistRecommendScoreReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if cast.ToFloat64(req.Score) <= 0 {
		response.ResponseFail("口碑分必须大于0")
		return
	}

	if err := models.ArtistRecommendScoreConfigDal.UpdateParamsById(c, req.Id, map[string]any{"score": req.Score}); err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}
