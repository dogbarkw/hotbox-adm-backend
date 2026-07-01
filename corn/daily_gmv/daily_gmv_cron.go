package daily_gmv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/corn/target_gmv"

	"hotbox-adm-backend/pkg/constant"

	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/hd_task_models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 统计当日数据
type DailyGmvCornJob struct{}

func (p *DailyGmvCornJob) Run() {
	startTime := time.Now().Format("2006-01-02 00:00:00")
	endTime := time.Now().Format(time.DateTime)
	ctx := context.Background()
	StartStatGmv(ctx, startTime, endTime)
	StartStatNftCategoryGmv(ctx, startTime, endTime)
	if err := target_gmv.NewDgTargetGmvCornJob().StartUpdateDailyGmv(ctx, time.Now()); err != nil {
		logrus.Error(err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, "更新分区gmv失败", err.Error())
	}
}

// 统计前一天的数据
type DailyBeforeGmvCornJob struct{}

func (p *DailyBeforeGmvCornJob) Run() {
	dayBefore := time.Now().AddDate(0, 0, -1)
	startTime := dayBefore.Format("2006-01-02 00:00:00")
	endTime := fmt.Sprintf("%s 23:59:59.999", dayBefore.Format(time.DateOnly))
	ctx := context.Background()
	StartStatGmv(ctx, startTime, endTime)
	StartStatNftCategoryGmv(ctx, startTime, endTime)

	c := target_gmv.NewDgTargetGmvCornJob()
	if err := c.StartStatPartitionIncome(ctx, dayBefore); err != nil {
		logrus.Error(err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, "统计分区进账失败", err.Error())
	}
	if err := c.StartUpdateDailyGmv(ctx, dayBefore); err != nil {
		logrus.Error(err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, "更新分区gmv失败", err.Error())
	}
	if err := c.GenNewDateTargetGmv(ctx); err != nil {
		logrus.Error(err.Error())
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, "生成新的一天的目标gmv失败", err.Error())
	}
}

func StartStatGmv(ctx context.Context, startTime string, endTime string) {
	msgType := "统计GMV"
	dailyGMVList, err := models.AiMatchProductOrderDal.GetDailyGMV(ctx, startTime, endTime)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("获取统计记录失败，err=%s", err.Error())
		logrus.Error(msg)
		httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, msg)
		return
	}
	for _, dailyGMV := range dailyGMVList {
		ymd := dailyGMV.Dt
		if dailyGMV.Dt.Unix() <= 0 {
			ymd, err = time.ParseInLocation(time.DateTime, startTime, time.Local)
			if err != nil {
				httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, fmt.Sprintf("解析时间失败，startTime=%s", startTime))
				return
			}
		}
		dailyGmvModel := hd_task_models.CmDailyGmv{
			Ymd:      ymd.Format(time.DateOnly),
			Gmv:      dailyGMV.Gmv,
			UserCnt:  dailyGMV.UserCnt,
			RGmv:     dailyGMV.RGmv,
			RUserCnt: dailyGMV.RUserCnt,
		}
		rowsAffected, err := hd_task_models.CmDailyGmvDal.FirstOrCreate(ctx, dailyGmvModel)
		if err != nil {
			httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, fmt.Sprintf("创建统计记录失败，err=%s", err.Error()))
			return
		}
		if rowsAffected == 0 { // update
			err = hd_task_models.CmDailyGmvDal.UpdateByParams(ctx, map[string]any{
				"ymd": ymd.Format(time.DateOnly),
			}, map[string]any{
				"gmv":        dailyGMV.Gmv,
				"user_cnt":   dailyGMV.UserCnt,
				"r_gmv":      dailyGMV.RGmv,
				"r_user_cnt": dailyGMV.RUserCnt,
			})
			if err != nil {
				httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, fmt.Sprintf("更新统计记录失败，err=%s", err.Error()))
				return
			}
		}
	}
}
