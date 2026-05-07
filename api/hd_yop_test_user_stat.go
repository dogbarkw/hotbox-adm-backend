package api

import (
	"time"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
)

// @Summary 特殊账号累记收入列表
// @Description 特殊账号累记收入列表
// @Tags 特殊账号收入统计
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} any
// @Router /yop_test_user/stat/list [post]
func YopTestUserStatList(c *gin.Context) {
	response := until.NewResponse(c)

	result, err := hd_task_models.HdYopTestUserIncomeRecordDal.StatByMouth(c)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	response.ResponseSuccess(result)
}

// @Summary 特殊账号收入详情
// @Description 特殊账号收入详情
// @Tags 特殊账号收入统计
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.YopTestUserStatDetailReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/stat/detail [post]
func YopTestUserStatDetail(c *gin.Context) {
	response := until.NewResponse(c)
	req := form.YopTestUserStatDetailReq{}
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	startTime, err := time.ParseInLocation("2006-01", req.Mouth, time.Local)
	if err != nil {
		response.ResponseFail("时间格式错误")
		return
	}

	timeEnd := startTime.AddDate(0, 1, 0)
	result, err := hd_task_models.HdYopTestUserIncomeRecordDal.GetIncomeRecord(c, startTime, timeEnd)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(result)
}
