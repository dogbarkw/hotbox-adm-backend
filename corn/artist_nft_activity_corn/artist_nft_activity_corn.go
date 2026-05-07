package artist_nft_activity_corn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type ArtistNftActivityCountCornJob struct{}

func (p *ArtistNftActivityCountCornJob) Run() {
	ctx := &gin.Context{}
	lockKey := "hd:artist_nft_activity_count_job:lock"
	lock := cli.HotDogRedis.SetNX(ctx, lockKey, "lock", 5*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, lockKey)
	startArtistNftActivityCornJob()
}

func startArtistNftActivityCornJob() {
	ctx := &gin.Context{}
	// 统计藏品数，更新
	artistNftNums, err := models.SaleCalendarProduct{Ctx: ctx}.CountNftNum()
	if err != nil {
		logrus.Errorf("统计藏品数错误，err=%v", err.Error())
		httpReq.FeiShuDebugRootBot(fmt.Sprintf("统计藏品数错误，err=%v", err.Error()))
		return
	}
	for _, data := range artistNftNums {
		// logrus.Infof("data:%+v", data)
		if err := models.ArtistRecommendScoreConfigDal.SaveNftNum(ctx, &models.ArtistRecommendScoreConfig{
			ArtistId:   data.ArtistId,
			ArtistName: data.ArtistName,
			NftNum:     data.NftNum,
		}); err != nil {
			logrus.Errorf("保存艺术家藏品数错误，err=%v", err.Error())
			httpReq.FeiShuDebugRootBot(fmt.Sprintf("保存艺术家藏品数错误，err=%v", err.Error()))
			return
		}
	}

	// 每次获取1000条记录，查询藏品所属艺术家，并更新艺术家名下活动次数
	eg := errgroup.Group{}
	// 统计合成活动数
	eg.Go(func() error {
		return countCombinationActivity(ctx)
	})
	// 统计升级活动数
	eg.Go(func() error {
		return countUpgradeActivity(ctx)
	})
	// 统计置换，分解活动
	eg.Go(func() error {
		return countNftReplaceContent(ctx)
	})
	if err = eg.Wait(); err != nil {
		logrus.Errorf("统计艺术家活动报错,startArtistNftActivityCornJob, err=%+v", err)
		httpReq.FeiShuDebugRootBot(fmt.Sprintf("统计艺术家活动报错,startArtistNftActivityCornJob, err=%+v", err))
		return
	}
}

func countUpgradeActivity(ctx *gin.Context) error {
	limit := 1000
	activityUpgradeMaterialTableName := models.ActivityUpgradeMaterialDal.TableName()
	activityUpgradeMaterialId, err := models.TableScanRecordDal.GetLastScanId(ctx, activityUpgradeMaterialTableName, 1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "统计升级活动数,startArtistNftActivityCornJob, GetActivityUpgradeMaterialId")
	}
	for {
		list, err := models.ActivityUpgradeMaterialDal.GetByParamsLimit(ctx, map[string]any{"id > ?": cast.ToInt64(activityUpgradeMaterialId)}, limit)
		if err != nil {
			return errors.Wrap(err, "统计升级活动数,startArtistNftActivityCornJob, GetByParamsLimit")
		}
		for _, materialItem := range list {
			var materialInfo []dto.UpgradeMaterial
			if err = json.Unmarshal(materialItem.Detail, &materialInfo); err != nil {
				return errors.Wrap(err, "统计升级活动数,startArtistNftActivityCornJob, json.Unmarshal")
			}

			artistMap := map[int64]bool{} // 对艺术家ID去重，防止一个活动多次统计
			for _, material := range materialInfo {
				data, err := models.SaleCalendarProduct{Ctx: ctx}.One(material.ProductID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) { // 没有找到，跳过
						logrus.Warnf("统计升级活动数, SaleCalendarProduct, getByProductIdRecordNotFound， id=%d", material.ProductID)
						continue
					}
					return errors.Wrapf(err, "统计升级活动数, SaleCalendarProduct, getByProductId， id=%d", material.ProductID)
				}
				artistMap[int64(data.BrandId)] = true
			}

			err = cli.HotDogGormDB.Transaction(func(tx *gorm.DB) error {
				for artistId := range artistMap {
					err = models.ArtistRecommendScoreConfigDal.UpdateParamsByArtistId(ctx, artistId,
						map[string]any{"activity_num": gorm.Expr("activity_num + ?", 1)})
					if err != nil {
						return err
					}
				}

				return models.TableScanRecordDal.UpdateLastScanId(tx, &models.TableScanRecord{
					Table:    activityUpgradeMaterialTableName,
					TaskType: 1,
					LastId:   materialItem.ID,
				})
			})
			if err != nil {
				return errors.Wrap(err, "统计升级活动数,startArtistNftActivityCornJob, 事务操作")
			}

		}
		if len(list) < limit {
			break
		}
	}
	return nil
}

