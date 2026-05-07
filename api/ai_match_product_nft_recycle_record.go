package api

import (
	"errors"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func GetRecycleRecordList(c *gin.Context) {
	ntNftRecycleModel := models.AiMatchProductNftRecycleRecord{}
	req := form.RecycleRecordListReq{}
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

	res, total, err := ntNftRecycleModel.GetRecycleRecordList(c, req, nil)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccessWithList(res, int(total))
}

func GetRecycleRecordById(c *gin.Context) {
	response := until.NewResponse(c)
	id, ok := c.Params.Get("id")
	if !ok {
		response.Responses(errno.Error, errno.MsgFlags[errno.InvalidParams], nil)
	}
	data, err := models.AiMatchProductNftRecycleRecord{}.One(c, cast.ToInt(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ResponseFails("暂无此数据", 404)
			return
		}
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(data)
}
