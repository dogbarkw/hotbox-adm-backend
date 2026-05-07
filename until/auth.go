package until

import (
	"context"
	"fmt"
	"os"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"

	"hotbox-adm-backend/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var ctx = context.Background()

type LoginStatusCode struct {
	Code string `json:"code"`
}

func unAuthorizationRes(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": errno.ErrorToken,
		"msg":  msg,
	})
}

func LoginChk() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("token")
		if tokenStr == "" {
			unAuthorizationRes(c, errno.GetMsg(errno.ErrorToken))
			c.Abort()
			return
		}
		isExist, _ := cli.HotDogRedis.Exists(ctx, tokenStr).Result()
		if isExist != 1 {
			logrus.Errorf("token->%s", tokenStr)
			unAuthorizationRes(c, errno.GetMsg(errno.ErrorToken))
			c.Abort()
			return
		}
		info := cli.HotDogRedis.HGetAll(ctx, tokenStr)
		userId, ok := info.Val()["user_id"]
		if !ok {
			unAuthorizationRes(c, errno.GetMsg(errno.ErrorToken))
			c.Abort()
			return
		}
		c.Set("user_id", userId)
		cacheKey := fmt.Sprintf("admin:token:%s", tokenStr)
		if cli.HotDogRedis.Exists(ctx, cacheKey).Val() == 1 {
			c.Next()
			return
		}
		u, err := models.User{Ctx: c}.FindAdmUserById(cast.ToUint64(userId))
		c.Set("adm_user_name", u.Name)
		if err != nil {
			unAuthorizationRes(c, err.Error())
			c.Abort()
			return
		}
		if u.OrgId == constant.IN_VALID_ORG_ID {
			unAuthorizationRes(c, err.Error())
			c.Abort()
			return
		}
		cli.HotDogRedis.Set(c, cacheKey, tokenStr, 3*time.Hour)
		c.Next()
	}
}

func GTokenChk() gin.HandlerFunc {
	return func(c *gin.Context) {
		gToken := c.Request.Header.Get("gtoken")
		if gToken != os.Getenv("G_TOKEN") || gToken == "" {
			unAuthorizationRes(c, errno.GetMsg(errno.ErrorToken))
			c.Abort()
			return
		}
		c.Next()
	}
}
