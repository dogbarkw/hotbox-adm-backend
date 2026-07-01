package target_gmv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/internal/service"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/util"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func NewDgTargetGmvCornJob() *DgTargetGmvCornJob {
	return &DgTargetGmvCornJob{}
}

type DgTargetGmvCornJob struct{}

func (p *DgTargetGmvCornJob) Run() {
	ctx := context.Background()
	lockKey := "hotbox:DgTargetGmvCornJob:lock"
	lock := cli.HotDogRedis.SetNX(ctx, lockKey, "lock", 3*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, lockKey)

	currentTime := time.Now()
	// 统计分区进账
	err := p.StartStatPartitionIncome(ctx, currentTime)
	if err != nil {
		logrus.Error(err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, "统计分区进账失败", err.Error())
	}
}

func (p *DgTargetGmvCornJob) StartStatPartitionIncome(ctx context.Context, countTime time.Time) error {
	start, end := util.GetStartAndEndOfDay(countTime)
	gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(ctx, map[string]any{
		"date": countTime.Format(util.DefaultFormat),
	}, nil, nil)
	if err != nil {
		return err
	}
	gmvStatsMap := make(map[string]*hd_task_models.HdPartitionDailyGmvStat)
	for _, v := range gmvStats {
		gmvStatsMap[fmt.Sprintf("%d-%d", v.MainId, v.ChildId)] = &v
		v.ShareIncome = 0
	}
	partitionIncomeList, err := hd_task_models.HdPartitionIncomeRecordDal.SumPartitionIncomeByTimeRange(ctx, [][]any{}, start, end)
	if err != nil {
		return err
	}
	for _, incomeRecord := range partitionIncomeList {
		shareIncome := decimal.NewFromFloat(incomeRecord.Income).Div(decimal.NewFromFloat(0.07)).Round(2).InexactFloat64() // 进账/0.07
		if incomeRecord.ChildId != 0 {
			// 用户配置子区间，GMV 配置一级区间
			if v, ok := gmvStatsMap[fmt.Sprintf("%d-%d", incomeRecord.MainId, 0)]; ok {
				v.ShareIncome += shareIncome
			}
		}
		if v, ok := gmvStatsMap[fmt.Sprintf("%d-%d", incomeRecord.MainId, incomeRecord.ChildId)]; ok {
			v.ShareIncome += shareIncome
		}
	}
	for _, v := range gmvStatsMap {
		where := map[string]any{
			"date":     v.Date,
			"main_id":  v.MainId,
			"child_id": v.ChildId,
		}
		status := v.Status
		if v.ShareIncome+v.CurrentGmv+v.PreGmv > v.TargetGmv && v.Status != 2 {
			status = 2
		}
		err = hd_task_models.HdPartitionDailyGmvStatDal.UpdateByParams(ctx, where, map[string]any{
			"share_income": v.ShareIncome,
			"status":       status,
		})
		if err != nil {
			return err
		}
		// 关闭补偿
		if status == 2 && v.Status != 2 {
			if err = service.HdPartitionGmvService.SwitchPartitionTestUserCompensateStatus(ctx, v.MainId, v.ChildId, constant.COMPENSATE_STATUS_CLOSE); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *DgTargetGmvCornJob) StartUpdateDailyGmv(ctx context.Context, gmvTime time.Time) error {
	gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(ctx, map[string]any{
		"date": gmvTime.Format(util.DefaultFormat),
	}, nil, nil)
	if err != nil {
		return err
	}

	eg, ctx2 := errgroup.WithContext(ctx)
	for _, v := range gmvStats {
		stat := v
		eg.Go(func() error {
			categoryPath := cast.ToString(stat.MainId)
			where := map[string]any{
				"ymd": gmvTime.Format(time.DateOnly),
			}
			if stat.ChildId > 0 {
				categoryPath = cast.ToString(stat.MainId) + "," + cast.ToString(stat.ChildId)
				where["category_path"] = categoryPath
			} else {
				// 大分区，统计大分区GMV
				where["category_path  LIKE ?"] = cast.ToString(stat.MainId) + "%"
			}
			gmv, err := hd_task_models.HdDailyNftCategoryGmvDal.GetGMVByParams(ctx2, where)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Infof("该分区找不到对应的gmv, ymd:%s, category_path:%s", gmvTime.Format(time.DateOnly), categoryPath)
				return nil
			}
			if err != nil {
				return err
			}

			status := v.Status
			if gmv+stat.ShareIncome+stat.PreGmv >= stat.TargetGmv && v.Status != 2 {
				status = 2
			}

			err = hd_task_models.HdPartitionDailyGmvStatDal.UpdateByParams(ctx, map[string]any{"id": v.Id}, map[string]any{
				"current_gmv": gmv,
				"status":      status,
			})
			if err != nil {
				return err
			}

			// 关闭补偿
			if status == 2 && v.Status != 2 {
				if err = service.HdPartitionGmvService.SwitchPartitionTestUserCompensateStatus(ctx, stat.MainId, stat.ChildId, constant.COMPENSATE_STATUS_CLOSE); err != nil {
					return err
				}
			}
			return nil
		})

	}

	return eg.Wait()
}

// 生成新的一天的目标gmv
func (p *DgTargetGmvCornJob) GenNewDateTargetGmv(ctx context.Context) error {
	currentTime := time.Now()
	lockKey := fmt.Sprintf("hotbox:GenNewDateTargetGmv_%s", currentTime.Format(time.DateOnly))
	lock := cli.HotDogRedis.SetNX(ctx, lockKey, "lock", time.Hour*24).Val()
	if !lock {
		return nil
	}

	// 查询昨天的数据
	gmvStats, err := hd_task_models.HdPartitionDailyGmvStatDal.GetByParams(ctx, map[string]any{
		"date": currentTime.AddDate(0, 0, -1).Format(util.DefaultFormat),
	}, nil, nil)
	if err != nil {
		cli.HotDogRedis.Del(ctx, lockKey)
		return err
	}

	for _, lastRecord := range gmvStats {
		preGmv := lastRecord.PreGmv + lastRecord.CurrentGmv + lastRecord.ShareIncome - lastRecord.TargetGmv
		if preGmv <= 0 || currentTime.Weekday().String() == "Monday" {
			// 周一不累计
			preGmv = 0
		}

		compensateStatus := 1
		cacheStatus := cli.HotDogRedis.HGet(ctx, constant.REDIS_HD_PARTITION_TARGET_GMV_SWITCH_KEY, fmt.Sprintf("%d-%d", lastRecord.MainId, lastRecord.ChildId)).Val()
		if cast.ToInt(cacheStatus) == 2 {
			compensateStatus = 2
		}

		err = hd_task_models.HdPartitionDailyGmvStatDal.UpSertNextGmv(ctx, &hd_task_models.HdPartitionDailyGmvStat{
			Date:      cast.ToInt(currentTime.Format(util.DefaultFormat)),
			MainId:    lastRecord.MainId,
			ChildId:   lastRecord.ChildId,
			TargetGmv: lastRecord.TargetGmv,
			PreGmv:    decimal.NewFromFloat(preGmv).Round(2).InexactFloat64(),
			Status:    compensateStatus,
		})
		if err != nil {
			cli.HotDogRedis.Del(ctx, lockKey)
			return err
		}
	}

	return nil
}
