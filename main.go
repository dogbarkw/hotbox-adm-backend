package main

import (
	"os"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/corn/target_gmv"
	cornJob "hotbox-adm-backend/corn"
	"hotbox-adm-backend/route"

	"github.com/gin-gonic/gin"
)

var port = ":8888"

func init() {
	cli.InitEnv()
	// cli.InitGormDB()
	cli.InitHDGormDB()
	cli.InitHDTaskDB()
	cli.InitHDRedis()
	cli.InitHDADBGormDB()
	cli.InitSpecialUserIds()

	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}
}

func main() {
	logrus.Infof("ENV: %s \n", os.Getenv("ENV"))
	if os.Getenv("ENV") != "dev" && os.Getenv("START_CRON") != "false" {
		go func() {
			c := cron.New(cron.WithSeconds())
			// 按数量回收
			c.AddJob("*/1 * * * * *", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewRecycleRecordCornJob()))
			// 按数量空投
			c.AddJob("@every 1s", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewAirdropRecordCornJob()))
			c.AddJob("*/1 * * * * *", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewNftReservedCornJob()))
			// c.AddJob("0 */1 * * * *", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewArtistNftActivityCountCornJob())) // 艺术家藏品数，活动数统计
			// c.AddJob("0 */1 * * * *", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewActivityScoreInitCornJob()))      // 活动结束数据快照
			c.AddJob("@every 5m", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewDailyGmvCornJob()))      // 统计GMV
			c.AddJob("0 0 2 * * *", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewPreDailyGmvCornJob())) // 统计前日GMV
			c.AddJob("@every 5m", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cornJob.NewYopTestUserJob()))       // 特殊账号累计进账
			c.AddJob("@every 5m", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(target_gmv.NewDgTargetGmvCornJob())) // 目标gmv统计

			c.Start()
			select {}
		}()
	}
	router := gin.Default()
	route.RegisterRouters(router)
	route.RegisterDocRouters(router)
	// route.RegisterGptRouters(router)
	// route.RegisterActivityScoreRouters(router)
	route.RegisterYopTestUserRouters(router)
	route.RegisterPartitionGmvRouters(router)
	router.Run(port)
}
