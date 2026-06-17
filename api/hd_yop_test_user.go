package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/util"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	"hotbox-adm-backend/cli"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @BasePath /aiera/v2/operation

// @Summary 添加
// @Description 添加
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserAddReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/create [post]
func YopTestUserAdd(c *gin.Context) {
	req := form.YopTestUserAddReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if user.UserId == 0 {
		response.ResponseFail("无此用户，请检查")
		return
	}

	testUser := &hd_task_models.HdYopTestUser{
		UserId:     user.UserId,
		RealName:   user.RealName,
		Mobile:     req.Mobile,
		UserType:   req.UserType,
		Rate:       req.Rate,
		FreezeRate: req.FreezeRate,
		Remark:     req.Remark,
	}
	row, err := hd_task_models.HdYopTestUserDal.FirstOrCreate(c, testUser)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	if row == 0 {
		response.ResponseFail("该用户已添加")
		return
	}

	if err = hd_task_models.HdYopTestUserRateRecordDal.Create(c, &hd_task_models.HdYopTestUserRateRecord{
		YopTestUserId: testUser.Id,
		Rate:          testUser.Rate,
	}); err != nil {
		response.ResponseFail(err.Error())
		return
	}

	hdYopTestUserCacheKey := fmt.Sprintf("hotbox:yop_divide_test_user:%d", user.UserId)
	cli.HotDogRedis.Set(c, hdYopTestUserCacheKey, req.Rate, 0)

	// 设置冻结比例缓存
	if req.FreezeRate > 0 {
		hdYopTestUserFreezeRateKey := fmt.Sprintf("matrix:yop_divide_test_user_freeze_rate:%d", user.UserId)
		cli.HotDogRedis.Set(c, hdYopTestUserFreezeRateKey, req.FreezeRate, 0)
	}

	// 更新到adm操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "添加特殊账号",
		Scenes:      constant.OPERATE_YOP_TEST_USER,
		AssociateId: cast.ToInt64(testUser.Id),
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}

	response.ResponseSuccess(nil)
}

// @Summary 编辑
// @Description 编辑
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserUpdateReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/update [post]
func YopTestUserUpdate(c *gin.Context) {
	req := form.YopTestUserUpdateReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	// 更新到adm操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err := models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "修改特殊账号信息",
		Scenes:      constant.OPERATE_YOP_TEST_USER,
		AssociateId: cast.ToInt64(req.Id),
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}

	err = hd_task_models.HdYopTestUserDal.UpdateByParams(c, map[string]any{"id": req.Id}, map[string]any{
		"rate":        req.Rate,
		"freeze_rate": req.FreezeRate,
		"remark":      req.Remark,
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	if err = hd_task_models.HdYopTestUserRateRecordDal.Create(c, &hd_task_models.HdYopTestUserRateRecord{
		YopTestUserId: req.Id,
		Rate:          req.Rate,
	}); err != nil {
		response.ResponseFail(err.Error())
		return
	}

	testUser, _ := hd_task_models.HdYopTestUserDal.One(c, req.Id)
	hdYopTestUserCacheKey := fmt.Sprintf("hotbox:yop_divide_test_user:%d", testUser.UserId)
	cli.HotDogRedis.Set(c, hdYopTestUserCacheKey, req.Rate, 0)

	// 设置冻结比例缓存
	if req.FreezeRate > 0 {
		hdYopTestUserFreezeRateKey := fmt.Sprintf("matrix:yop_divide_test_user_freeze_rate:%d", testUser.UserId)
		cli.HotDogRedis.Set(c, hdYopTestUserFreezeRateKey, req.FreezeRate, 0)
	} else {
		hdYopTestUserFreezeRateKey := fmt.Sprintf("matrix:yop_divide_test_user_freeze_rate:%d", testUser.UserId)
		cli.HotDogRedis.Del(c, hdYopTestUserFreezeRateKey)
	}

	response.ResponseSuccess(nil)
}

// @Summary 删除
// @Description 删除
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserDelReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/del [post]
func YopTestUserDel(c *gin.Context) {
	req := form.YopTestUserDelReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	list, err := hd_task_models.HdYopTestUserDal.GetHdYopTestUsers(c, map[string][]any{
		"id in ?": {req.Ids},
	}, []string{}, len(req.Ids))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	if err = hd_task_models.HdYopTestUserDal.Delete(c, req.Ids); err != nil {
		response.ResponseFail(err.Error())
		return
	}

	delUserId := make([]int64, 0, len(list))
	keys := make([]string, 0, len(list))
	for _, user := range list {
		keys = append(keys, fmt.Sprintf("hotbox:yop_divide_test_user:%d", user.UserId))
		// 只有测试用户才标记删除
		if user.UserType == 2 {
			delUserId = append(delUserId, user.UserId)
		}
	}
	cli.HotDogRedis.Del(c, keys...)

	// 禁用账号
	if len(delUserId) > 0 {
		err = models.User{Ctx: c}.DeleteTestUserByIds(delUserId)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}

	// 删除用户的分成比例记录
	hd_task_models.HdYopTestUserRateRecordDal.Delete(c, req.Ids)

	// 更新到adm操作日志
	operatorId, _ := c.Get("user_id")
	jsonStr, _ := json.Marshal(req)
	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      cast.ToInt64(operatorId),
		Username:    c.GetString("adm_user_name"),
		Remark:      "批量删除特殊账号",
		Scenes:      constant.OPERATE_YOP_TEST_USER,
		AssociateId: 0,
		RequestData: string(jsonStr),
	})
	if err != nil {
		klog.Error(err)
		response.ResponseFail(errno.MsgFlags[errno.Error])
		return
	}
	response.ResponseSuccess(nil)
}

