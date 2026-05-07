package common

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// 获取活动产物、艺术家
func GetOutPutProductAndArtist(activityType int, activityData any) (outPutProduct []dto.OutPutProduct, artist []dto.Artist, err error) {
	ctx := &gin.Context{}
	switch activityType {
	// 合成
	case constant.ACTIVITY_TYPE_COMBINATION:
		activityData, ok := activityData.(models.AiMatchProductNftCombinationModel)
		if !ok {
			return nil, nil, errors.New("不合法数据类型")
		}
		combinationIds := []int64{activityData.ID}
		if !lo.Contains([]string{"null", "0", "[]"}, activityData.SlaveId.String()) {
			slaveIds := make([]int64, 0)
			err := json.Unmarshal([]byte(activityData.SlaveId), &slaveIds)
			if err != nil {
				return nil, nil, err
			}
			combinationIds = append(combinationIds, slaveIds...)
		}
		combinationList, err := models.NewAiMatchProductNftCombination().GetNftCombinationListByParams(ctx, map[string][]any{
			"id":        {combinationIds},
			"is_delete": {0},
		})
		if err != nil {
			return nil, nil, err
		}
		productIds := []int{}
		for _, v := range combinationList {
			outPutProductData := dto.OutPutProduct{
				ProductId:        int(v.ProductId),
				NftProductSizeId: int(v.NftSizeId),
				PropId:           int(v.PropId),
				Num:              int(v.Cnt),
			}
			if v.ProductId > 0 {
				productIds = append(productIds, int(v.ProductId))
				outPutProductData.Type = "product"
				// 盲盒
				if v.NftSizeId == 0 {
					data, err := models.SaleCalendarProduct{Ctx: ctx}.One(v.ProductId)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return nil, nil, err
					}
					outPutProductData.Name = data.ProductName
				} else {
					data, err := models.NewSaleProductNftSize().GetOneByParams(ctx, map[string]any{
						"product_id":          v.ProductId,
						"nft_product_size_id": v.NftSizeId,
					})
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return nil, nil, err
					}
					outPutProductData.Name = data.ProductTitle
				}
			}
			if v.PropId > 0 {
				outPutProductData.Type = "prop"
				// 获取道具类型
				propData, err := models.NewAiMatchProductNftProp().One(ctx, cast.ToUint64(v.PropId), "")
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil, err
				}
				outPutProductData.Name = propData.ProductTitle
			}
			outPutProduct = append(outPutProduct, outPutProductData)
		}
		if len(productIds) > 0 {
			artistData, err := models.SaleCalendarProduct{Ctx: ctx}.GetListJoinWithNftArtist(map[string]any{
				"sale_calendar_product.id": productIds,
			})
			if err != nil {
				return nil, nil, err
			}
			for _, v := range artistData {
				artistData := dto.Artist{
					ArtistId:   int(v.ArtistId),
					ArtistName: v.ArtistName,
				}
				artist = append(artist, artistData)
				for kk, vv := range outPutProduct {
					if vv.ProductId == int(v.Id) {
						outPutProduct[kk].ArtistData = artistData
					}
				}
			}
		}
		// 置换、分解
	case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
		activityData, ok := activityData.(models.AiMatchProductNftReplaceModel)
		if !ok {
			return nil, nil, errors.New("不合法数据类型")
		}
		displaceMaterialList, err := models.AiMatchProductNftReplaceContentDal.GetByParams(ctx, map[string]any{
			"replace_type": "to",
			"replace_id":   activityData.ReplaceId,
			"is_delete":    0,
		})
		if err != nil {
			return nil, nil, err
		}
		productIds := []int64{}
		for _, v := range displaceMaterialList {
			outPutProductData := dto.OutPutProduct{
				Name:             v.TargetName,
				Type:             v.TargetType,
				ProductId:        int(v.ProductId),
				NftProductSizeId: int(v.NftProductSizeId),
				PropId:           int(v.PropId),
				Num:              int(v.TargetCount),
			}
			productIds = append(productIds, v.ProductId)
			outPutProduct = append(outPutProduct, outPutProductData)
		}
		if len(productIds) > 0 {
			artistData, err := models.SaleCalendarProduct{Ctx: ctx}.GetListJoinWithNftArtist(map[string]any{
				"sale_calendar_product.id": productIds,
			})
			if err != nil {
				return nil, nil, err
			}
			for _, v := range artistData {
				artistData := dto.Artist{
					ArtistId:   int(v.ArtistId),
					ArtistName: v.ArtistName,
				}
				artist = append(artist, artistData)
				for kk, vv := range outPutProduct {
					if vv.ProductId == int(v.Id) {
						outPutProduct[kk].ArtistData = artistData
					}
				}
			}
		}
		// 升级
	// case constant.ACTIVITY_TYPE_UPGRADE:
	default:
		return nil, nil, errors.New("不合法活动类型")
	}
	return outPutProduct, artist, nil
}

// 获取活动材料
func GetMaterials[T int | int64](activityType, activityId T) (result dto.ActivityScoreMaterials, data error) {
	ctx := context.Background()
	result.ActivityType = int(activityType)
	result.ActivityId = int(activityId)
	switch activityType {
	// 合成
	case constant.ACTIVITY_TYPE_COMBINATION:
		materialList, err := models.CombinationMaterialDal.GetByParams(ctx, map[string]any{
			"combination_id": activityId,
			"is_delete":      0,
		})
		if err != nil {
			return result, err
		}
		if len(materialList) == 0 {
			return result, nil
		}
		for _, v := range materialList {
			mt := []dto.Materials{}
			mp := []map[string]any{}
			json.Unmarshal(v.MaterialInfo, &mp)
			for _, mm := range mp {
				material := dto.Materials{
					MaterialUuid:     v.MaterialUuid,
					Name:             cast.ToString(mm["name"]),
					ProductId:        cast.ToInt(mm["sub_product_id"]),
					NftProductSizeId: cast.ToInt(mm["nft_product_size_id"]),
					Num:              cast.ToInt(mm["need_num"]),
					PropId:           cast.ToInt(mm["prop_id"]),
					MaterialNum:      v.MaterialNum,
					MaterialType:     v.MaterialType,
				}
				if material.ProductId > 0 {
					material.Type = "product"
				}
				if material.PropId > 0 {
					material.Type = "prop"
				}
				mt = append(mt, material)
			}
			result.MaterialsData.Materials = append(result.MaterialsData.Materials, mt)
		}

		// 置换、分解
	case constant.ACTIVITY_TYPE_REPLACE, constant.ACTIVITY_TYPE_DECOMPOSE:
		materialList, err := models.AiMatchProductNftReplaceContentDal.GetByParams(ctx, map[string]any{
			"replace_type": "from",
			"replace_id":   activityId,
			"is_delete":    0,
		})
		if err != nil {
			return result, err
		}
		mt := []dto.Materials{}
		// 置换、分解的多组材料是或的关系
		for _, v := range materialList {
			material := dto.Materials{
				MaterialId:       int(v.ReplaceContentId),
				Name:             v.TargetName,
				ProductId:        int(v.ProductId),
				NftProductSizeId: int(v.NftProductSizeId),
				PropId:           int(v.PropId),
				Num:              int(v.TargetCount),
			}
			if material.ProductId > 0 {
				material.Type = "product"
			}
			if material.PropId > 0 {
				material.Type = "prop"
			}
			mt = append(mt, material)
		}
		result.MaterialsData.Materials = append(result.MaterialsData.Materials, mt)
		// 升级
	// case constant.ACTIVITY_TYPE_UPGRADE:

	default:
		return result, errors.New("不合法活动类型")
	}
	return result, nil
}
