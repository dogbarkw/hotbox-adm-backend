package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/pkg/constant"
	pkg "hotbox-adm-backend/pkg/util"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

var keyPrefixCache = "new_backend"

func RateLimitMiddleware(d time.Duration, target int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(os.Getenv("ENV")) != "production" {
			ctx.Next()
			return
		}
		ip := ctx.GetHeader("real_ip_header")
		if ip == "" {
			ip = pkg.ClientIP(ctx)
		}
		if ip == "" {
			ctx.Next()
			return
		}
		ok, _ := cli.HotDogRedis.SIsMember(ctx, constant.REDIS_IP_WHITELIST_KEY, ip).Result()
		if ok {
			ctx.Next()
			return
		}
		response := until.NewResponse(ctx)
		if target == 0 {
			response.Responses(1001, "操作太频繁，请稍后重试~", nil)
			ctx.Abort()
			return
		}

		key := fmt.Sprintf("%s:%s:%s", keyPrefixCache, ctx.Request.URL.Path, ip)
		defer func() {
			cli.HotDogRedis.Expire(ctx, key, d)
		}()
		rate, err := cli.HotDogRedis.Incr(ctx, key).Result()
		if err != nil && err.Error() != redis.Nil.Error() {
			response.Responses(1001, "操作太频繁，请稍后重试~", nil)
			ctx.Abort()
			return
		}
		if rate > int64(target) {
			//if ip != "" {
			//	d := models.BlockIpModel{Ip: ip}
			//	models.BlockIp{Ctx: ctx}.FirstOrCreate(&d)
			//}
			response.Responses(1001, "操作太频繁，请稍后重试~", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
