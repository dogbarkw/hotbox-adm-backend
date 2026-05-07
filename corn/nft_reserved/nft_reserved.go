package nft_reserved

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/util"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type NftReservedCornJob struct{}

func (p *NftReservedCornJob) Run() {
	StarNftReservedCornJob()
}

func StarNftReservedCornJob() {
	ctx := &gin.Context{}
	lock := cli.HotDogRedis.SetNX(ctx, "hd:nft_reserved_job:lock", "lock", 5*time.Minute).Val()
	if !lock {
		//_ = httpReq.FeiShuRootBot(fmt.Sprintf("[StarNftReservedCornJob] 没有获取到锁, time: %s", time.Now().Format(time.DateTime)))
		return
	}
	defer cli.HotDogRedis.Del(ctx, "hd:nft_reserved_job:lock")
	allReserved, err := models.ActivityMaterialReserveDal.GetByParams(ctx, map[string]any{
		"exec_status":     0,
		"exec_time <= ? ": time.Now().UnixMilli(),
		"exec_end_time":   0,
	})
	if err != nil {
		logrus.Error(err.Error())
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[StarNftReservedCornJob] ActivityMaterialReserveDal.GetByParams fail, time: %s, err: %v", time.Now().Format(time.DateTime), err))
		return
	}
	for _, reserve := range allReserved {
		switch reserve.ActivityType {
		case constant.ACTIVITY_TYPE_COMBINATION:
			reserveDetailJoinCombination, err := models.ActivityMaterialReserveDetailDal.GetReserveDetailJoinCombinationByParams(ctx, map[string]any{
				"ai_match_product_nft_activity_material_reserve_detail.activity_type": reserve.ActivityType,
				"ai_match_product_nft_activity_material_reserve_detail.reserve_id":    reserve.Id,
				"ai_match_product_nft_activity_material_reserve_detail.status":        0,
				"ai_match_product_nft_combination.on_sale_status":                     1,
			})
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			allReservedDetail := lo.Map(reserveDetailJoinCombination, func(item models.AiMatchProductNftActivityMaterialReserveDetailJoinCombination, index int) models.AiMatchProductNftActivityMaterialReserveDetail {
				return item.AiMatchProductNftActivityMaterialReserveDetail
			})
			err = handleNftReserved(ctx, reserve, allReservedDetail)
			if err != nil {
				logrus.Errorf("+%v", err.Error())
				return
			}
		case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
			reserveDetailJoinReplace, err := models.ActivityMaterialReserveDetailDal.GetReserveDetailJoinReplaceByParams(ctx, map[string]any{
				"ai_match_product_nft_activity_material_reserve_detail.activity_type": reserve.ActivityType,
				"ai_match_product_nft_activity_material_reserve_detail.reserve_id":    reserve.Id,
				"ai_match_product_nft_activity_material_reserve_detail.status":        0,
				"ai_match_product_nft_replace.on_sale_status":                         1,
			})
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			allReservedDetail := lo.Map(reserveDetailJoinReplace, func(item models.AiMatchProductNftActivityMaterialReserveDetailJoinReplace, index int) models.AiMatchProductNftActivityMaterialReserveDetail {
				return item.AiMatchProductNftActivityMaterialReserveDetail
			})
			err = handleNftReserved(ctx, reserve, allReservedDetail)
			if err != nil {
				logrus.Errorf("+%v", err.Error())
				return
			}
		}
	}
}

