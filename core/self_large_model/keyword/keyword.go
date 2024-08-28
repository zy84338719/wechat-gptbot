package keyword

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/sdk/conf"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/self_large_model"
)

var (
	open string
	gpt  *self_large_model.SelfGpt
)

func Init(ctx context.Context) error {
	// 初始化关键词
	gpt = self_large_model.NewSelfGpt(config.C.SelfGpt.BaseUrl, config.C.SelfGpt.KeyWordAuthorization)
	profile, err := gpt.Profile(context.Background())
	if err != nil {
		return err
	}
	open, err = conf.RedisClient.Get(ctx, "wechat_bot_keyword_token").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(open) != 0 {
		return nil
	}
	open, err = gpt.ChatOpen(ctx, profile.Id)
	if err != nil {
		return err
	}

	conf.RedisClient.Set(ctx, "wechat_bot_keyword_token", open, 30*24*time.Hour)
	return nil
}

func GetKeyword(ctx context.Context, content string) (string, error) {
	message, err := gpt.ChatMessage(ctx, open, self_large_model.ChatMessageRequest{
		Message: content,
		ReChat:  false,
		Stream:  false,
	})
	if err != nil {
		return "", err
	}
	return message.Content, nil
}
