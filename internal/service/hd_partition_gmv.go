package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/util"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var HdPartitionGmvService = &HdPartitionGmvSrv{}

type HdPartitionGmvSrv struct{}

// 用户选的小区，大区gmv要包含这个用户
// 用户选的大区，小区gmv不会包含这个用户

// SwitchPartitionTestUserCompensateStatus 分区切换补偿状态,并更新用户的分成比例
func (s *HdPartitionGmvSrv) SwitchPartitionTestUserCompensateStatus(ctx context.Context, mainId int64, childId int64, compensateStatus int) error {
	where := map[string][]any{
		"main_id":  {mainId},
		"child_id": {childId},
		"status":   {0}, // 非暂停状态
	}
	if compensateStatus == constant.COMPENSATE_STATUS_OPEN {
		where["compensate_status"] = []any{constant.COMPENSATE_STATUS_CLOSE} // 未开启分账状态
	} else {
		where["compensate_status"] = []any{constant.COMPENSATE_STATUS_OPEN} // 开启分账状态
	}
	// 配置了子分区
	if childId != 0 {
		where["child_id"] = []any{childId}
	}
	yopTestUsers, err := hd_task_models.HdYopTestUserDal.GetHdYopTestUsers(ctx, where, []string{}, 0)
	if err != nil {
		return err
	}

	eg, ctx2 := errgroup.WithContext(ctx)
	eg.SetLimit(50)
	for _, u := range yopTestUsers {
		testUser := u
		eg.Go(func() error {
			err = s.SwitchTestUserCompensateStatus(ctx2, testUser.Id, compensateStatus)
			if err != nil {
				return err
			}
			return s.RefreshUserRate(ctx2, testUser.Id)
		})
	}
	return eg.Wait()
}

// SwitchTestUserCompensateStatus 切换用户补偿状态
func (s *HdPartitionGmvSrv) SwitchTestUserCompensateStatus(ctx context.Context, testUserId int64, compensateStatus int) error {
	updateMap := make(map[string]any)
	whereMap := make(map[string]any)
	if compensateStatus == constant.COMPENSATE_STATUS_OPEN { // 补偿中
		whereMap = map[string]any{
			"id":                testUserId,
			"compensate_status": constant.COMPENSATE_STATUS_CLOSE,
		}
		updateMap = map[string]any{
			"pre_compensate_rate": gorm.Expr("rate"),               // 记录补偿前的分成比例
			"rate":                100,                             // 100进入公司账户
			"compensate_status":   constant.COMPENSATE_STATUS_OPEN, // 开启补偿
		}
	} else {
		whereMap = map[string]any{
			"id":                testUserId,
			"compensate_status": constant.COMPENSATE_STATUS_OPEN,
		}
		// 恢复正常分成比例
		updateMap = map[string]any{
			"rate":              gorm.Expr("pre_compensate_rate"),
			"compensate_status": constant.COMPENSATE_STATUS_CLOSE, // 关闭补偿
		}
	}
	err := hd_task_models.HdYopTestUserDal.UpdateByParams(ctx, whereMap, updateMap)
	if err != nil {
		return err
	}
	return nil
}

// CheckPartitionAndUpdateRate 添加/修改分区，检查并操作特殊用户的分成比例
func (s *HdPartitionGmvSrv) CheckPartitionAndUpdateRate(ctx context.Context, yopTestUserid int64) error {
	testUser, err := hd_task_models.HdYopTestUserDal.One(ctx, yopTestUserid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if testUser.Id == 0 {
		return nil
	}
	if testUser.MainId <= 0 { // 没有分区
		return nil
	}
	partitionIds := [][]any{}
	partitionIds = append(partitionIds, []any{testUser.MainId, testUser.ChildId})
	where := map[string]any{
		"status":   1,
		"date = ?": time.Now().Format(util.DefaultFormat),
		"main_id":  testUser.MainId,
	}

	childIds := []int64{testUser.ChildId}
	if testUser.ChildId != 0 {
		// 用户小分区可以对应大GMV分区
		childIds = append(childIds, 0)
	}
	mainIds := []int64{testUser.MainId}
	gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(ctx, where, mainIds, childIds)
	if err != nil {
		return err
	}

	// 该用户处于补偿中的分区, 开启补偿
	if len(gmvStats) > 0 {
		if testUser.CompensateStatus == constant.COMPENSATE_STATUS_CLOSE {
			err = s.SwitchTestUserCompensateStatus(ctx, yopTestUserid, constant.COMPENSATE_STATUS_OPEN)
		}
	} else {
		// 该用户处于非补偿中的分区, 关闭补偿
		if testUser.CompensateStatus == constant.COMPENSATE_STATUS_OPEN {
			err = s.SwitchTestUserCompensateStatus(ctx, yopTestUserid, constant.COMPENSATE_STATUS_CLOSE)
		}
	}
	return err
}

// 更新用户当前分成比例
func (s *HdPartitionGmvSrv) RefreshUserRate(ctx context.Context, yopTestUserid int64) (err error) {
	testUser, err := hd_task_models.HdYopTestUserDal.One(ctx, yopTestUserid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if testUser.Id == 0 {
		return errors.New("用户不存在")
	}

	if err = hd_task_models.HdYopTestUserRateRecordDal.Create(ctx, &hd_task_models.HdYopTestUserRateRecord{
		YopTestUserId: yopTestUserid,
		Rate:          testUser.Rate,
		MainId:        testUser.MainId,
		ChildId:       testUser.ChildId,
	}); err != nil {
		return err
	}

	hdYopTestUserCacheKey := fmt.Sprintf("cardmart:yop_divide_test_user:%d", testUser.UserId)
	cli.HotDogRedis.Set(ctx, hdYopTestUserCacheKey, testUser.Rate, 0)

	// 设置冻结比例缓存
	if testUser.FreezeRate > 0 {
		hdYopTestUserFreezeRateKey := fmt.Sprintf("matrix:yop_divide_test_user_freeze_rate:%d", testUser.UserId)
		cli.HotDogRedis.Set(ctx, hdYopTestUserFreezeRateKey, testUser.FreezeRate, 0)
	} else {
		// 如果冻结比例为0，删除缓存
		hdYopTestUserFreezeRateKey := fmt.Sprintf("matrix:yop_divide_test_user_freeze_rate:%d", testUser.UserId)
		cli.HotDogRedis.Del(ctx, hdYopTestUserFreezeRateKey)
	}
	return nil
}
