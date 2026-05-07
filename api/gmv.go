package api

import (
	"strings"

	"github.com/shopspring/decimal"
	"hotbox-adm-backend/dto"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
)

// @Summary GMV统计列表
// @Description GMV统计列表
// @Tags GMV统计
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GmvListReq true "查询参数"
// @Success 200 {object} any
// @Router /hotbox/v2/operation/gmv/list [post]
func GmvList(c *gin.Context) {
	req := form.GmvListReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	queryMap := make(map[string]any)
	if req.YMD != "" {
		queryMap["ymd"] = req.YMD
	}
	if req.DateRange != "" {
		split := strings.Split(req.DateRange, ",")
		if len(split) > 1 {
			queryMap["ymd >= ?"] = split[0]
			queryMap["ymd <= ?"] = split[1]
		}
	}
	res, total, err := hd_task_models.CmDailyGmvDal.GetCmDailyGmvList(c, queryMap, []string{"ymd desc"}, req.PageNumber, req.PageSize)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	var (
		totalGmv  float64 = 0
		totalRGmv float64 = 0
	)

	for _, item := range res {
		totalGmv += item.Gmv
		totalRGmv += item.RGmv
	}
	data := dto.GmvResponse{
		List:     res,
		Total:    total,
		TotalGmv: totalGmv,
		Profit:   decimal.NewFromFloat(totalRGmv).Mul(decimal.NewFromFloat(0.05)).InexactFloat64(),
	}
	// response.ResponseSuccessWithList(res, int(total))
	response.ResponseSuccess(data)
}
