package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// @Summary 获取合成活动预留材料任务信息
// @Description 合成活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetNftCombinationMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/collection/nft/combination/material_reserve/query [get]
func GetNftCombinationMaterialReserveTask(c *gin.Context) {
	req := form.GetNftCombinationMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	combId := req.CombinationId
	// 查询主信息
	infos, err := models.ActivityMaterialReserveDal.GetByParams(c, map[string]any{
		"activity_id":   combId,
		"activity_type": 1,
		"exec_status":   0,
	})
	if err != nil {
		logrus.Error(err)
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[GetNftCombinationMaterialReserveTask] ActivityMaterialReserveDal.GetByParams comb_id:%d, err: %v", combId, err))
		response.ResponseFail(err.Error())
		return
	}

	var info models.AiMatchProductNftActivityMaterialReserve
	var details []models.AiMatchProductNftActivityMaterialReserveDetail
	// 从合成活动材料组组装数据
	if len(infos) == 0 {
		details, err = getCombMaterialReserveInfoDetailFromGroup(c, combId)
		if err != nil {
			_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[GetNftCombinationMaterialReserveTask] getCombMaterialReserveInfoDetailFromGroup, comb_id:%d, err: %v", combId, err))
			response.ResponseFail(err.Error())
			return
		}
	} else {
		info = infos[0]
		details, err = models.ActivityMaterialReserveDetailDal.GetByParams(c, map[string]any{
			"reserve_id": info.Id,
		})
		if err != nil {
			_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[GetNftCombinationMaterialReserveTask] ActivityMaterialReserveDetailDal.GetByParams, comb_id:%d, err: %v", combId, err))
			response.ResponseFail(err.Error())
			return
		}
	}

	eg := errgroup.Group{}
	lock := sync.Mutex{}

	// 查询每个材料的任务和剩余份数
	reserveDetails := make([]dto.CombinationMaterialReserveMaterialInfo, 0, len(details))
	for _, v := range details {
		vv := v
		eg.Go(func() error {
			switch vv.MaterialType {
			case "nft", "product":
				// 获取藏品所有任务
				taskList, err := getProductAllReverseTask(vv.ProductId, vv.ProductSizeId)
				if err != nil {
					return err
				}
				// 所有任务(排除掉自己的)的预留份数总和
				totalReverseNum := int64(0)
				for _, task := range taskList {
					if task.ActivityId == combId && task.ActivityType == constant.ACTIVITY_TYPE_COMBINATION && !task.IsCountReset {
						continue
					}
					totalReverseNum += task.ReserveNum
				}
				// 获取藏品剩余份数
				userCount, appCount, err := getProductRemainNumber(vv.ProductId, vv.ProductSizeId)
				if err != nil {
					return err
				}

				reserveDetail := dto.CombinationMaterialReserveMaterialInfo{
					MaterialUuid:        vv.MaterialUuid,
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
				reserveDetail := dto.CombinationMaterialReserveMaterialInfo{
					MaterialUuid:        vv.MaterialUuid,
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
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[GetNftCombinationMaterialReserveTask] errgroup.Group.Wait fail, comb_id:%d, err: %v", combId, err))
		response.ResponseFail(err.Error())
		return
	}

	groupIdToDetailsMap := make(map[string]dto.CombinationMaterialReserveMaterialGroup)
	// 要根据材料组进行聚合
	for _, v := range reserveDetails {
		rd := v
		item := groupIdToDetailsMap[rd.MaterialUuid]
		item.MaterialUuid = rd.MaterialUuid
		item.MaterialInfos = append(item.MaterialInfos, rd)
		groupIdToDetailsMap[rd.MaterialUuid] = item
	}

	respItems := make([]dto.CombinationMaterialReserveMaterialGroup, 0, len(groupIdToDetailsMap))
	// 要根据材料组进行聚合
	for k := range groupIdToDetailsMap {
		respItems = append(respItems, groupIdToDetailsMap[k])
	}

	resp := dto.GetNftCombinationMaterialReserveTaskResp{
		CombinationMaterialReserveInfo: dto.CombinationMaterialReserveInfo{
			CombinationId:      combId,
			ExecTime:           info.ExecTime,
			Remark:             info.Remark,
			ActivityReserveNum: info.ActivityReserveNum,
			MaterialGroups:     respItems,
		},
	}
	response.ResponseSuccess(resp)
}

// 从活动材料组组装数据
func getCombMaterialReserveInfoDetailFromGroup(c context.Context, combId int64) (details []models.AiMatchProductNftActivityMaterialReserveDetail, err error) {
	materialGroups, err := models.CombinationMaterialDal.GetByParams(c, map[string]any{
		"combination_id": combId,
		"is_delete":      0,
	})
	if err != nil {
		return nil, err
	}
	for _, m := range materialGroups {
		// 序列化材料组的材料要求
		materialList := make([]map[string]interface{}, 0)
		err = json.Unmarshal([]byte(m.MaterialInfo), &materialList)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("[getCombMaterialReserveInfoDetailFromGroup] material_uuid:%s", m.MaterialUuid))
		}
		if len(materialList) == 0 {
			return nil, fmt.Errorf("[getCombMaterialReserveInfoDetailFromGroup] material_info empty, material_uuid:%s", m.MaterialUuid)
		}
		for _, val := range materialList {
			name := cast.ToString(val["name"])
			mType := cast.ToString(val["type"])
			productId := cast.ToUint64(val["sub_product_id"])
			nftProductSizeId := cast.ToUint64(val["nft_product_size_id"])
			propId := cast.ToUint64(val["prop_id"])

			newDetail := models.AiMatchProductNftActivityMaterialReserveDetail{
				MaterialUuid:  m.MaterialUuid,
				MaterialType:  mType,
				MaterialName:  name,
				ProductId:     productId,
				ProductSizeId: nftProductSizeId,
				PropId:        propId,
			}
			details = append(details, newDetail)
		}
	}
	return details, nil
}

// @Summary 新建/更新合成活动预留材料任务信息
// @Description 合成活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.UpdateNftCombinationMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/collection/nft/combination/material_reserve/create [post]
func UpdateNftCombinationMaterialReserveTask(c *gin.Context) {
	req := form.UpdateNftCombinationMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	// 加锁
	lockKey := fmt.Sprintf("hd:comb_material_reserve:lock:%d_%d", req.CombinationId, constant.ACTIVITY_TYPE_COMBINATION)

	lock := cli.HotDogRedis.SetNX(c, lockKey, "lock", 30*time.Second).Val()
	if !lock {
		response.ResponseFail("操作过于频繁，请30s后再试")
		return
	}

	// 查询活动是否已开始，已开始则不能新增
	combination, err := models.NftCombinationDal.One(c, req.CombinationId, "start_time")
	if err != nil {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] AiMatchProductNftCombination.One, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail("获取活动信息失败")
		return
	}
	if combination.StartTime <= time.Now().UnixMilli() {
		response.ResponseFail("活动已开始，不能进行修改")
		return
	}

	// 先检查各个材料总数是否够扣
	failMaterialNameList := make([]string, 0)
	mp := make(map[string]int)
	for _, v := range req.MaterialGroups {
		for _, vv := range v.MaterialInfos {
			ok, err := checkNftCombinationMaterialReserveTaskMaterial(mp, req.CombinationId, vv)
			if err != nil {
				_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] checkNftCombinationMaterialReserveTaskMaterial, comb_id:%d, err: %v", req.CombinationId, err))
				response.ResponseFail(err.Error())
				return
			}
			if !ok {
				failMaterialNameList = append(failMaterialNameList, vv.MaterialName)
			}
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
		"activity_id":   req.CombinationId,
		"activity_type": 1,
		"exec_status":   0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] ActivityMaterialReserveDal update is_deleted fail, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}

	err = tx.Model(models.ActivityMaterialReserveDetailDal).Where(map[string]any{
		"activity_id":   req.CombinationId,
		"activity_type": 1,
		"status":        0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] ActivityMaterialReserveDetailDal update is_deleted fail, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}

	// 新增
	if req.ExecTime == 0 {
		req.ExecTime = time.Now().UnixMilli()
	}
	operatorId, _ := c.Get("user_id")
	newMainInfo := models.AiMatchProductNftActivityMaterialReserve{
		ActivityId:         req.CombinationId,
		ActivityType:       1,
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
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] ActivityMaterialReserveDal create fail, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}
	newDetails := buildCombinationMaterialReserveDetails(newMainInfo.Id, req.CombinationId, req.MaterialGroups)
	err = tx.Model(models.ActivityMaterialReserveDetailDal).Create(newDetails).Error
	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] ActivityMaterialReserveDetailDal create fail, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] transation commit fail, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
	}
	response.ResponseSuccess(nil)
}

