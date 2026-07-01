package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/internal/service"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"
	"hotbox-adm-backend/util"

	"github.com/shopspring/decimal"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @Summary GMV目标列表
// @Description GMV目标列表
// @Tags GMV目标
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} any
// @Router /target_gmv/list [post]
func TargetGmvList(c *gin.Context) {
	response := until.NewResponse(c)
	partitionData, err := models.PartitionDataDal.GetByParams(c, map[string]any{"is_delete": 0, "main_id > ?": 0}, []string{"main_id asc", "child_id asc"})
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	partitionIds := make([][]any, 0, len(partitionData))
	for _, v := range partitionData {
		partitionIds = append(partitionIds, []any{v.MainId, v.ChildId})
	}
	dataMap, err := hd_task_models.HdPartitionDailyGmvStatDal.GetPartitionTodayGmvStatList(c, partitionIds)
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	startWeek, endWeek := util.GetWeekRange()
	weekStatMap, err := hd_task_models.HdPartitionDailyGmvStatDal.SumWeeklyGmvStat(c, cast.ToInt(startWeek.Format(util.DefaultFormat)), cast.ToInt(endWeek.Format(util.DefaultFormat)))
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	// 获取全局量化配比（不再按分区）
	globalRatio, err := hd_task_models.HdPartitionGmvQuantRatioDal.GetByPartitionId(c, 0, 0)
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	type dataItem struct {
		MainId     int64  `json:"main_id"`
		ChildId    int64  `json:"child_id"`
		MainTitle  string `json:"main_title"`
		ChildTitle string `json:"child_title"`

		TargetGmv    float64 `json:"target_gmv"`
		CurrentGmv   float64 `json:"current_gmv"`
		ShareIncome  float64 `json:"share_income"`
		WeeklyGmv    float64 `json:"weekly_gmv"`
		WeeklyIncome float64 `json:"weekly_income"`
		PreGmv       float64 `json:"pre_gmv"`
		QuantRatio   float64 `json:"quant_ratio"`
		Status       int     `json:"status"`
		TargetId     int64   `json:"target_id"`
	}
	dataList := make([]dataItem, 0, len(partitionData))
	for _, datum := range partitionData {
		data := dataItem{
			MainId:     datum.MainId,
			ChildId:    datum.ChildId,
			MainTitle:  datum.MainTitle,
			ChildTitle: datum.ChildTitle,
		}
		dailyGmvStat, ok := dataMap[fmt.Sprintf("%d-%d", datum.MainId, datum.ChildId)]
		if ok {
			data.TargetId = dailyGmvStat.Id
			data.TargetGmv = dailyGmvStat.TargetGmv
			data.CurrentGmv = dailyGmvStat.CurrentGmv + dailyGmvStat.PreGmv
			data.ShareIncome = dailyGmvStat.ShareIncome
			data.PreGmv = dailyGmvStat.PreGmv
			data.Status = dailyGmvStat.Status

			weekStat, ok1 := weekStatMap[fmt.Sprintf("%d-%d", datum.MainId, datum.ChildId)]
			if ok1 {
				data.WeeklyGmv = weekStat.WeeklyGmv
				data.WeeklyIncome = weekStat.WeeklyShareIncome
			}
		}
		quantRatio, ok2 := quantRatioMap[fmt.Sprintf("%d-%d", datum.MainId, datum.ChildId)]
		if ok2 {
			data.QuantRatio = quantRatio.QuantRatio
		}
		dataList = append(dataList, data)
	}
	response.ResponseSuccess(dataList)
}

// @Summary 获取全局量化配比
// @Description 获取全局量化配比
// @Tags GMV目标
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} any
// @Router /target_gmv/quant_ratio/info [post]
func GetTargetGmvQuantRatio(c *gin.Context) {
	response := until.NewResponse(c)
	quantRatio, err := hd_task_models.HdPartitionGmvQuantRatioDal.GetByPartitionId(c, 0, 0)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(gin.H{
		"quant_ratio": quantRatio.QuantRatio,
	})
}

