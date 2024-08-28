package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/sdk/conf"
	"time"
)

// RateLimiter 定义限流器结构体
type RateLimiter struct {
	redisClient *redis.Client
	keyPrefix   string
	windows     []time.Duration
	limits      []int
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(keyPrefix string, windows []time.Duration, limits []int) *RateLimiter {
	if len(windows) != len(limits) {
		panic("windows and limits must have the same length")
	}
	return &RateLimiter{
		redisClient: conf.RedisClient,
		keyPrefix:   keyPrefix,
		windows:     windows,
		limits:      limits,
	}
}

// Allow 判断是否允许请求
func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()
	pipeline := r.redisClient.Pipeline()

	for _, window := range r.windows {
		windowKey := fmt.Sprintf("%s:%s:%d", r.keyPrefix, key, window)
		// 获取窗口内的请求数
		pipeline.ZCard(ctx, windowKey)
	}

	// 执行管道命令
	results, err := pipeline.Exec(ctx)
	if err != nil {
		return false, err
	}

	// 检查所有窗口内的请求数是否超过限制
	for i := 0; i < len(r.windows); i++ {
		count := results[i].(*redis.IntCmd).Val()
		if count >= int64(r.limits[i]) {
			return false, nil
		}
	}

	for _, window := range r.windows {
		windowStart := now - int64(window.Seconds())
		windowKey := fmt.Sprintf("%s:%s:%d", r.keyPrefix, key, window)
		// 移除窗口外的请求
		pipeline.ZRemRangeByScore(ctx, windowKey, "0", fmt.Sprintf("%d", windowStart))
		// 添加当前请求
		pipeline.ZAdd(ctx, windowKey, redis.Z{Score: float64(now), Member: now})
		// 设置key的过期时间
		pipeline.Expire(ctx, windowKey, window)
	}

	// 执行管道命令
	results, err = pipeline.Exec(ctx)
	if err != nil {
		return true, err
	}
	return true, nil
}
