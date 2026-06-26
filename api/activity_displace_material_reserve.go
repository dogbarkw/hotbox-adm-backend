package api

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/common"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"
	"hotbox-adm-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
)

// @Summary 获取置换、分解活动预留材料任务信息
// @Description 置换、分解活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetNftDisplaceMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/operation/nft/displace/material_reserve/query [get]
func GetNftDisplaceMaterialReserveTask(c *gin.Context) {
	req := form.GetNftDisplaceMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	activityId := req.DisplaceId
	materials, err := common.GetMaterials(req.ActivityType, activityId)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 获取预留材料表信息
	infos, err := models.ActivityMaterialReserveDal.GetByParams(c, map[string]any{
		"activity_id":   activityId,
		"activity_type": req.ActivityType,
		"exec_status":   0,
	})
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}
	var info models.AiMatchProductNftActivityMaterialReserve
	var details []models.AiMatchProductNftActivityMaterialReserveDetail
	// 如果材料预留表信息为0，则从材料里面组装出来
	if len(infos) == 0 {
		for _, v := range materials.MaterialsData.Materials {
			for _, vv := range v {
				newDetail := models.AiMatchProductNftActivityMaterialReserveDetail{
					MaterialId:    int64(vv.MaterialId),
					MaterialType:  vv.Type,
					MaterialName:  vv.Name,
					ProductId:     uint64(vv.ProductId),
					ProductSizeId: uint64(vv.NftProductSizeId),
					PropId:        uint64(vv.PropId),
				}
				details = append(details, newDetail)
			}
		}
	} else {
		info = infos[0]
		details, err = models.ActivityMaterialReserveDetailDal.GetByParams(c, map[string]any{
			"reserve_id": info.Id,
		})
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}

	eg := errgroup.Group{}
	lock := sync.Mutex{}
	// 查询每个材料的任务和剩余份数
	reserveDetails := make([]dto.DisplaceMaterialReserveMaterialInfo, 0, len(details))
	for _, v := range details {
		vv := v
		eg.Go(func() error {
			switch vv.MaterialType {
			case "nft", "product":
				// 获取藏品所有任务
				taskList, err := getProductAllReverseTask(c, vv.ProductId, vv.ProductSizeId)
				if err != nil {
					return err
				}
				// 所有任务(排除掉自己的)的预留份数总和
				totalReverseNum := int64(0)
				for _, task := range taskList {
					if task.ActivityId == activityId && task.ActivityType == int(req.ActivityType) && !task.IsCountReset {
						continue
					}
					totalReverseNum += task.ReserveNum
				}
				// 获取藏品剩余份数
				userCount, appCount, err := getProductRemainNumber(vv.ProductId, vv.ProductSizeId)
				if err != nil {
					return err
				}
				reserveDetail := dto.DisplaceMaterialReserveMaterialInfo{
					MaterialId:          int(vv.MaterialId),
					MaterialName:        vv.MaterialName,
					AppRemainNum:        appCount,
					UserRemainNum:       userCount,
					MaterialType:        vv.MaterialType,
					ProductId:           vv.ProductId,
					ProductSizeId:       vv.ProductSizeId,
					ReserveNum:          vv.ReserveNum,
					AvailableReserveNum: appCount - userCount - totalReverseNum,
					TaskList:            taskList,
				}
				lock.Lock()
				reserveDetails = append(reserveDetails, reserveDetail)
				lock.Unlock()
				return nil
			case "prop":
				totalPropNum, err := models.NftPropUserDal.UserTotalPropNum(c, vv.PropId)
				if err != nil {
					return err
				}
				reserveDetail := dto.DisplaceMaterialReserveMaterialInfo{
					MaterialId:          int(vv.MaterialId),
					MaterialName:        vv.MaterialName,
					MaterialType:        vv.MaterialType,
					PropId:              vv.PropId,
					ReserveNum:          vv.ReserveNum,
					AppRemainNum:        99999,
					UserRemainNum:       totalPropNum,
					AvailableReserveNum: 99999,
				}
				lock.Lock()
				reserveDetails = append(reserveDetails, reserveDetail)
				lock.Unlock()
				return nil
			default:
				return fmt.Errorf("[GetNftCombinationMaterialReserveTask] range details material_type invalid, material_uuid:%s", vv.MaterialUuid)
			}
		})
	}
	err = eg.Wait()
	if err != nil {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[GetNftCombinationMaterialReserveTask] errgroup.Group.Wait fail, [%s]comb_id:%d, err: %v", util.StringActivityType(req.ActivityType), activityId, err))
		response.ResponseFail(err.Error())
		return
	}
	sort.Slice(reserveDetails, func(i, j int) bool {
		return reserveDetails[i].MaterialId < reserveDetails[j].MaterialId
	})
	resp := dto.DisplaceMaterialReserveInfo{
		DisplaceId:         activityId,
		ActivityType:       int(req.ActivityType),
		ExecTime:           info.ExecTime,
		Remark:             info.Remark,
		ActivityReserveNum: info.ActivityReserveNum,
		Materials:          reserveDetails,
	}
	response.ResponseSuccess(resp)
}

