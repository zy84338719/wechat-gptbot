package self_large_model

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"time"
)

const (
	chatMessageUrl = "/application/chat_message/%s" // chat_id
	profileUrl     = "/application/profile"
	chatOpenUrl    = "/application/%s/chat/open"
)

type SelfGpt struct {
	client *req.Client
}

func NewSelfGpt(baseUrl string, authorization string) *SelfGpt {
	return &SelfGpt{
		client: req.C().SetBaseURL(baseUrl).SetTimeout(time.Second*time.Duration(180)).
			SetCommonHeader("AUTHORIZATION", authorization),
	}

}

func (s *SelfGpt) ChatMessage(ctx context.Context, chatId string, chatMsg ChatMessageRequest) (*ChatMessageData, error) {
	var cmp ChatMessageResp
	err := s.client.Post(fmt.Sprintf(chatMessageUrl, chatId)).SetBody(chatMsg).Do(ctx).Into(&cmp)
	if err != nil {
		return nil, err
	}
	if cmp.Code != http.StatusOK {
		return nil, fmt.Errorf("code:%d, message:%s", cmp.Code, cmp.Message)
	}
	if cmp.ChatMessageData == nil {
		return nil, fmt.Errorf("chatMessageData is nil")
	}
	return cmp.ChatMessageData, nil
}

func (s *SelfGpt) Profile(ctx context.Context) (*ProfileData, error) {
	var cmp ProfileResp
	err := s.client.Get(profileUrl).Do(ctx).Into(&cmp)
	if err != nil {
		return nil, err
	}
	if cmp.Code != http.StatusOK {
		return nil, fmt.Errorf("code:%d, message:%s", cmp.Code, cmp.Message)
	}
	if cmp.ProfileData == nil {
		return nil, fmt.Errorf("profileData is nil")
	}
	return cmp.ProfileData, nil
}

func (s *SelfGpt) ChatOpen(ctx context.Context, applicationId string) (string, error) {
	var cmp ChatOpen
	err := s.client.Get(fmt.Sprintf(chatOpenUrl, applicationId)).Do(ctx).Into(&cmp)
	if err != nil {
		return "", err
	}
	if cmp.Code != http.StatusOK {
		return "", fmt.Errorf("code:%d, message:%s", cmp.Code, cmp.Message)
	}
	if cmp.Data == "" {
		return "", fmt.Errorf("chatOpen data is nil")
	}
	return cmp.Data, nil
}
