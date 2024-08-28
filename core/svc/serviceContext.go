package svc

import (
	"wechat-gptbot/core/gpt"
)

type ServiceContext struct {
	Session gpt.Session
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		Session: gpt.NewSession(),
	}
}