// @Summary 新建/更新置换、分解活动预留材料任务信息
// @Description 置换、分解活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.UpdateNftDisplaceMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/operation/nft/displace/material_reserve/create [post]
func UpdateNftDisplaceMaterialReserveTask(c *gin.Context) {
	req := form.UpdateNftDisplaceMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	// 加锁
	lockKey := fmt.Sprintf("hd:comb_material_reserve:lock:%d_%d", req.DisplaceId, req.ActivityType)

	lock := cli.HotDogRedis.SetNX(c, lockKey, "lock", 30*time.Second).Val()
	if !lock {
		response.ResponseFail("操作过于频繁，请30s后再试")
		return
	}
	activityData, err := models.AiMatchProductNftReplaceDal.GetByReplaceId(c, int(req.DisplaceId))
	if err != nil {
		logrus.Errorf("%+v", err.Error())
		response.ResponseFail("获取活动信息失败")
		return
	}
	if activityData.StartTime <= time.Now().UnixMilli() {
		response.ResponseFail("活动已开始，不能进行修改")
		return
	}

	// 先检查各个材料总数是否够扣
	failMaterialNameList := make([]string, 0)
	mp := make(map[string]int)
	for _, v := range req.Materials {
		ok, err := checkNftDisplaceMaterialReserveTaskMaterial(c, mp, req.DisplaceId, req.ActivityType, v)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
		if !ok {
			failMaterialNameList = append(failMaterialNameList, v.MaterialName)
		}
	}

	if len(failMaterialNameList) != 0 {
		response.Responses(errno.CombinationMaterialReserveNumOverload, "预留材料数量超出库存值", dto.UpdateNftCombinationMaterialReserveTaskResp{FailMaterialNameList: failMaterialNameList})
		return
	}

	tx := cli.HotDogGormDB.Begin()
	now := time.Now()
	// 把旧数据软删除
	err = tx.Model(models.ActivityMaterialReserveDal).Where(map[string]any{
		"activity_id":   req.DisplaceId,
		"activity_type": req.ActivityType,
		"exec_status":   0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	err = tx.Model(models.ActivityMaterialReserveDetailDal).Where(map[string]any{
		"activity_id":   req.DisplaceId,
		"activity_type": req.ActivityType,
		"status":        0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	// 新增
	if req.ExecTime == 0 {
		req.ExecTime = time.Now().UnixMilli()
	}
	operatorId, _ := c.Get("user_id")
	newMainInfo := models.AiMatchProductNftActivityMaterialReserve{
		ActivityId:         req.DisplaceId,
		ActivityType:       int(req.ActivityType),
		ExecTime:           req.ExecTime,
		Remark:             req.Remark,
		ActivityReserveNum: req.ActivityReserveNum,
		UserId:             cast.ToInt64(operatorId),
		UserName:           c.GetString("adm_user_name"),
	}
	err = tx.Model(models.ActivityMaterialReserveDal).Create(&newMainInfo).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}
	newDetails := buildDisplaceMaterialReserveDetails(newMainInfo.Id, req.DisplaceId, req.ActivityType, req.DisplaceMaterialReserveInfo.Materials)
	err = tx.Model(models.ActivityMaterialReserveDetailDal).Create(newDetails).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		response.ResponseFail(err.Error())
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.ResponseFail(err.Error())
	}
	response.ResponseSuccess(nil)
}

// 查询每个材料是否够扣，countMap的value是材料的剩余可扣除数量
func checkNftDisplaceMaterialReserveTaskMaterial(c *gin.Context, countMap map[string]int, activityId, activityType int64, req dto.DisplaceMaterialReserveMaterialInfo) (ok bool, err error) {
	reserveMaterialTotalNum := int64(0) // 本次要预留的份数 + 其他任务要预留的份数
	if lo.Contains([]string{"nft", "product"}, req.MaterialType) {
		// 现有任务的预留数量
		taskNums, err := getCombMaterialReserveOtherTaskTotal(c, activityId, req.ProductId, req.ProductSizeId, int(activityType))
		if err != nil {
			return false, err
		}
		reserveMaterialTotalNum += taskNums
		// 传入的预留数量
		reserveMaterialTotalNum += req.ReserveNum

		// 查询可以扣的数量
		userCount, totalNum, err := getProductRemainNumber(req.ProductId, req.ProductSizeId)
		if err != nil {
			return false, err
		}
		// 第一次都不够扣，直接返回失败
		if (totalNum - userCount) < reserveMaterialTotalNum {
			return false, nil
		}

		key := fmt.Sprintf("%d_%d", req.ProductId, req.ProductSizeId)
		resNum, getOk := countMap[key]
		if !getOk {
			countMap[key] = int(totalNum-userCount) - int(taskNums) - int(req.ReserveNum)
		} else {
			// 已存在，则对比这次传入的剩余数量够不够扣
			if int64(resNum) < req.ReserveNum {
				return false, nil
			} else {
				// 足够扣
				countMap[key] = resNum - int(req.ReserveNum)
			}
		}

		return true, nil
	}
	return true, nil
}

func buildDisplaceMaterialReserveDetails(reserveId int64, activityId, activityType int64, materials []dto.DisplaceMaterialReserveMaterialInfo) (resp []models.AiMatchProductNftActivityMaterialReserveDetail) {
	for _, v := range materials {
		d := models.AiMatchProductNftActivityMaterialReserveDetail{
			ReserveId:     reserveId,
			ActivityId:    activityId,
			ActivityType:  int(activityType),
			MaterialId:    int64(v.MaterialId),
			MaterialType:  v.MaterialType,
			MaterialName:  v.MaterialName,
			ProductId:     v.ProductId,
			ProductSizeId: v.ProductSizeId,
			PropId:        v.PropId,
			ReserveNum:    v.ReserveNum,
		}
		resp = append(resp, d)
	}

	return resp
}

// @Summary 删除置换、分解预留材料任务信息
// @Description 置换、分解
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.DeleteNftDisplaceMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/operation/nft/combination/displace/delete [post]
func DeleteNftDisplaceMaterialReserveTask(c *gin.Context) {
	req := form.DeleteNftDisplaceMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	// 查询活动是否已开始，已开始则不能新增
	activityData, err := models.AiMatchProductNftReplaceDal.GetByReplaceId(c, int(req.DisplaceId))
	if err != nil {
		logrus.Errorf("%+v", err.Error())
		response.ResponseFail("获取活动信息失败")
		return
	}
	if activityData.StartTime <= time.Now().UnixMilli() {
		response.ResponseFail("活动已开始，不能进行修改")
		return
	}

	tx := cli.HotDogGormDB.Begin()
	now := time.Now()
	// 把旧数据软删除
	err = tx.Model(models.ActivityMaterialReserveDal).Where(map[string]any{
		"activity_id":   req.DisplaceId,
		"activity_type": req.ActivityType,
		"exec_status":   0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		response.ResponseFail(err.Error())
		return
	}

	err = tx.Model(models.ActivityMaterialReserveDetailDal).Where(map[string]any{
		"activity_id":   req.DisplaceId,
		"activity_type": req.ActivityType,
		"status":        0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		response.ResponseFail(err.Error())
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}

// @Summary 计算置换、分解活动预留材料任务信息合成次数
// @Description 置换、分解
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.CalculateNftDisplaceMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/operation/nft/displace/material_reserve/calculate [post]
func CalculateNftDisplaceMaterialReserveTask(c *gin.Context) {
	req := form.CalculateNftDisplaceMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	var (
		eg                 errgroup.Group
		lock               sync.Mutex
		materialArr        = []models.AiMatchProductNftReplaceContent{} // 材料数组
		totalDisplaceTimes = 0                                          // 传进来的材料总共能和多少次
	)

	// 查询活动信息
	activityInfo, err := models.AiMatchProductNftReplaceDal.GetByReplaceId(c, int(req.DisplaceId))
	if err != nil {
		response.ResponseFail("获取活动信息失败")
		return
	}

	// 根据材料组id查询材料组信息，得出最小合成数
	for _, v := range req.Materials {
		vv := v
		eg.Go(func() error {
			// 查询材料组信息
			materialInfo, err := models.AiMatchProductNftReplaceContentDal.One(c, vv.MaterialId)
			if err != nil {
				return err
			}
			// 只要有1个材料满足即可
			needNum := materialInfo.TargetCount
			reserveNum := vv.ReserveNum
			// 置换次数
			displaceTimes := int(util.SafeDivision(reserveNum, needNum))
			lock.Lock()
			materialArr = append(materialArr, materialInfo)
			totalDisplaceTimes += displaceTimes
			lock.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if totalDisplaceTimes == 0 {
		response.ResponseFail("没有计算出最小合成数")
		return
	}
	enough := false
	// 扣减每组材料的数量
	for i, info := range req.Materials {
		if enough {
			req.Materials[i].ExpectReserveNum = info.ReserveNum
			req.Materials[i].ReserveNum = 0
			continue
		}
		// 置换一次要这个材料多少个
		needNum := 0
		for _, v := range materialArr {
			if v.ReplaceContentId == int64(info.MaterialId) {
				needNum = int(v.TargetCount)
			}
		}
		displaceTimes := util.SafeDivision(info.ReserveNum, int64(needNum))
		req.Materials[i].ExpectReserveNum = info.ReserveNum
		req.Materials[i].ReserveNum = displaceTimes * int64(needNum)
		// 足够了
		if displaceTimes >= int64(totalDisplaceTimes) {
			enough = true
			continue
		}
	}
	if totalDisplaceTimes > int(activityInfo.TotalCount) {
		response.ResponseFail(fmt.Sprintf("超过%s次数", util.StringActivityType(req.ActivityType)))
		return
	}
	// 增加最小合成数字段
	req.ActivityReserveNum = int64(totalDisplaceTimes)
	response.ResponseSuccess(req)
}