func buildCombinationMaterialReserveDetails(reserveId int64, activityId int64, groups []dto.CombinationMaterialReserveMaterialGroup) (resp []models.AiMatchProductNftActivityMaterialReserveDetail) {
	for _, v := range groups {
		for _, vv := range v.MaterialInfos {
			d := models.AiMatchProductNftActivityMaterialReserveDetail{
				ReserveId:     reserveId,
				ActivityId:    activityId,
				ActivityType:  1,
				MaterialUuid:  vv.MaterialUuid,
				MaterialType:  vv.MaterialType,
				MaterialName:  vv.MaterialName,
				ProductId:     vv.ProductId,
				ProductSizeId: vv.ProductSizeId,
				PropId:        vv.PropId,
				ReserveNum:    vv.ReserveNum,
			}
			resp = append(resp, d)
		}
	}
	return resp
}

// @Summary 删除合成活动预留材料任务信息
// @Description 合成活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.DeleteNftCombinationMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/collection/nft/combination/material_reserve/delete [post]
func DeleteNftCombinationMaterialReserveTask(c *gin.Context) {
	req := form.DeleteNftCombinationMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	// 查询活动是否已开始，已开始则不能新增
	combination, err := models.NftCombinationDal.One(c, req.CombinationId, "start_time")
	if err != nil {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[UpdateNftCombinationMaterialReserveTask] AiMatchProductNftCombination.One, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail("获取活动信息失败")
		return
	}
	if combination.StartTime <= time.Now().UnixMilli() {
		response.ResponseFail("活动已开始，不能进行修改")
		return
	}

	tx := cli.HotDogGormDB.Begin()
	now := time.Now()
	// 把旧数据软删除
	err = tx.Model(models.ActivityMaterialReserveDal).Where(map[string]any{
		"activity_id":   req.CombinationId,
		"activity_type": 1,
		"exec_status":   0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[DeleteNftCombinationMaterialReserveTask] ActivityMaterialReserveDal update deleted_at fail. combination_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}

	err = tx.Model(models.ActivityMaterialReserveDetailDal).Where(map[string]any{
		"activity_id":   req.CombinationId,
		"activity_type": 1,
		"status":        0,
	}).Updates(map[string]any{
		"deleted_at": now,
	}).Error
	if err != nil {
		tx.Rollback()
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[DeleteNftCombinationMaterialReserveTask] ActivityMaterialReserveDetailDal update deleted_at fail. combination_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[DeleteNftCombinationMaterialReserveTask] transaction commit fail. combination_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}

// 查询每个材料是否够扣，countMap的value是材料的剩余可扣除数量
func checkNftCombinationMaterialReserveTaskMaterial(countMap map[string]int, combId int64, req dto.CombinationMaterialReserveMaterialInfo) (ok bool, err error) {
	reserveMaterialTotalNum := int64(0) // 本次要预留的份数 + 其他任务要预留的份数
	if req.MaterialType == "nft" {
		// 现有任务的预留数量
		taskNums, err := getCombMaterialReserveOtherTaskTotal(combId, req.ProductId, req.ProductSizeId, constant.ACTIVITY_TYPE_COMBINATION)
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

// 查询这个活动之外已预留的藏品数量
func getCombMaterialReserveOtherTaskTotal(activityId int64, productId, productSizeId uint64, activityType int) (num int64, err error) {
	// 获取总数
	allTask, err := getProductAllReverseTask(productId, productSizeId)
	if err != nil {
		return 0, err
	}
	// 排除掉活动本身原有的数量
	for _, task := range allTask {
		if !task.IsCountReset && task.ActivityType == activityType && task.ActivityId == activityId {
			continue
		}
		num += task.ReserveNum
	}
	return num, nil
}

// @Summary 计算合成活动预留材料任务信息合成次数
// @Description 合成活动
// @Tags 材料预留
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.CalculateNftCombinationMaterialReserveTaskReq true "查询参数"
// @Success 200 {object} any
// @Router /aiera/v2/collection/nft/combination/material_reserve/calculate [post]
func CalculateNftCombinationMaterialReserveTask(c *gin.Context) {
	req := form.CalculateNftCombinationMaterialReserveTaskReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	var (
		eg           errgroup.Group
		lock         sync.Mutex
		minCombNum   = math.MaxInt32                                                // 初始化取最大值，然后取各组的最小值                                                  // 算出各个材料组的最小合成次数
		groupInfoMap = make(map[string]models.AiMatchProductNftCombinationMaterial) // 材料组
		needNumMap   = make(map[string]int)
	)

	// 查询活动信息
	combination, err := models.NftCombinationDal.One(c, req.CombinationId, "*")
	if err != nil {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[CalculateNftCombinationMaterialReserveTask] AiMatchProductNftCombination.One, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail("获取活动信息失败")
		return
	}
	// 多产物的暂不支持
	// 序列化材料组的材料要求
	if len(combination.SlaveId) != 0 {
		slaveIds := make([]int64, 0)
		err = json.Unmarshal([]byte(combination.SlaveId), &slaveIds)
		if err != nil {
			response.ResponseFail("活动信息多产物信息异常")
			return
		}
		if len(slaveIds) > 0 {
			response.ResponseFail("暂不支持多产物")
			return
		}
	}

	// 根据材料组id查询材料组信息，得出最小合成数
	for _, v := range req.MaterialGroups {
		g := v
		eg.Go(func() error {
			// 查询材料组信息
			group, err := models.CombinationMaterialDal.One(c, g.MaterialUuid)
			if err != nil {
				return errors.WithMessage(err, fmt.Sprintf("[CalculateNftCombinationMaterialReserveTask] CombinationMaterialDal.One fail, material_uuid:%s", g.MaterialUuid))
			}
			// 判断能合多少次
			// 0是总数，1是其中一个达到
			gType := group.MaterialType
			lock.Lock()
			groupInfoMap[g.MaterialUuid] = group
			lock.Unlock()

			if gType == 0 {
				combNeedNum := group.MaterialNum // 满足这个总数就可以合成一次
				mTotalNum := 0
				for _, mInfo := range g.MaterialInfos {
					mTotalNum += int(mInfo.ReserveNum)
				}
				// 合成次数
				combTimes := mTotalNum / int(combNeedNum)
				lock.Lock()
				if combTimes < minCombNum {
					minCombNum = combTimes
				}
				lock.Unlock()
			}
			if gType == 1 {
				// 序列化材料组的材料要求
				materialList := make([]map[string]interface{}, 0)
				err = json.Unmarshal([]byte(group.MaterialInfo), &materialList)
				if err != nil {
					return err
				}
				if len(materialList) == 0 {
					return errors.New("材料组数据错误")
				}

				for _, val := range materialList {
					mType := cast.ToString(val["type"])
					productId := cast.ToUint64(val["sub_product_id"])
					nftProductSizeId := cast.ToUint64(val["nft_product_size_id"])
					propId := cast.ToUint64(val["prop_id"])
					needNum := cast.ToInt(val["need_num"])

					if mType == "nft" {
						key := fmt.Sprintf("%s_nft_%d_%d", group.MaterialUuid, productId, nftProductSizeId)
						needNumMap[key] = needNum
					} else {
						key := fmt.Sprintf("%s_prop_%d", group.MaterialUuid, propId)
						needNumMap[key] = needNum
					}
				}

				totalTimes := 0
				for _, info := range g.MaterialInfos {
					needNum := 0
					key := ""
					if info.MaterialType == "nft" {
						key = fmt.Sprintf("%s_nft_%d_%d", info.MaterialUuid, info.ProductId, info.ProductSizeId)
					} else {
						key = fmt.Sprintf("%s_prop_%d", info.MaterialUuid, info.PropId)
					}
					needNum = needNumMap[key]
					singleTimes := int(info.ReserveNum) / needNum
					totalTimes += singleTimes
				}

				lock.Lock()
				if totalTimes < minCombNum {
					minCombNum = totalTimes
				}
				lock.Unlock()
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[CalculateNftCombinationMaterialReserveTask] checkNftCombinationMaterialReserveTaskMaterial, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail(err.Error())
		return
	}

	if minCombNum == math.MaxInt32 {
		_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[CalculateNftCombinationMaterialReserveTask] 没有计算出最小合成数, comb_id:%d, err: %v", req.CombinationId, err))
		response.ResponseFail("没有计算出最小合成数")
		return
	}

	// 扣减每组材料的数量
	for _, v := range req.MaterialGroups {
		groupMaterialType := groupInfoMap[v.MaterialUuid].MaterialType

		switch groupMaterialType {
		case 0:
			// 总数满足条件
			gNum := groupInfoMap[v.MaterialUuid].MaterialNum
			delTotal := int(gNum) * minCombNum // 这组内要扣的总数

			enough := false
			for i := range v.MaterialInfos {
				// 如果预加已经够了，剩下的全部置零
				if enough {
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
					v.MaterialInfos[i].ReserveNum = 0
					continue
				}

				// 足够扣
				if int(v.MaterialInfos[i].ReserveNum) >= delTotal {
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
					v.MaterialInfos[i].ReserveNum = int64(delTotal)
					enough = true
				} else {
					delTotal -= int(v.MaterialInfos[i].ReserveNum)
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
				}
			}
		case 1:
			// 或关系
			minTimes := minCombNum
			enough := false
			for i := range v.MaterialInfos {
				info := v.MaterialInfos[i]
				if enough {
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
					v.MaterialInfos[i].ReserveNum = 0
					continue
				}

				// 合成一次要这个材料多少个
				needNum := 0
				key := ""
				if info.MaterialType == "nft" {
					key = fmt.Sprintf("%s_nft_%d_%d", info.MaterialUuid, info.ProductId, info.ProductSizeId)
				} else {
					key = fmt.Sprintf("%s_prop_%d", info.MaterialUuid, info.PropId)
				}
				needNum = needNumMap[key]

				combTimes := info.ReserveNum / int64(needNum)

				// 足够了
				if combTimes >= int64(minTimes) {
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
					v.MaterialInfos[i].ReserveNum = int64(minTimes * needNum)
					enough = true
					continue
				} else {
					// 不够，直接取最大次数
					v.MaterialInfos[i].ExpectReserveNum = v.MaterialInfos[i].ReserveNum
					v.MaterialInfos[i].ReserveNum = combTimes * int64(needNum)
					minTimes -= int(combTimes)
				}

			}
		}
	}
	if int64(minCombNum)*combination.Cnt > combination.TotalTimeValue {
		response.ResponseFail("超过合成次数")
		return
	}

	// 增加最小合成数字段
	req.ActivityReserveNum = int64(minCombNum) * combination.Cnt
	response.ResponseSuccess(req)
}

func getProductAllReverseTask(productId, productSizeId uint64) (result []dto.MaterialReserveMaterialTaskInfo, err error) {
	// 查询预留任务
	list, err := models.ActivityMaterialReserveDetailDal.GetJoinMainWithParams(map[string]any{
		"product_id":      productId,
		"product_size_id": productSizeId,
		"exec_status":     0,
		"status":          0,
	})
	if err != nil {
		return nil, err
	}
	for _, dl := range list {
		result = append(result, dto.MaterialReserveMaterialTaskInfo{
			ActivityType: dl.ActivityType,
			ExecTime:     dl.ExecTime,
			ReserveNum:   dl.ReserveNum,
			ActivityId:   dl.ActivityId,
			IsCountReset: false,
		})
	}
	// 藏品数量定时变更
	req := form.SecondPriceListReq{
		ProductId:        int(productId),
		NftProductSizeId: int(productSizeId),
		OnSaleStatus:     -1,
	}
	priceList, _, err := models.Nft{Ctx: &gin.Context{}}.NftProductPriceList(req)
	if err != nil {
		return
	}
	if len(priceList) > 0 {
		item := priceList[0]
		execTime := cast.ToInt64(item.CountResetTime)
		if execTime > 0 {
			// 如果是定时更新库存， 活动id、活动类型都为0
			result = append(result, dto.MaterialReserveMaterialTaskInfo{
				ExecTime:     execTime,
				ReserveNum:   int64(item.CountResetValue),
				IsCountReset: true,
			})
		}
	}
	return result, nil
}

// 获取藏品持有数量 realUserSurplus用户持有总量；nftCount藏品总量
func getProductRemainNumber(productId, productSizeId uint64) (realUserSurplus, nftCount int64, err error) {
	c := &gin.Context{}
	realUserSurplus, nftCount, err = models.Nft{Ctx: c}.GetNftUserOwnNumAndAppShowNum(int(productId), int(productSizeId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, errors.WithMessage(err, fmt.Sprintf("查询藏品APP展示剩余份数和用户拥有该藏品数失败，product_id:%d, product_size_id:%d", productId, productSizeId))
	}
	return
}
