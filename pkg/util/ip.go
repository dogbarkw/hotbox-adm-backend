package pkg

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func ClientIP(c *gin.Context) string {
	// 第三方抗D
	httpXForwardedFor := c.Request.Header.Get("HTTP_X_FORWARDED_FOR")
	ip := strings.TrimSpace(strings.Split(httpXForwardedFor, ",")[0])
	if ip != "" {
		log.Println("第三方抗D获取的客户端ip: ", ip)
		return ip
	}

	// X-Forwarded-For
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	// X-Real-Ip
	ip = strings.TrimSpace(c.Request.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	// X-Connecting-Ip
	ip = strings.TrimSpace(c.Request.Header.Get("X-Connecting-Ip"))
	if ip != "" {
		return ip
	}

	// client_ip
	return c.ClientIP()
}
