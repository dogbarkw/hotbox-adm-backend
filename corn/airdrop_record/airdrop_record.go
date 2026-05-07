package airdrop_record

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const AIRDROP_RECORD_LOCK_KEY = "hd:airdrop_record_lock"

var SourceTypeMap = map[string]string{
	"PFP_video": "PFP_video",
	"PFP_photo": "PFP_photo",
	"SYG_photo": "SYG_photo",
}

type AirdropRecordCornJob struct{}

func (p *AirdropRecordCornJob) Run() {
	ctx := &gin.Context{}
	lock := cli.HotDogRedis.SetNX(ctx, "hd:nft_airdrop_record_job:lock", "lock", 5*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, "hd:nft_airdrop_record_job:lock")
	StartAirdropRecordCornJob(ctx)
}

func StartAirdropRecordCornJob(ctx *gin.Context) {
	msgType := "批量空投"
	allAirdropRecordList, err := models.AiMatchProductNftRecycleRecord{}.GetByParams(ctx, map[string]any{
		"type":   constant.BATCH_RECORD_AIRDROP_TYPE,
		"status": 0,
	}, 5)
	if err != nil {
		logrus.Errorf("StartAirdropRecordCornJob 查询回收记录表报错:%s", err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("查询回收记录表报错:%s", err))

		return
	}
	g, _ := errgroup.WithContext(ctx)
	for _, ar := range allAirdropRecordList {
		airdropRecord := ar
		g.Go(func() error {
			// 对藏品枷锁
			key := fmt.Sprintf("%s_%d_%d", AIRDROP_RECORD_LOCK_KEY, airdropRecord.ProductId, airdropRecord.NftProductSizeId)
			lock := cli.HotDogRedis.SetNX(ctx, key, "lock", 15*time.Minute).Val()
			if !lock {
				return nil
			}
			defer cli.HotDogRedis.Del(ctx, key)
			// 避免重启时重复获取，这里用的是目标数-已完成数来进行查询
			targetCount := airdropRecord.RecycleTargetCount - airdropRecord.RecycleCount
			nftLimit := 10
			successCounter := 0
			// 每次只能获取 10条，不够就继续
			for successCounter < targetCount {
				minLimit := lo.Min([]int{targetCount - successCounter, nftLimit})
				if minLimit <= 0 {
					return nil
				}
				nftIdsData, err := httpReq.NftRestPop(os.Getenv("NO_EXPIRE_USER_TOKEN"), dto.NftRestPopSendReq{
					ProductId:     airdropRecord.ProductId,
					ProductSizeId: airdropRecord.NftProductSizeId,
					Limit:         minLimit,
				})
				if err != nil {
					httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("id:%d,请求nftRestPop接口报错:%s", airdropRecord.ID, err))
					logrus.Error(err.Error())
					models.AiMatchProductNftRecycleRecord{}.Update(ctx, int(airdropRecord.ID), map[string]any{
						"msg":    err.Error(),
						"status": -1,
					})
					return err
				}
				if len(nftIdsData.Data) == 0 {
					logrus.Error(nftIdsData.Msg)
					httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("id:%d,%s", airdropRecord.ID, nftIdsData.Msg))
					models.AiMatchProductNftRecycleRecord{}.Update(ctx, int(airdropRecord.ID), map[string]any{
						"msg":    nftIdsData.Msg,
						"status": -1,
					})
					return err
				}
				nftData, err := models.AiMatchRestNftProduct{Ctx: ctx}.GetRestNftProductWithParams(map[string]any{
					"id": nftIdsData.Data,
				})
				if err != nil {
					httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("id:%d,获取nftRestNftProduct数据失败，id:%v:%s", airdropRecord.ID, nftIdsData.Data, err))
					return err
				}
				if len(nftData) == 0 {
					httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("id:%d,获取nftRestNftProduct数据为空，id:%v:%s", airdropRecord.ID, nftIdsData.Data, err))
					return nil
				}
				// 处理空投，方法就已经将成功计数加上
				err = handleAirdropOrder(ctx, &successCounter, targetCount, nftData, airdropRecord)
				if err != nil {
					httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("id:%d,处理空投失败，id:%v:%s", airdropRecord.ID, nftIdsData.Data, err))
					return err
				}

				if successCounter >= targetCount {
					break
				}

			}
			if successCounter > 0 {
				// 发送飞书推送
				userData, _ := models.User{Ctx: ctx}.FindSysUserById(uint64(airdropRecord.UserId))

				envStr := "测试环境"
				if strings.ToUpper(os.Getenv("ENV")) == "PRODUCTION" {
					envStr = "生产环境"
				}
				msg := fmt.Sprintf("[%s] %s(%d)批量空投（目标：%d）给%s(mobile:%s)如下藏品: %s %d个", envStr, airdropRecord.AdmUserName,
					airdropRecord.OperatorId, targetCount, userData.RealName, userData.Mobile, airdropRecord.ProductTitle, successCounter)
				httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_AIRDROP_URL, msgType, msg)
				p := struct {
					Title        string      `json:"title"`
					Nums         int         `json:"nums"`
					UserInfoList [][2]string `json:"userInfoList"`
				}{
					Title: airdropRecord.ProductTitle,
					Nums:  targetCount,
					UserInfoList: [][2]string{
						{cast.ToString(userData.UserId), userData.Mobile},
					},
				}
				jsonStr, _ := json.Marshal(p)
				models.OperateRecord{Ctx: ctx}.CreateRecord(models.AiMatchBackendOperateRecord{
					UserId:      int64(airdropRecord.OperatorId),
					Username:    airdropRecord.AdmUserName,
					Remark:      "批量空投-任务",
					Scenes:      10001,
					AssociateId: int64(airdropRecord.ProductId),
					RequestData: string(jsonStr),
				})
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		logrus.Error(err)
	}
}