func countCombinationActivity(ctx *gin.Context) error {
	limit := 1000
	combinationTableName := models.NftCombinationModel.TableName()
	combinationActivityLastId, err := models.TableScanRecordDal.GetLastScanId(ctx, combinationTableName, 1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "统计合成活动数,startArtistNftActivityCornJob, GetCombinationMaterialLastId")
	}
	for {
		list, err := models.NftCombinationDal.GetByParamsLimit(ctx, map[string]any{"id > ?": cast.ToInt64(combinationActivityLastId)}, limit) //
		if err != nil {
			return errors.Wrap(err, "统计合成活动数,startArtistNftActivityCornJob, combinationActivityLastId")
		}
		for _, item := range list {

			combinationMaterials, err := models.CombinationMaterialDal.GetByParams(ctx, map[string]any{"combination_id": item.ID})
			if err != nil {
				return errors.Wrap(err, "查询合成活动材料,startArtistNftActivityCornJob, CombinationMaterialDal.GetByParams")
			}

			artistMap := map[int64]bool{} // 对艺术家ID去重，防止一个活动多次统计
			for _, materialItem := range combinationMaterials {
				var materialInfo []dto.CombinationMaterial
				if err = json.Unmarshal(materialItem.MaterialInfo, &materialInfo); err != nil {
					return errors.Wrap(err, "统计合成活动数,startArtistNftActivityCornJob, json.Unmarshal")
				}

				for _, material := range materialInfo {
					if material.SubProductID == 0 {
						continue
					}
					data, err := models.SaleCalendarProduct{Ctx: ctx}.One(material.SubProductID)
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) { // 没有找到，跳过
							logrus.Warnf("统计合成活动数, SaleCalendarProduct, getByProductIdRecordNotFound， id=%d", material.SubProductID)
							continue
						}
						return errors.Wrapf(err, "统计合成活动数, SaleCalendarProduct, getByProductId，id=%d", material.SubProductID)
					}
					artistMap[data.BrandId] = true
				}
			}

			err = cli.HotDogGormDB.Transaction(func(tx *gorm.DB) error {
				for artistId := range artistMap {
					err = models.ArtistRecommendScoreConfigDal.UpdateParamsByArtistId(ctx, artistId,
						map[string]any{"activity_num": gorm.Expr("activity_num + ?", 1)})
					if err != nil {
						return err
					}
				}
				return models.TableScanRecordDal.UpdateLastScanId(tx, &models.TableScanRecord{
					Table:    combinationTableName,
					TaskType: 1,
					LastId:   item.ID,
				})
			})
			if err != nil {
				return errors.Wrap(err, "统计合成活动数,startArtistNftActivityCornJob, 事务操作")
			}

		}
		if len(list) < limit {
			break
		}
	}
	return nil
}

func countNftReplaceContent(ctx *gin.Context) error {
	limit := 1000
	nftReplaceContent := models.AiMatchProductNftReplaceContentDal.TableName()
	nftReplaceContentId, err := models.TableScanRecordDal.GetLastScanId(ctx, nftReplaceContent, 1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "统计置换分解活动,startArtistNftActivityCornJob, GetNftReplaceContentId")
	}
	for {
		list, err := models.AiMatchProductNftReplaceContentDal.GetByParamsLimit(ctx, map[string]any{"replace_content_id > ?": cast.ToInt64(nftReplaceContentId)}, limit)
		if err != nil {
			return errors.Wrap(err, "统计置换分解活动,startArtistNftActivityCornJob, GetByParamsLimit")
		}

		for _, content := range list {
			if content.ProductId == 0 {
				continue
			}
			data, err := models.SaleCalendarProduct{Ctx: ctx}.One(content.ProductId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) { // 没有找到，跳过
					logrus.Warnf("统计合成活动数, SaleCalendarProduct, getByProductIdRecordNotFound， id=%d,nftReplaceContentId=%d", content.ProductId, content.ReplaceContentId)
					continue
				}
				return errors.Wrapf(err, "统计置换分解活动, SaleCalendarProduct, getByProductId，id=%d,nftReplaceContentId=%d", content.ProductId, content.ReplaceContentId)
			}

			err = cli.HotDogGormDB.Transaction(func(tx *gorm.DB) error {
				err = models.ArtistRecommendScoreConfigDal.UpdateParamsByArtistId(ctx, data.BrandId,
					map[string]any{"activity_num": gorm.Expr("activity_num + ?", 1)})
				if err != nil {
					return err
				}

				return models.TableScanRecordDal.UpdateLastScanId(tx, &models.TableScanRecord{
					Table:    nftReplaceContent,
					TaskType: 1,
					LastId:   content.ReplaceContentId,
				})
			})
			if err != nil {
				return errors.Wrap(err, "统计置换分解活动, 事务操作err")
			}
		}

		if len(list) < limit {
			break
		}
	}
	return nil
}
