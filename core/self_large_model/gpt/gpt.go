package gpt

import (
	"context"
	"wechat-gptbot/config"
	"wechat-gptbot/core/self_large_model"
)

var (
	profileId string
	gpt       *self_large_model.SelfGpt
)

func Init() error {
	// 初始化关键词
	gpt = self_large_model.NewSelfGpt(config.C.SelfGpt.BaseUrl, config.C.SelfGpt.GptAuthorization)
	profile, err := gpt.Profile(context.Background())
	if err != nil {
		return err
	}
	profileId = profile.Id
	return nil
}

func GPT(ctx context.Context, chatId string, content string) (string, string, error) {
	if chatId == "" {
		var err error
		chatId, err = gpt.ChatOpen(ctx, profileId)
		if err != nil {
			return "", "", err
		}
	}
	message, err := gpt.ChatMessage(ctx, chatId, self_large_model.ChatMessageRequest{
		Message: content,
		ReChat:  false,
		Stream:  false,
	})
	if err != nil {
		return "", "", err
	}
	return message.Content, message.ChatId, nil
}
