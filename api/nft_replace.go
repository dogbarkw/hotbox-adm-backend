package api

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"
	"hotbox-adm-backend/util"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func GetDisplaceList(c *gin.Context) {
	req := form.GetDisplaceListReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	where := map[string]any{
		"is_delete": 0,
	}
	if req.Id > 0 {
		where["replace_id"] = req.Id
	}

	if req.IsSecondaryActive > 0 {
		where["is_secondary_active"] = 1
		if req.IsUnused > 0 {
			where["use_by_main_id"] = 0
		}
	}
	switch req.ActiveType {
	case 0:
		where["active_type"] = []uint32{0, 2, 4, 6}
	case 1:
		where["active_type"] = []uint32{1, 3, 5}
	}

	switch req.Status {
	case 1:
		where["on_sale_status"] = 1
	case 2:
		where["on_sale_status"] = 0
	}
	if req.Name != "" {
		where["replace_name LIKE ?"] = "%" + req.Name + "%"
	}
	offset := (req.PageNumber - 1) * req.PageSize
	list, total, err := models.NewAiMatchProductNftReplace().GetAiMatchProductNftReplaceList(c, where, []string{"replace_id DESC"}, &req.PageSize, &offset)
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	var result []dto.NftReplaceListData

	for _, item := range list {
		var status uint32
		if item.OnSaleStatus == 1 {
			status = 1 // 上架
		} else {
			status = 2 // 下架
		}
		var images []string
		if item.DetailPictures != "" {
			_ = json.Unmarshal([]byte(item.DetailPictures), &images)
		}
		temp := &dto.NftReplaceListData{
			Id:                 uint64(item.ReplaceId),
			Eid:                util.EncodeInt64(int64(item.ReplaceId)),
			Name:               item.ReplaceName,
			CoverImg:           item.HeadPicture,
			Images:             images,
			TotalNum:           uint32(item.TotalCount),
			DisplaceNum:        uint32(item.ReplaceCount),
			MinVersion:         uint32(item.MinVersion),
			Status:             status,
			ActiveType:         uint32(item.ActiveType),
			DrawEndTime:        item.DrawEndTime.Format(time.DateTime),
			PublishTime:        item.PublishTime.Format(time.DateTime),
			IsSecondaryActive:  uint32(item.IsSecondaryActive),
			UsedMainActive:     uint64(item.UseByMainId),
			ReserveNum:         item.ReserveNum,
			ReferenceActiveIds: []uint64{},
		}
		if len(item.ReferenceActiveId) > 0 {
			rids := make([]uint64, 0)
			for _, id := range strings.Split(item.ReferenceActiveId, ",") {
				_id, _ := strconv.Atoi(id)
				rids = append(rids, uint64(_id))
			}
			temp.ReferenceActiveIds = rids
		}
		result = append(result, *temp)
	}
	response.ResponseSuccessWithList(result, int(total))
}

func UpdateNftReplaceDisplaceCount(c *gin.Context) {
	req := form.UpdateNftReplaceDisplaceCountReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	replaceId, ok := c.Params.Get("replaceId")
	if !ok {
		response.Responses(errno.Error, errno.MsgFlags[errno.InvalidParams], nil)
	}
	payload := map[string]any{}
	where := map[string]any{}
	// 二者只能存在一个
	if (req.AlterRemainCount != nil) == (req.EditRemainCount != nil) {
		response.Responses(errno.Error, errno.GetMsg(errno.InvalidParams), nil)
		return
	}
	existRecord, err := models.NewAiMatchProductNftReplace().GetByReplaceId(c, cast.ToInt(replaceId))
	var whereTarget int
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	if req.EditRemainCount != nil {
		whereTarget = int(*req.EditRemainCount)
		where["total_count - reserve_num  >= ?"] = whereTarget
		payload["replace_count"] = gorm.Expr("total_count - reserve_num - ?", whereTarget)
	}
	if req.AlterRemainCount != nil {
		whereTarget = int(*req.AlterRemainCount)
		where["(total_count - reserve_num + ? - replace_count)  >=  0"] = whereTarget
		payload["replace_count"] = gorm.Expr("replace_count - ?", *req.AlterRemainCount)
	}
	rowAffected, err := models.NewAiMatchProductNftReplace().UpdateAiMatchProductNftReplaceByReplaceId(c, cast.ToInt(replaceId), where, payload)
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	if req.EditRemainCount != nil && rowAffected == 0 && existRecord.ReplaceCount != (existRecord.TotalCount-existRecord.ReserveNum-int64(whereTarget)) {
		response.ResponseFail("置换剩余数量不足")
		return
	}
	// 更新到adm操作日志
	operatorId := c.GetUint64("user_id")
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "修改置换剩余数量",
		Scenes:      1,
		AssociateId: cast.ToInt64(replaceId),
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	response.ResponseSuccess(nil)
}

func UpdateNftReplaceReserveNum(c *gin.Context) {
	req := form.UpdateNftReplaceReserveNumReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	replaceId, ok := c.Params.Get("replaceId")
	if !ok {
		response.Responses(errno.Error, errno.MsgFlags[errno.InvalidParams], nil)
	}
	payload := map[string]any{
		"reserve_num": req.ReserveNum,
	}
	where := map[string]any{
		"total_count - replace_count >= ?": req.ReserveNum,
	}
	existRecord, err := models.NewAiMatchProductNftReplace().GetByReplaceId(c, cast.ToInt(replaceId))
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	rowAffected, err := models.NewAiMatchProductNftReplace().UpdateAiMatchProductNftReplaceByReplaceId(c, cast.ToInt(replaceId), where, payload)
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	if rowAffected == 0 && existRecord.ReserveNum != int64(*req.ReserveNum) {
		response.ResponseFail("预留数量不足")
		return
	}
	// 更新到adm操作日志
	operatorId := c.GetUint64("user_id")
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "修改置换预留数量",
		Scenes:      1,
		AssociateId: cast.ToInt64(replaceId),
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	response.ResponseSuccess(nil)
}
