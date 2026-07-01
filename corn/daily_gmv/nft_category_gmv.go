package daily_gmv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func StartStatNftCategoryGmv(ctx context.Context, startTime string, endTime string) {
	msgType := "统计分区GMV"
	dailyGMVList, err := models.AiMatchProductOrder{Ctx: &gin.Context{}}.GetNftCategoryDailyGMV(startTime, endTime, cli.SpecialUserIds)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("获取统计分区gmv记录失败，err=%s", err.Error())
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
		dailyGmvModel := hd_task_models.HdDailyNftCategoryGmv{
			Ymd:          ymd.Format(time.DateOnly),
			UserCnt:      dailyGMV.UserCnt,
			Gmv:          dailyGMV.Gmv,
			CategoryPath: dailyGMV.CategoryPath,
			Category:     dailyGMV.Category,
			Rk:           dailyGMV.Rk,
		}
		rowsAffected, err := hd_task_models.HdDailyNftCategoryGmvDal.FirstOrCreate(ctx, dailyGmvModel)
		if err != nil {
			httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, fmt.Sprintf("创建统计分区gmv记录失败，err=%s", err.Error()))
			return
		}
		if rowsAffected == 0 { // update
			err = hd_task_models.HdDailyNftCategoryGmvDal.UpdateByParams(ctx, map[string]any{
				"ymd":           ymd.Format(time.DateOnly),
				"category_path": dailyGMV.CategoryPath,
			}, map[string]any{
				"category": dailyGMV.Category,
				"gmv":      dailyGMV.Gmv,
				"user_cnt": dailyGMV.UserCnt,
				"rk":       dailyGMV.Rk,
			})
			if err != nil {
				httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_GMV_URL, msgType, fmt.Sprintf("更新统计分区gmv记录失败，err=%s", err.Error()))
				return
			}
		}
	}
}
