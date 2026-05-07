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

func RegisterGptRouters(r *gin.Engine) {
	r.Use(until.CORS())
	poker := r.Group("/hotbox/operation", until.LoginChk())
	if os.Getenv("ENV") == "dev" {
		// 本地环境不验token
		poker = r.Group("/hotbox/operation")
	}

	gpt := poker.Group("/gpt", middleware.RateLimitMiddleware(time.Duration(cast.ToInt64(os.Getenv("RATE_LIMIT_SECOND")))*time.Second, cast.ToInt(os.Getenv("RATE_LIMIT_MAX_COUNT"))))
	{
		gpt.POST("/event_planning/schema", api.GetGptEventPlanningSchema)
		gpt.POST("/event_planning/collection_info", api.GetGptEventPlanningCollectionInfo)
		gpt.POST("/event_planning/schema_msg", api.GetGptEventPlanningSchemaMsg)
		gpt.POST("/event_planning/schema_stream", api.GetGptEventPlanningSchemaStream)
		gpt.POST("/send_gpt_msg", api.SendGptMsg)
		gpt.POST("/send_gpt_msg_stream", api.SendGptMsgStream)
		gpt.POST("/gen_ai_article", api.GenAiArticle)
	}
	eventPlanningConversation := poker.Group("/gpt/event_planning/conversation")
	{
		eventPlanningConversation.POST("/list", api.GetEventPlanningConversationList)
		eventPlanningConversation.POST("/create", api.CreateEventPlanningConversation)
	}
}
