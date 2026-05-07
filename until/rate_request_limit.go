package until

import (
	"context"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

var keyPrefixCache = "new_back_end_"

func RequestRateLimit(ctx context.Context, key string, duration time.Duration) bool {
	key = fmt.Sprintf("%s:%s", keyPrefixCache, key)
	defer func() {
		cli.HotDogRedis.Expire(ctx, key, duration)
	}()
	rate, err := cli.HotDogRedis.Incr(ctx, key).Result()
	if err != nil && err.Error() != redis.Nil.Error() {
		klog.Errorf("RequestRateLimit key:%s, error:%v", key, err)
		return false
	}
	if rate > 1 {
		klog.Errorf("RequestRateLimit too often, key:%s, val:%d", key, rate)
		return false
	}
	return true
}
