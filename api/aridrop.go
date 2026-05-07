package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
)

// TODO: 暂无使用
func AddUserAirdropTask(c *gin.Context) {
	req := form.AddUserAirdropTaskReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	id, ok := c.Params.Get("uid")
	if !ok {
		response.Responses(errno.Error, errno.MsgFlags[errno.InvalidParams], nil)
	}
	uid, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		response.Responses(errno.Error, errno.MsgFlags[errno.InvalidParams], nil)
		return
	}
	if req.DropTime < time.Now().Add(time.Minute*constant.TaskTime).Format(time.DateTime) {
		response.Responses(errno.TaskTimeFailed, errno.MsgFlags[errno.TaskTimeFailed], nil)
		return
	}
	if !until.RequestRateLimit(c, fmt.Sprintf("operation:add_%d_airdrop_task", uid), time.Second*3) {
		response.Responses(errno.RequestTooOfter, errno.MsgFlags[errno.RequestTooOfter], nil)
		return
	}
	operatorId := c.GetUint64("user_id")

	_, err = models.User{Ctx: c}.GetSysUserByUserId(uid)
	if err != nil {
		response.Responses(errno.Error, err.Error(), nil)
		return
	}
	dropCount := req.SendNum * req.Count
	_, err = models.AirdropTask{Ctx: c}.AddAirdropTask(models.AirdropTaskModel{
		OperatorId:    operatorId,
		OperatorName:  c.GetString("adm_user_name"),
		Title:         req.Title,
		BenefitType:   req.BenefitType,
		PropId:        req.PropId,
		PropName:      req.PropName,
		Source:        constant.ManualInput,
		UserList:      cast.ToString(uid),
		DropUserNum:   1,
		SendNum:       uint32(req.SendNum),
		DropCount:     dropCount,
		DropTime:      null.TimeFromPtr(until.StrTimeToTimePtr(req.DropTime)),
		AirdropType:   req.AirdropType,
		ProductSizeId: req.ProductSizeId,
	})
	if err != nil {
		response.Responses(errno.Error, err.Error(), nil)
		return
	}
	response.ResponseSuccess(nil)
}

func AirdropTreasure(c *gin.Context) {
	req := form.AirdropTreasureReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.Responses(errno.Error, err.Error(), nil)
		return
	}
	operatorId := c.GetUint64("user_id")
	admUser, err := models.User{Ctx: c}.FindAdmUserById(cast.ToUint64(operatorId))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	err = models.SendProp(c, models.SendPropReq{
		UserId:     uint64(user.UserId),
		PropId:     req.PropId,
		Number:     uint32(req.Number),
		Tm:         time.Now(),
		Source:     "DROP",
		PropSource: 3,
		Remark:     fmt.Sprintf("后台空投#操作员:%d, 空投用户:%d", operatorId, user.UserId),
	})
	if err != nil {
		response.Responses(errno.Error, err.Error(), nil)
		return
	}
	target, _ := models.NewAiMatchProductNftProp().One(c, req.PropId, "")
	httpReq.FeiShuRootBot("%s(%d) 空投 %s (mobile: %s) 如下道具: 《%s(%s)》%d 份", admUser.Name, admUser.UserId, user.RealName, user.Mobile, target.ProductTitle, target.SubTitle, req.Number)
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "空投道具",
		Scenes:      68,
		AssociateId: int64(req.PropId),
		RequestData: string(jsonStr),
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}