func handleAirdropOrder(ctx *gin.Context, successCounter *int, targetCount int, nftData []models.AiMatchRestNftProductModel, airdropRecord models.AiMatchProductNftRecycleRecord) (err error) {
	for _, v := range nftData {
		tx := cli.HotDogGormDB.WithContext(ctx).Begin()
		// 先更新，锁住数据
		updateQuery := tx.Model(&models.AiMatchRestNftProductModel{}).
			Where("id = ? AND is_release = ? AND is_delete = ?", v.ID, 0, 0).
			Updates(map[string]any{
				"is_release":   1,
				"operator_id":  airdropRecord.OperatorId,
				"release_time": time.Now(),
			})
		if updateQuery.Error != nil {
			tx.Rollback()
			return
		}
		if updateQuery.RowsAffected == 0 {
			tx.Rollback()
			continue
		}
		// 优惠卷
		coupons := []models.CouponBusiness{}
		if len(v.CouponIds) > 0 {
			couponIds := []int{}
			json.Unmarshal(v.CouponIds, &couponIds)
			if len(couponIds) > 0 {
				coupons, err = models.CouponBusinessDal.GetByParams(ctx, map[string]any{
					"id": couponIds,
				})
			}
		}
		boxCode := ""
		// 盲盒
		if v.NftProductSizeId == 0 {
			if len(v.ReceiverProvince) >= 4 {
				boxCode = fmt.Sprintf("#%s#%s/%s", v.ReceiverProvince[len(v.ReceiverProvince)-4:], v.ReceiverCity, v.ReceiverRegion)
			}
		}
		nftSize, err := models.NewSaleProductNftSize().GetOneByParams(ctx, map[string]any{
			"nft_product_size_id": v.IosSourceFile,
		})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		err = sendNftProductToUser(tx, sendNftProductToUserStruct{
			UserId:         airdropRecord.UserId,
			RestNftProduct: v,
			BoxCode:        boxCode,
			CouPons:        coupons,
			Note:           "ADMIN_DROP",
			LogisticsInfo:  v.CombineInfo,
			PlayMethod:     v.PlayMethod,
			NewFlashId:     v.CombineActiveId,
			SourceType:     SourceTypeMap[nftSize.SourceType],
		})
		if err != nil {
			tx.Rollback()
			return err
		}
		*successCounter++
		// 更新回收记录表
		updatePayload := map[string]any{
			"recycle_count": gorm.Expr("recycle_count + ?", 1),
			"status":        1,
		}
		if targetCount <= *successCounter {
			updatePayload["status"] = -1
		}
		err = tx.Model(&models.AiMatchProductNftRecycleRecord{}).
			Where("id", airdropRecord.ID).
			Updates(updatePayload).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if targetCount <= *successCounter {
			return nil
		}
	}
	return nil
}

