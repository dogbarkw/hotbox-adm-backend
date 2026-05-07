package recycle_record

import (
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/api"
	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const RECYCLE_RECORD_LOCK_KEY = "hd:recycle_record_lock:"

type RecycleRecordCornJob struct{}

func (p *RecycleRecordCornJob) Run() {
	StartNftRecycleRecordJob()
}

func StartNftRecycleRecordJob() {
	ctx := &gin.Context{}
	lock := cli.HotDogRedis.SetNX(ctx, "hd:nft_recycle_record_job:lock", "lock", 5*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, "hd:nft_recycle_record_job:lock")
	allRecycleRecordList, err := models.AiMatchProductNftRecycleRecord{}.GetByParams(ctx, map[string]any{
		"type":   constant.BATCH_RECORD_RECYCLE_TYPE,
		"status": 0,
	}, 5)
	if err != nil {
		logrus.Errorf("StartNftRecycleRecordJob 查询回收记录表报错:%s", err.Error())
		return
	}
	g, _ := errgroup.WithContext(ctx)
	for _, ar := range allRecycleRecordList {
		recycleRecord := ar
		g.Go(func() error {
			// 对每个用户单独藏品枷锁
			key := fmt.Sprintf("%s%d_%d_%d", RECYCLE_RECORD_LOCK_KEY, recycleRecord.UserId, recycleRecord.ProductId, recycleRecord.NftProductSizeId)
			lock := cli.HotDogRedis.SetNX(ctx, key, "lock", 5*time.Minute).Val()
			if !lock {
				return nil
			}
			defer cli.HotDogRedis.Del(ctx, key)
			// 避免重启时重复获取，这里用的是目标数-已完成数来进行查询
			targetCount := recycleRecord.RecycleTargetCount - recycleRecord.RecycleCount
			// 获取相应的订单
			orderList, err := models.AiMatchProductOrder{Ctx: ctx}.GetProductOrderList(map[string]any{
				"receiver_name":       cast.ToString(recycleRecord.UserId),
				"product_type":        "NFT",
				"status":              2,
				"is_delete":           0,
				"apply_refund_status": 0,
				"nft_product_size_id": recycleRecord.NftProductSizeId,
				"product_id":          recycleRecord.ProductId,
			}, &targetCount)
			if err != nil {
				logrus.Error(err)
				return nil
			}
			if len(orderList) != targetCount {
				return nil
			}
			for _, v := range orderList {
				if v.MarketType == "land" {
					models.AiMatchProductNftRecycleRecord{}.Update(ctx, int(recycleRecord.ID), map[string]any{
						"status": -1,
						"msg":    "暂不支持回收土地",
					})
					return nil

				} else if lo.Contains[string]([]string{"***", ""}, v.ReceiverCity) {
					models.AiMatchProductNftRecycleRecord{}.Update(ctx, int(recycleRecord.ID), map[string]any{
						"status": -1,
						"msg":    "订单还没发编号，不能回收",
					})
					return nil
				}
			}
			err = handleRecycleOrder(ctx, orderList, recycleRecord)
			if err != nil {
				models.AiMatchProductNftRecycleRecord{}.Update(ctx, int(recycleRecord.ID), map[string]any{
					"status": -1,
					"msg":    err.Error(),
				})
				return nil
			}
			return nil
		})
	}
	g.Wait()
}

func handleRecycleOrder(c *gin.Context, orders []models.AiMatchProductOrderModel, recycleRecord models.AiMatchProductNftRecycleRecord) error {
	// 组装事务sql

	orderChunk := lo.Chunk(orders, 50)
	for k, v := range orderChunk {
		orderIds := make([]int64, 0)
		tx := cli.HotDogGormDB.WithContext(c).Begin()
		for _, order := range v {

			orderIds = append(orderIds, order.ID)
			err := tx.Model(&models.AiMatchProductOrderModel{}).Where("id", order.ID).
				Updates(map[string]any{
					"status":    67,
					"note":      "RECYCLING",
					"is_delete": 1,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			var count int64
			err = tx.Model(&models.AiMatchRestNftProductModel{}).Where(map[string]any{
				"product_id":          order.ProductId,
				"nft_product_size_id": order.NftProductSizeId,
				"receiver_city":       order.ReceiverCity,
				"receiver_province":   order.ReceiverProvince,
				"is_release":          0,
				"is_delete":           0,
			}).Count(&count).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			if count > 0 {
				continue
			}
			if order.NftProductSizeId > 0 {
				iosSourceFile := ""
				androidSourceFile := ""
				if lo.Contains[string](api.PFP_LIST, order.SourceType) {
					var productSize models.SaleProductNftSizePfpModel
					err := tx.Where(map[string]any{
						"product_id":          order.ProductId,
						"nft_product_size_id": order.NftProductSizeId,
						"receiver_city":       order.ReceiverCity,
						"is_delete":           0,
					}).First(&productSize).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							tx.Rollback()
							return api.RecycleNoNftSizeError
						}
						tx.Rollback()
						return err
					}
					iosSourceFile = productSize.OriginMedia
					androidSourceFile = productSize.OriginMedia

				} else {
					var productSize models.SaleProductNftSizeModel
					err := tx.Where(map[string]any{
						"nft_product_size_id": order.NftProductSizeId,
					}).First(&productSize).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							tx.Rollback()
							return api.RecycleNoNftSizeError
						}
						tx.Rollback()
						return err
					}
					iosSourceFile = productSize.IosSourceFile
					androidSourceFile = productSize.AndroidSourceFile
				}
				err := tx.Create(&models.AiMatchRestNftProductModel{
					ProductId:         order.ProductId,
					ProductTitle:      order.ProductTitle,
					SizeId:            order.SizeId,
					NftProductSizeId:  order.NftProductSizeId,
					ActivePicture:     order.ProductPicture,
					IosSourceFile:     iosSourceFile,
					AndroidSourceFile: androidSourceFile,
					ReceiverProvince:  order.ReceiverProvince,
					ReceiverCity:      order.ReceiverCity,
					ReceiverRegion:    order.ReceiverRegion,
					CombineInfo:       order.LogisticsInfo,
					CombineActiveId:   order.NewFlashId,
					PlayMethod:        order.PlayMethod,
					CouponIds:         datatypes.JSON([]byte(`[]`)),
				}).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			} else {
				err := tx.Create(&models.AiMatchRestNftProductModel{
					ProductId:         order.ProductId,
					ProductTitle:      order.ProductTitle,
					SizeId:            order.SizeId,
					NftProductSizeId:  order.NftProductSizeId,
					ActivePicture:     order.ProductPicture,
					IosSourceFile:     "",
					AndroidSourceFile: "",
					ReceiverProvince:  order.ReceiverProvince,
					ReceiverCity:      order.ReceiverCity,
					ReceiverRegion:    order.ReceiverRegion,
					BoxContent:        order.BoxContent,
					CouponIds:         datatypes.JSON([]byte(`[]`)),
				}).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		err := tx.Model(&models.AiMatchProductNftSecondModel{}).
			Where("order_id IN (?)", orderIds).
			Where("is_delete", 0).
			Where("status", "on_shelf").
			Update("status", "off_shelf").Error
		if err != nil {
			tx.Rollback()
			return err
		}
		payload := map[string]any{
			"recycle_count": gorm.Expr("recycle_count + ?", len(v)),
			"status":        1,
		}
		if k == len(orderChunk)-1 {
			payload["status"] = -1
		}

		err = tx.Model(&models.AiMatchProductNftRecycleRecord{}).
			Where("id", recycleRecord.ID).
			Updates(payload).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}
