package route

import (
	"os"

	_ "hotbox-adm-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterDocRouters(r *gin.Engine) {
	if os.Getenv("ENV") != "production" {
		r.GET("/hotbox/operation/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