// @Summary 列表
// @Description 列表
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserListReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/list [post]
func YopTestUserList(c *gin.Context) {
	req := form.YopTestUserListReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
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

	where := map[string][]any{}
	if req.Mobile != "" {
		where["mobile"] = []any{req.Mobile}
	}

	var (
		data                  []hd_task_models.HdYopTestUser
		count                 int64
		total                 float64
		todayIncome           float64
		selectTotalDateIncome float64
		err                   error
	)
	start, end := util.GetStartAndEndOfDay(time.Now().Local())

	selectStartTime := time.Now()
	if len(req.StartTime) > 0 {
		selectStartTime, err = time.ParseInLocation(util.StandardFormat, req.StartTime, time.Local)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}

	selectEndTime := time.Now()
	if len(req.EndTime) > 0 {
		selectEndTime, err = time.ParseInLocation(util.StandardFormat, req.EndTime, time.Local)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
	}

	eg, _ := errgroup.WithContext(c)
	eg.Go(func() error {
		data, count, err = hd_task_models.HdYopTestUserDal.GetList(c, where, []string{"id desc"}, req.PageNumber, req.PageSize)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		total, err = hd_task_models.HdYopTestUserDal.SumTotalIncome(c)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		todayIncome, err = hd_task_models.HdYopTestUserIncomeRecordDal.SumTodayIncome(c, start, end)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
			selectTotalDateIncome, err = hd_task_models.HdYopTestUserIncomeRecordDal.SumTodayIncome(c, selectStartTime, selectEndTime)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		return nil
	})

	if err = eg.Wait(); err != nil {
		response.ResponseFail(err.Error())
		return
	}

	testUserIds := make([]int64, 0, len(data))
	for _, v := range data {
		testUserIds = append(testUserIds, v.Id)
	}
	userIncomeMap, err := hd_task_models.HdYopTestUserIncomeRecordDal.SumIncomeByTimeRange(c, testUserIds, start, end)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseFail(err.Error())
		return
	}

	var selectDateUserIncomeMap map[int64]float64
	if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
		selectDateUserIncomeMap, err = hd_task_models.HdYopTestUserIncomeRecordDal.SumIncomeByTimeRange(c, testUserIds, selectStartTime, selectEndTime)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.ResponseFail(err.Error())
			return
		}

	}

	type YopTestUserData struct {
		hd_task_models.HdYopTestUser
		DailyIncome      float64 `json:"daily_income"`
		SelectDateIncome float64 `json:"select_date_income"`
		ForceFreeze                int64   `json:"force_freeze"`
		UserSelectDateForceFreeze  int64   `json:"user_select_date_force_freeze"`
	}
	list := make([]YopTestUserData, 0, len(data))
	err = copier.Copy(&list, &data)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	// 查询所有用户的冻结金额
	userIds := make([]int64, 0, len(data))
	for _, v := range data {
		userIds = append(userIds, v.UserId)
	}
	forceFreezeMap := make(map[int64]int64)
	if len(userIds) > 0 {
		wallets, err := models.SysUserWalletDal.GetUserWallets(c, userIds)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Errorf("YopTestUserList GetUserWallets err: %v", err)
		}
		for _, w := range wallets {
			forceFreezeMap[w.UserId] = w.ForceFreeze
		}
	}

	var selectDateForceFreezeSum int64
	var selectDateForceFreezeUserMap map[int64]int64
	if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
		selectDateForceFreezeSum, err = models.SysUserWalletForceFreezeRecordDal.SumFreezeByUserIdsAndTimeRange(c, userIds, selectStartTime, selectEndTime)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Errorf("YopTestUserList SumFreezeByUserIdsAndTimeRange err: %v", err)
		}
		selectDateForceFreezeUserMap, err = models.SysUserWalletForceFreezeRecordDal.SumFreezeGroupByUserAndTimeRange(c, userIds, selectStartTime, selectEndTime)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Errorf("YopTestUserList SumFreezeGroupByUserAndTimeRange err: %v", err)
		}
	}

	var totalForceFreeze int64
	for i := range list {
		income, ok := userIncomeMap[list[i].Id]
		if ok {
			list[i].DailyIncome = income
		}

		if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
			selectDateIncome, ok := selectDateUserIncomeMap[list[i].Id]
			if ok {
				list[i].SelectDateIncome = selectDateIncome
			}
		}

		if ff, ok := forceFreezeMap[list[i].UserId]; ok {
			if ff < 0 {
				ff = -ff
			}
			list[i].ForceFreeze = ff
			totalForceFreeze += ff
		}

		if len(req.StartTime) > 0 && len(req.EndTime) > 0 {
			if userFreeze, ok := selectDateForceFreezeUserMap[list[i].UserId]; ok {
				list[i].UserSelectDateForceFreeze = userFreeze
			}
		}
	}

	code := errno.Success
	msg := errno.GetMsg(errno.Success)
	c.JSON(200, gin.H{
		"code":                     code,
		"msg":                      msg,
		"data":                     list,
		"res_count":                count,
		"all":                      total,
		"daily_all":                todayIncome,
		"select_total_date_income": selectTotalDateIncome,
		"total_force_freeze":       totalForceFreeze,
		"select_date_force_freeze": selectDateForceFreezeSum,
	})
}

// @Summary 检查用户是否注册
// @Description 检查用户是否注册
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserCheckReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/check [post]
func YopTestUserCheck(c *gin.Context) {
	req := form.YopTestUserCheckReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(user)
}

// @Summary 获取平台用户零钱余额
// @Description 获取平台用户零钱余额
// @Tags 特殊账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query form.YopTestUserIdReq true "查询参数"
// @Success 200 {object} any
// @Router /yop_test_user/balance [post]
func YopTestUserBalance(c *gin.Context) {
	req := form.YopTestUserIdReq{}
	response := until.NewResponse(c)
	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	yopTestUser, err := hd_task_models.HdYopTestUserDal.One(c, req.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseFail(err.Error())
		return
	}
	if yopTestUser.UserId == 0 {
		response.ResponseFail("无此用户")
		return
	}
	user, err := models.SysUserWalletDal.GetUserWallet(c, yopTestUser.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseSuccess(models.SysUserWallet{TotalBalance: 0, FreezeBalance: 0, ForceFreeze: 0})
		return
	}
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(user)
}
