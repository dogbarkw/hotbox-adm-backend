package dg_yop_test_user

import (
	"context"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"github.com/shopspring/decimal"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type DgYopTestUserIncomeCornJob struct{}

func (p *DgYopTestUserIncomeCornJob) Run() {
	ctx := context.Background()
	lockKey := "hd:DgYopTestUserIncomeCornJob"
	lock := cli.HotDogRedis.SetNX(ctx, lockKey, "lock", 10*time.Minute).Val()
	if !lock {
		return
	}
	defer cli.HotDogRedis.Del(ctx, lockKey)

	p.StartSumIncome(ctx)
}

func (p *DgYopTestUserIncomeCornJob) StartSumIncome(ctx context.Context) {
	var (
		msgType   = "统计特殊用户累计进账"
		lastIndex int64
		limit     = 1000
		flag      bool
	)

	for !flag {
		list, err := hd_task_models.HdYopTestUserDal.GetHdYopTestUsers(ctx, map[string][]any{"id > ?": {lastIndex}}, []string{"id asc"}, limit)
		if err != nil {
			httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("获取特殊用户列表失败:%s", err))
			logrus.Error(err)
			return
		}

		// 统计用户的累计进账
		eg, ctx2 := errgroup.WithContext(ctx)
		for _, yopTestUserItem := range list {
			eg.Go(func() error {
				return p.SumUserTotalIncome(ctx2, yopTestUserItem)
			})
		}
		err = eg.Wait()
		if err != nil {
			httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_ROBOT_URL, msgType, fmt.Sprintf("获取统计特殊用户累计进账失败:%s", err))
			logrus.Error(err)
			return
		}

		if len(list) < limit {
			flag = true
			return
		}
		lastIndex = list[len(list)-1].Id
	}
}

type CountTimeRange struct {
	StartTime time.Time
	EndTime   time.Time
	Rate      int
}

func (p *DgYopTestUserIncomeCornJob) SumUserTotalIncome(ctx context.Context, testUser hd_task_models.HdYopTestUser) error {
	timeNow := time.Now().Local()

	lastCountTime := testUser.CountTime.Time

	if testUser.CountTime.Time.UnixMilli() < testUser.CreatedAt.UnixMilli() {
		lastCountTime = testUser.CreatedAt
	}

	// 获取期间所有的进账记录
	orderLogs, err := models.AiMatchProductOrder{Ctx: &gin.Context{}}.GetTestUserOrderLogs(testUser.UserId, lastCountTime, timeNow, []string{"pay_amount", "payment_time"})
	if err != nil {
		return errors.Wrapf(err, "统计特殊用户累计进账错误，user_id=%d", testUser.UserId)
	}

	if len(orderLogs) == 0 {
		return nil
	}

	// 获取期间特殊用户的比例记录
	rateRecords, err := hd_task_models.HdYopTestUserRateRecordDal.GetByParams(ctx, map[string][]any{
		"yop_test_user_id": {testUser.Id},
	}, []string{"id asc"}, 0)
	if err != nil {
		return errors.Wrapf(err, "获取用户期间比例记录错误，user_id=%d, created_at=%v", testUser.UserId, testUser.CreatedAt)
	}

	//特殊处理账号
	if testUser.Mobile == "16600006666" {
		testUser.Rate = 100
		for i := range rateRecords {
			rateRecords[i].Rate = 100
		}
	}

	var totalIncome float64
	incomeRecords := make([]*hd_task_models.HdYopTestUserIncomeRecord, 0, len(orderLogs))
	if len(rateRecords) == 0 {
		// 期间没有新增比例记录，取用户最新比例
		for _, orderLog := range orderLogs {
			income := p.CalculateTotalIncome(testUser.Rate, orderLog)
			totalIncome += income
			incomeRecords = append(incomeRecords, &hd_task_models.HdYopTestUserIncomeRecord{
				YopTestUserId: testUser.Id,
				Fee:           orderLog.PayAmount,
				Rate:          testUser.Rate,
				Income:        income,
				IncomeTime:    orderLog.PaymentTime,
			})
		}
	} else {
		for _, orderLog := range orderLogs {
			// 以对应区间的比例计算
			for i := 0; i < len(rateRecords); i += 1 {
				if i < len(rateRecords)-1 {
					if orderLog.PaymentTime.UnixNano() >= rateRecords[i].CreatedAt.UnixNano() && orderLog.PaymentTime.UnixNano() < rateRecords[i+1].CreatedAt.UnixNano() {
						income := p.CalculateTotalIncome(rateRecords[i].Rate, orderLog)
						totalIncome += income
						incomeRecords = append(incomeRecords, &hd_task_models.HdYopTestUserIncomeRecord{
							YopTestUserId: testUser.Id,
							Fee:           orderLog.PayAmount,
							Rate:          rateRecords[i].Rate,
							Income:        income,
							IncomeTime:    orderLog.PaymentTime,
						})
						break
					}
				} else {
					if orderLog.PaymentTime.UnixNano() >= rateRecords[i].CreatedAt.UnixNano() {
						income := p.CalculateTotalIncome(rateRecords[i].Rate, orderLog)
						totalIncome += income
						incomeRecords = append(incomeRecords, &hd_task_models.HdYopTestUserIncomeRecord{
							YopTestUserId: testUser.Id,
							Fee:           orderLog.PayAmount,
							Rate:          rateRecords[i].Rate,
							Income:        income,
							IncomeTime:    orderLog.PaymentTime,
						})
						break
					}
				}
			}
		}
	}

	if len(incomeRecords) > 0 {
		if err = hd_task_models.HdYopTestUserIncomeRecordDal.BatchCreate(ctx, incomeRecords); err != nil {
			return errors.Wrapf(err, "批量生成进账记录错误，user_id=%d", testUser.UserId)
		}
	}

	updateParams := map[string]any{
		"count_time": timeNow,
	}

	if testUser.CountTime.Time.UnixMilli() > 0 {
		updateParams["total_income"] = gorm.Expr("total_income + ?", totalIncome)
	} else {
		// 初次统计
		updateParams["total_income"] = totalIncome
	}

	err = hd_task_models.HdYopTestUserDal.UpdateByParams(ctx, map[string]any{"user_id": testUser.UserId}, updateParams)
	if err != nil {
		return errors.Wrapf(err, "更新特殊用户累计进账错误，user_id=%d", testUser.UserId)
	}
	return nil
}

// 计算进账
func (p *DgYopTestUserIncomeCornJob) CalculateTotalIncome(userRate int, userWalletLog models.AiMatchProductOrderModel) float64 {
	// 特殊账号分成比例
	rate := decimal.NewFromInt(int64(userRate)).Mul(decimal.NewFromFloat(0.01)).InexactFloat64()
	return decimal.NewFromFloat(userWalletLog.PayAmount).Mul(decimal.NewFromFloat(0.93 * rate)).Round(2).InexactFloat64()
}