// 预留合材料
func handleNftReserved(ctx *gin.Context, reserve models.AiMatchProductNftActivityMaterialReserve, allReservedDetail []models.AiMatchProductNftActivityMaterialReserveDetail) error {
	reserveId := reserve.Id
	g, _ := errgroup.WithContext(ctx)
	for _, v := range allReservedDetail {
		item := v
		// 道具不需要预留
		if item.MaterialType == "prop" {
			continue
		}
		g.Go(func() error {
			productNftSecondPrice, _ := models.NewAiMatchProductNftSecondPrice(ctx).GetByProductIdAndNftProductSizeId(int(item.ProductId), int(item.ProductSizeId))
			product, err := models.SaleCalendarProduct{Ctx: ctx}.One(int64(item.ProductId))
			if err != nil {
				return errors.Wrapf(err, "sale_calendar_product query fail, 预留 product_id: %d", item.ProductId)
			}
			tx := cli.HotDogGormDB.WithContext(ctx).Begin()
			if product.MarketType == "copyright" {
				err := tx.Model(&models.SaleCalendarProductModel{}).Where("id = ?", item.ProductId).Updates(map[string]any{
					"sold_count": gorm.Expr("sold_count - ?", item.ReserveNum),
				}).Error
				if err != nil {
					tx.Rollback()
					return errors.Wrapf(err, "sale_calendar_product update sold_count fail, id: %d, ReserveNum: %d", item.ProductId, item.ReserveNum)
				}
			}
			businessNftMarketWareHouseTotalCountData := models.BusinessNftMarketWarehouseTotalCount{}
			err = tx.Model(&models.BusinessNftMarketWarehouseTotalCount{}).
				Where("product_id", item.ProductId).
				Where("product_size_id", item.ProductSizeId).
				First(&businessNftMarketWareHouseTotalCountData).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "business_nft_market_warehouse_total_count query fail,预留 product_id: %d, product_size_id: %d", item.ProductId, item.ProductSizeId)
			}
			err = tx.Model(&models.BusinessNftMarketWarehouseTotalCount{}).
				Where("product_id", item.ProductId).
				Where("product_size_id", item.ProductSizeId).
				Updates(map[string]any{
					"nft_count": gorm.Expr("nft_count - ?", item.ReserveNum),
				}).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "business_nft_market_warehouse_total_count update nft_count fail, product_id: %d, product_size_id: %d, ReserveNum=%d", item.ProductId, item.ProductSizeId, item.ReserveNum)
			}
			err = tx.Model(&models.AiMatchProductNftActivityMaterialReserveDetail{}).Where("id", item.Id).Update("status", 1).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "ai_match_product_nft_activity_material_reserve_detail update status fail, id: %d", item.Id)
			}
			requestData := struct {
				Remark string `json:"remark"`
			}{
				Remark: reserve.Remark,
			}
			jsonByt, _ := json.Marshal(&requestData)
			record := &models.AiMatchBackendOperateRecord{
				UserId:      reserve.UserId,
				Username:    reserve.UserName,
				AssociateId: productNftSecondPrice.ID,
				Scenes:      1,
				RequestData: string(jsonByt),
				Remark:      fmt.Sprintf("%s活动[id:%d] 调整剩余份数,变更前剩余份数: %d, 变更后剩余份数:%d", util.StringActivityType(reserve.ActivityType), reserve.ActivityId, businessNftMarketWareHouseTotalCountData.NftCount, businessNftMarketWareHouseTotalCountData.NftCount-item.ReserveNum),
			}
			err = tx.Model(&models.AiMatchBackendOperateRecord{}).Create(record).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "ai_match_backend_operate_record create fail, record: %+v", record)
			}
			return tx.Commit().Error
		})
	}
	if err := g.Wait(); err != nil {
		msg := fmt.Sprintf("[StarNftReservedCornJob] %s活动[id:%d],预留任务id:%d, group work fail, time: %s, err: %v",
			util.StringActivityType(reserve.ActivityType), reserve.ActivityId, reserveId, time.Now().Format(time.DateTime), err)
		logrus.Error(msg)
		_ = httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_NFT_RESERVED_URL, "活动预留", msg)
		return err
	}
	affectRow, err := models.ActivityMaterialReserveDal.UpdateByParams(map[string]any{
		"id": reserveId,
	}, map[string]any{
		"exec_status":   1,
		"exec_end_time": time.Now().UnixMilli(),
	})
	if err != nil {
		errMsg := fmt.Sprintf("%d 更新预留任务失败:%s", reserveId, err.Error())
		logrus.Error(errMsg)
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, "活动预留", errMsg)

	}
	logrus.Infof("%d 预留成功,allReservedDetail 数量为:%d", reserveId, len(allReservedDetail))
	if affectRow == 0 {
		errMsg := fmt.Sprintf("%d 更新预留任务影响数为0", reserveId)
		logrus.Error(errMsg)
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, "活动预留", errMsg)

	}
	return nil
}
