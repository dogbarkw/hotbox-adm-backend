package route

import (
	"os"
	"time"

	"hotbox-adm-backend/api"
	"hotbox-adm-backend/middleware"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func RegisterPartitionGmvRouters(r *gin.Engine) {
	r.Use(until.CORS())
	cardmartV2 := r.Group("/hotbox/v2/operation", until.LoginChk())
	if os.Getenv("ENV") == "dev" {
		// 本地环境不验token
		cardmartV2 = r.Group("/hotbox/v2/operation")
	}

	rivalTradeTask := cardmartV2.Group("/target_gmv", middleware.RateLimitMiddleware(time.Duration(cast.ToInt64(os.Getenv("RATE_LIMIT_SECOND")))*time.Second, cast.ToInt(os.Getenv("RATE_LIMIT_MAX_COUNT"))))
	{
		rivalTradeTask.POST("/list", api.TargetGmvList)
		rivalTradeTask.POST("/update", api.UpdateTargetGmv)
		rivalTradeTask.POST("/switch", api.TargetGmvSwitch)
		rivalTradeTask.POST("/quant_ratio/info", api.GetTargetGmvQuantRatio)
		rivalTradeTask.POST("/quant_ratio/update", api.UpdateTargetGmvQuantRatio)
	}
}
