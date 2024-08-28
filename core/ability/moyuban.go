package ability

import (
	"github.com/imroc/req/v3"
	"wechat-gptbot/core/ability/proto"
)

type MoYu1Resp struct {
	Success bool   `json:"success"`
	Url     string `json:"url"`
}

type moyu struct {
}

func (moyu) TextFunc(g2 string) string {
	resp := req.C().Get("https://api.vvhan.com/api/moyu?type=json").Do()
	if resp.StatusCode != 200 {
		return ""
	}
	var respData MoYu1Resp
	if err := resp.Into(&respData); err != nil {
		return ""
	}
	if respData.Url == "" {
		return ""
	}
	return respData.Url
}

func (moyu) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "摸鱼",
		Long:    "随机获取一张摸鱼图片",
		Keyword: "摸鱼",
	}
}
