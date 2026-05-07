package api

import (
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @BasePath /hotbox/operation

// @Summary 列表
// @Description 列表
// @Tags Gpt活动策划会话
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetEventPlanningConversationListReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/conversation/list [post]
func GetEventPlanningConversationList(c *gin.Context) {
	req := form.GetEventPlanningConversationListReq{}
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
	userId, _ := c.Get("user_id")
	res, total, err := models.EventPlanningConversation{Ctx: c}.GetEventPlanningConversationList(req, cast.ToInt(userId), []string{"id DESC"})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccessWithList(res, int(total))
}

// @Summary 创建
// @Description 创建
// @Tags Gpt活动策划会话
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.CreateEventPlanningConversationListReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/conversation/create [post]
func CreateEventPlanningConversation(c *gin.Context) {
	req := form.CreateEventPlanningConversationListReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	_, err := models.EventPlanningConversation{Ctx: c}.CreateEventPlanningConversation(cast.ToInt(userId), req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}