// @Summary 修改GMV目标
// @Description 修改GMV目标
// @Tags GMV目标
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.HdTargetGmvUpdateReq true "查询参数"
// @Success 200 {object} any
// @Router /target_gmv/update [post]
func UpdateTargetGmv(c *gin.Context) {
	req := form.HdTargetGmvUpdateReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	if req.ChildId > 0 {
		// 检查大区是否配置了gmv
		gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(c, map[string]any{
			"date":     time.Now().Format(util.DefaultFormat),
			"main_id":  req.MainId,
			"child_id": 0,
		}, nil, nil)
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
		if len(gmvStats) > 0 {
			response.ResponseFail("大区已配置目标gmv")
			return
		}
	} else {
		// 检查子分区是否配置了gmv
		gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(c, map[string]any{
			"date":         time.Now().Format(util.DefaultFormat),
			"main_id":      req.MainId,
			"child_id > ?": 0,
		}, nil, nil)
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
		if len(gmvStats) > 0 {
			response.ResponseFail("子分区已配置目标gmv")
			return
		}
	}

	// 更新到adm操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err := models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "设置分区目标gmv",
		Scenes:      constant.OPERATE_PARTITION_TARGET_GMV,
		AssociateId: 0,
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}

	// 检查参数
	if req.TargetGmv <= 0 { // gmv<=0视为删除
		err := hd_task_models.HdPartitionDailyGmvStatDal.DeleteByPartitionId(c, req.MainId, req.ChildId)
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
		cli.HotDogRedis.HDel(c, constant.REDIS_HD_PARTITION_TARGET_GMV_SWITCH_KEY, fmt.Sprintf("%d-%d", req.MainId, req.ChildId))

		// 更新分区用户的补偿状态和分成比例
		err = service.HdPartitionGmvService.SwitchPartitionTestUserCompensateStatus(c, req.MainId, req.ChildId, constant.COMPENSATE_STATUS_CLOSE)
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
		response.ResponseSuccess(nil)
		return
	}

	partitionDailyGmvdata, err := hd_task_models.HdPartitionDailyGmvStatDal.GetOneByParams(c, map[string]any{
		"main_id":  req.MainId,
		"child_id": req.ChildId,
		"date":     time.Now().Format(util.DefaultFormat),
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	// 创建
	if partitionDailyGmvdata.Id == 0 {
		err = hd_task_models.HdPartitionDailyGmvStatDal.Create(c, &hd_task_models.HdPartitionDailyGmvStat{
			MainId:    req.MainId,
			ChildId:   req.ChildId,
			Date:      cast.ToInt(time.Now().Format(util.DefaultFormat)),
			TargetGmv: decimal.NewFromInt(req.TargetGmv).InexactFloat64(),
			Status:    1,
		})
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}

		// 更新分区用户的补偿状态和分成比例
		err = service.HdPartitionGmvService.SwitchPartitionTestUserCompensateStatus(c, req.MainId, req.ChildId, constant.COMPENSATE_STATUS_OPEN)
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
	} else {
		// 更新
		err = hd_task_models.HdPartitionDailyGmvStatDal.UpdateByParams(c, map[string]any{"id": partitionDailyGmvdata.Id}, map[string]any{
			"target_gmv": req.TargetGmv,
		})
		if err != nil {
			logrus.Error(err)
			response.ResponseFail(err.Error())
			return
		}
	}

	response.ResponseSuccess(nil)
}

// @Summary 账号补偿，暂停/恢复
// @Description 账号补偿，暂停/恢复
// @Tags GMV目标
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.HdTargetGmvSwitchReq true "查询参数"
// @Success 200 {object} any
// @Router /target_gmv/switch [post]
func TargetGmvSwitch(c *gin.Context) {
	req := form.HdTargetGmvSwitchReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	gmvStat, err := hd_task_models.HdPartitionDailyGmvStatDal.GetOneByParams(c, map[string]any{"id": req.TargetId})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseFail("id 不存在")
		return
	}
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	remark := "操作补偿分区目标gmv"
	if req.Status == 1 {
		remark = remark + "恢复"
	} else {
		remark = remark + "暂停"
	}

	// 更新到adm操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      remark,
		Scenes:      constant.OPERATE_PARTITION_TARGET_GMV,
		AssociateId: req.TargetId,
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}

	// 更新
	err = hd_task_models.HdPartitionDailyGmvStatDal.UpdateByParams(c, map[string]any{"id": req.TargetId}, map[string]any{
		"status": req.Status,
	})
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	cli.HotDogRedis.HSet(c, constant.REDIS_HD_PARTITION_TARGET_GMV_SWITCH_KEY, fmt.Sprintf("%d-%d", gmvStat.MainId, gmvStat.ChildId), req.Status)

	// 操作特殊用户的补偿状态
	err = service.HdPartitionGmvService.SwitchPartitionTestUserCompensateStatus(c, gmvStat.MainId, gmvStat.ChildId, req.Status)
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	response.ResponseSuccess(nil)
}

// @Summary 修改量化配比倍数
// @Description 修改量化配比倍数
// @Tags GMV 目标
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.HdTargetGmvQuantRatioUpdateReq true "查询参数"
// @Success 200 {object} any
// @Router /target_gmv/quant_ratio/update [post]
func UpdateTargetGmvQuantRatio(c *gin.Context) {
	req := form.HdTargetGmvQuantRatioUpdateReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	// 更新到 adm 操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err := models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      fmt.Sprintf("修改量化配比倍数为%.2f", req.QuantRatio),
		Scenes:      constant.OPERATE_PARTITION_TARGET_GMV,
		AssociateId: 0,
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}

	// 更新或创建全局量化配比，保留两位小数，main_id、child_id 固定为 0
	err = hd_task_models.HdPartitionGmvQuantRatioDal.Upsert(c, 0, 0, decimal.NewFromFloat(req.QuantRatio).Round(2).InexactFloat64())
	if err != nil {
		logrus.Error(err)
		response.ResponseFail(err.Error())
		return
	}

	response.ResponseSuccess(nil)
}
