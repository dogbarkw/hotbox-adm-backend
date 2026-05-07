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

func RegisterActivityScoreRouters(r *gin.Engine) {
	r.Use(until.CORS())
	poker := r.Group("/hotbox/operation", until.LoginChk())
	if os.Getenv("ENV") == "dev" {
		// 本地环境不验token
		poker = r.Group("/hotbox/operation")
	}
	activityScore := poker.Group("/activity_score", middleware.RateLimitMiddleware(time.Duration(cast.ToInt64(os.Getenv("RATE_LIMIT_SECOND")))*time.Second, cast.ToInt(os.Getenv("RATE_LIMIT_MAX_COUNT"))))
	{
		activityScore.POST("/pending/list", api.GetPendingActivityScoreList)
		activityScore.POST("/pending/update", api.UpdatePendingActivityScore)
		activityScore.POST("/artist_recommend_score/list", api.GetArtistRecommendScoreList)
		activityScore.POST("/artist_recommend_score/update", api.UpdateArtistRecommendScore)
		activityScore.POST("/ended/list", api.GetEndedActivityScoreList)
	}
}