type sendNftProductToUserStruct struct {
	UserId         int
	RestNftProduct models.AiMatchRestNftProductModel
	BoxCode        string
	CouPons        []models.CouponBusiness
	Note           string
	PlayMethod     string
	LogisticsInfo  string
	NewFlashId     int
	SourceType     string
}

// 返回错误，tx不进行rollback，需要调用方处理
func sendNftProductToUser(tx *gorm.DB, payload sendNftProductToUserStruct) error {
	userId := payload.UserId
	receiverCity := payload.RestNftProduct.ReceiverCity
	productId := payload.RestNftProduct.ProductId
	sizeId := payload.RestNftProduct.SizeId
	iosSourceFile := payload.RestNftProduct.IosSourceFile
	androidSourceFile := payload.RestNftProduct.AndroidSourceFile
	receiverProvince := payload.RestNftProduct.ReceiverProvince
	receiverRegion := payload.RestNftProduct.ReceiverRegion
	productTitle := payload.RestNftProduct.ProductTitle
	nftProductSizeId := payload.RestNftProduct.NftProductSizeId
	activePicture := payload.RestNftProduct.ActivePicture
	boxContent := payload.RestNftProduct.BoxContent
	boxCode := payload.BoxCode
	note := payload.Note
	logisticsInfo := payload.LogisticsInfo
	playMethod := payload.PlayMethod
	newFlashId := payload.NewFlashId
	sourceType := payload.SourceType
	receiverDetailAddress := util.GenerateRandomString(64)

	productInfo := models.SaleCalendarProductModel{}
	err := tx.Model(&models.SaleCalendarProductModel{}).Where("id", productId).First(&productInfo).Error
	if err != nil {
		return err
	}
	order := models.AiMatchProductOrderModel{
		UserId:                int64(userId),
		ProductId:             productId,
		SizeId:                sizeId,
		Size:                  "",
		Note:                  note,
		DeliverySn:            iosSourceFile,
		DeliveryStatus:        1,
		AndroidNftSource:      androidSourceFile,
		ReceiverProvince:      receiverProvince,
		ReceiverCity:          receiverCity,
		ReceiverName:          cast.ToString(userId),
		ReceiveTime:           time.Now(),
		PayAmount:             0,
		ReceiverDetailAddress: receiverDetailAddress,
		ReceiverRegion:        receiverRegion,
		ProductTitle:          productTitle,
		NftProductSizeId:      nftProductSizeId,
		ProductPicture:        activePicture,
		LogisticsInfo:         logisticsInfo,
		Status:                2,
		ProductType:           "NFT",
		Version:               2,
		BuyCount:              1,
		OrderNo:               "HC" + util.GenerateRandomString(21),
		BoxCode:               boxCode,
		BoxContent:            boxContent,
		PlayMethod:            playMethod,
		NewFlashId:            newFlashId,
		SourceType:            sourceType,
		MarketType:            productInfo.MarketType,
		OrderSn:               nil,
		CreateTime:            time.Now(),
		UpdateTime:            time.Now(),
	}
	err = tx.Model(&models.AiMatchProductOrderModel{}).Create(&order).Error
	if err != nil {
		return err
	}
	// 创建nft扭转信息
	nftHolder := models.NftTorsion{
		ProductId:   order.ProductId,
		OrderId:     order.ID,
		FromUser:    0,
		ToUser:      order.UserId,
		OrderNo:     order.OrderNo,
		TradingHash: order.ReceiverDetailAddress,
		Status:      "drop",
		CreateTime:  order.CreateTime,
	}
	err = tx.Model(&models.NftTorsion{}).Create(&nftHolder).Error
	if err != nil {
		return err
	}
	return nil
}
