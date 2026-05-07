package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"
	"hotbox-adm-backend/util"

	"github.com/pkg/errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func GetNftCombinationList(c *gin.Context) {
	req := form.GetNftCombinationListReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	where := map[string]any{
		"combine_type": []string{"normal", "destroy", "group", "mix"},
		"is_delete":    0,
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SearchName != "" {
		where["nft_map_title like ?"] = "%" + req.SearchName + "%"
	}
	offset := (req.PageNumber - 1) * req.PageSize

	list, total, err := models.NewAiMatchProductNftCombination().GetAiMatchProductNftCombinationList(c, where, []string{"start_time desc"}, &req.PageSize, &offset)
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	if total == 0 {
		response.ResponseSuccessWithList(list, int(total))
		return
	}

	var nftSizeIds, boxProductIds, propIds []int64
	for _, v := range list {
		switch v.GenerationType {
		case "nft", "inscription":
			if v.NftSizeId == 0 {
				boxProductIds = append(boxProductIds, v.ProductId)
			} else {
				nftSizeIds = append(nftSizeIds, v.NftSizeId)
			}
		case "prop":
			propIds = append(propIds, v.PropId)
		}
	}
	var nftProductSize []models.SaleProductNftSizeModel
	var productSize []models.SaleCalendarProductSizeModel
	var propSize []models.AiMatchProductNftPropModel
	g, _ := errgroup.WithContext(c)
	g.Go(func() (err error) {
		nftProductSize, err = models.NewSaleProductNftSize().GetSaleProductNftSizeByParams(c, map[string]any{
			"nft_product_size_id": nftSizeIds,
		}, nil, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	g.Go(func() (err error) {
		productSize, err = models.NewSaleCalendarProductSize().GetSaleCalendarProductSizeByParams(c, map[string]any{
			"product_id": boxProductIds,
			"is_delete":  0,
		}, nil, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	g.Go(func() (err error) {
		propSize, err = models.NewAiMatchProductNftProp().GetSaleCalendarProductSizeByParams(c, map[string]any{
			"prop_id":   propIds,
			"is_delete": 0,
		}, nil, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		klog.Error(err)
		response.ResponseFail(errno.GetMsg(errno.Error))
		return
	}
	res := formatCombinationList(list, nftProductSize, productSize, propSize)
	response.ResponseSuccessWithList(res, int(total))
}

func formatCombinationList(list []models.AiMatchProductNftCombinationModel, nftProductSize []models.SaleProductNftSizeModel, productSize []models.SaleCalendarProductSizeModel, propSize []models.AiMatchProductNftPropModel) (result []dto.AiMatchProductNftCombinationRes) {
	type counters struct {
		StockCount int64
		TotalCount int64
	}
	propTypeMap := make(map[int64]string, 0)
	propSizeMap := make(map[int64]counters, 0)
	nftProductSizeMap := make(map[int64]counters, 0)
	productSizeMap := make(map[int64]counters, 0)
	for _, v := range propSize {
		propTypeMap[v.PropId] = v.BenefitType
		propSizeMap[v.PropId] = counters{
			StockCount: int64(v.StockCount),
			TotalCount: int64(v.TotalCount),
		}
	}
	for _, v := range nftProductSize {
		nftProductSizeMap[v.ID] = counters{
			StockCount: v.StockCount,
			TotalCount: v.TotalCount,
		}
	}
	for _, v := range productSize {
		productSizeMap[v.ProductId] = counters{
			StockCount: v.StockCount,
			TotalCount: v.TotalCount,
		}
	}
	for _, v := range list {
		publishTime := v.PublishTime.Format(time.DateTime)
		drawEndTime := v.DrawEndTime.Format(time.DateTime)
		if v.PublishTime.IsZero() {
			publishTime = "0000-00-00 00:00:00"
		}
		if v.DrawEndTime.IsZero() {
			drawEndTime = "0000-00-00 00:00:00"
		}

		nftMapTitleSubPic := []string{}
		_ = json.Unmarshal([]byte(v.NftMapTitleSubPic), &nftMapTitleSubPic)
		subProductInfo := []any{}
		_ = json.Unmarshal([]byte(v.SubProductInfo), &subProductInfo)
		subPropInfo := []any{}
		_ = json.Unmarshal([]byte(v.SubPropInfo), &subPropInfo)
		combineCode := []any{}
		_ = json.Unmarshal([]byte(v.CombineCode), &combineCode)
		// 初始化res
		res := dto.AiMatchProductNftCombinationRes{
			ID:                   int(v.ID),
			NftMapTitle:          v.NftMapTitle,
			NftMapTitleSub:       v.NftMapTitleSub,
			NftMapTitlePic:       v.NftMapTitlePic,
			NftMapTitleSubPic:    nftMapTitleSubPic,
			ProductID:            int(v.ProductId),
			NftSizeID:            int(v.NftSizeId),
			SubProductInfo:       subProductInfo,
			StartTime:            v.StartTime,
			EndTime:              v.EndTime,
			SecondSaleTime:       int(v.SecondSaleTime),
			Weight:               int(v.Weight),
			CreateTime:           v.CreateTime.Format(time.DateTime),
			UpdateTime:           v.UpdateTime.Format(time.DateTime),
			IsDelete:             int(v.IsDelete),
			OnSaleStatus:         int(v.OnSaleStatus),
			CombinationPicture:   v.CombinationPicture,
			CombineType:          v.CombineType,
			Version:              int(v.Version),
			SubPropInfo:          subPropInfo,
			SendScore:            int(v.SendScore),
			NewFeeRate:           int(v.NewFeeRate),
			NewEffectiveDays:     int(v.NewEffectiveDays),
			BoxOpenTime:          int(v.BoxOpenTime),
			Hot:                  int(v.Hot),
			IsLimitCombine:       cast.ToInt(v.IsLimitCombine),
			LimitNum:             int(v.LimitNum),
			CombineCode:          combineCode,
			GenerationType:       v.GenerationType,
			PropID:               int(v.PropId),
			SendPropID:           int(v.SendPropId),
			CombineMoreThanNum:   int(v.CombineMoreThanNum),
			CombineMoreThanRate:  int(v.CombineMoreThanRate),
			IsDisplayCount:       int(v.IsDisplayCount),
			IsDisplayTime:        int(v.IsDisplayTime),
			TotalTime:            int(v.TotalTime),
			TotalTimeValue:       int(v.TotalTimeValue),
			UserMaxTime:          int(v.UserMaxTime),
			UserMaxTimeValue:     int(v.UserMaxTimeValue),
			OriginTotalTimeValue: int(v.OriginTotalTimeValue),
			AppType:              int(v.AppType),
			RunRestStockStatus:   int(v.RunRestStockStatus),
			AdvanceReservation:   int(v.AdvanceReservation),
			ActivityType:         int(v.ActivityType),
			PublishTime:          publishTime,
			PublishFlag:          int(v.PublishFlag),
			DrawEndTime:          drawEndTime,
			RecomPrice:           int(v.RecomPrice),
			OncePrice:            int(v.OncePrice),
			IsNewMode:            int(v.IsNewMode),
			TemporaryReservation: int(v.TemporaryReservation),
			RemainNum:            int(v.RemainNum),
			SubscribeShow:        int(v.SubscribeShow),
		}
		res.HashID = util.EncodeInt64(v.ID)
		res.PropBenefitType = propTypeMap[v.PropId]
		res.CombineID = int(v.ID)
		if v.TotalTime == 2 {
			res.StockCount = res.TotalTimeValue
			res.TotalCount = int(v.OriginTotalTimeValue)
		} else {
			switch v.GenerationType {
			case "nft", "inscription":
				if v.IsLimitCombine {
					// 限制合成次数
					res.StockCount = int(v.LimitNum)

					res.TotalCount = len(combineCode)
				} else {
					if v.IsNewMode == 1 {
						res.StockCount = int(v.RemainNum)
						res.TotalCount = int(nftProductSizeMap[v.NftSizeId].TotalCount)
					} else {
						combineRestCountKey := fmt.Sprintf("combine_rest_nft_count:product:%d:nft_size_id:%d", v.ProductId, v.NftSizeId)
						restCount := cli.HotDogRedis.Get(context.Background(), combineRestCountKey).Val()
						if v.NftSizeId == 0 {
							// 合成盲盒
							if restCount != "" {
								res.StockCount = cast.ToInt(restCount)
							} else {
								res.StockCount = int(productSizeMap[v.ProductId].StockCount)
							}
							res.TotalCount = int(productSizeMap[v.ProductId].TotalCount)
						} else {
							// 合成正常藏品
							if restCount != "" {
								res.StockCount = cast.ToInt(restCount)
							} else {
								res.StockCount = int(nftProductSizeMap[v.NftSizeId].StockCount)
							}
							res.TotalCount = int(nftProductSizeMap[v.NftSizeId].TotalCount)
						}
					}
				}
			case "prop":
				if v.IsNewMode == 1 {
					res.StockCount = int(v.RemainNum)
					res.TotalCount = int(nftProductSizeMap[v.PropId].TotalCount)
				} else {
					combineRestCountKey := fmt.Sprintf("combine_rest_prop_count:prop:%d", v.PropId)
					restCount := cli.HotDogRedis.Get(context.Background(), combineRestCountKey).Val()
					if restCount != "" {
						res.StockCount = cast.ToInt(restCount)
					} else {
						res.StockCount = int(propSizeMap[v.PropId].StockCount)
					}
					res.TotalCount = int(propSizeMap[v.PropId].TotalCount)
				}
			}
		}
		result = append(result, res)
	}
	return
}
