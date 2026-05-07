package until

import (
	"hotbox-adm-backend/cli"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func GetAdmUserId(c *gin.Context) uint64 {
	token := c.GetHeader("token")
	val, err := cli.RedisCli.HGet(ctx, token, "user_id").Result()
	if err != nil {
		logrus.Panic("userid no.", err)
	}
	userId := cast.ToUint64(val)
	return userId
}
