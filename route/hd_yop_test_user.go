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

func RegisterYopTestUserRouters(r *gin.Engine) {
	r.Use(until.CORS())
	aieraV2 := r.Group("/hotbox/v2/operation", until.LoginChk())
	if os.Getenv("ENV") == "dev" {
		// 本地环境不验token
		aieraV2 = r.Group("/hotbox/v2/operation")
	}

	yopTestUser := aieraV2.Group("/yop_test_user", middleware.RateLimitMiddleware(time.Duration(cast.ToInt64(os.Getenv("RATE_LIMIT_SECOND")))*time.Second, cast.ToInt(os.Getenv("RATE_LIMIT_MAX_COUNT"))))
	{
		yopTestUser.POST("/list", api.YopTestUserList)
		yopTestUser.POST("/create", api.YopTestUserAdd)
		yopTestUser.POST("/update", api.YopTestUserUpdate)
		yopTestUser.POST("/del", api.YopTestUserDel)
		yopTestUser.POST("/check", api.YopTestUserCheck)
		yopTestUser.POST("/balance", api.YopTestUserBalance)

		yopTestUser.POST("/stat/list", api.YopTestUserStatList)
		yopTestUser.POST("/stat/detail", api.YopTestUserStatDetail)
	}
}
