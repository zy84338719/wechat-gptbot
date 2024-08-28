package ability

import (
	"fmt"
	"github.com/imroc/req/v3"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type hotday struct {
}

func newHotday() *hotday {
	return &hotday{}
}

func (a hotday) TextFunc(class string) string {
	var m hotdayData
	err := req.C().SetBaseURL((config.C.CrawlerDomain)).Get("/getBy/" + class).Do().Into(&m)
	if err != nil {
		//global.GVA_LOG.Error(fmt.Sprintf("请求今日热榜遇到错误 %+v", err))
		return ""
	}
	mk := class

	for i, content := range m.Content {
		mk += fmt.Sprintf("%d. ", i+1) + content.Title + "\n"
		if i > 10 {
			break
		}
	}
	mk += "\n\n欢迎访问今日热榜\nV1.0 \nhttps://hotday.murphyyi.com/ \nV2.0\nhttps://hotnews.murphyyi.com/#/"

	return mk
}

func (a hotday) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "新闻 综合",
		Long:    "输入 任意字符 即可获得 ",
		Keyword: "新闻",
	}
}

type hotdayData struct {
	HotName     string              `json:"hot_name"`
	EnHotName   string              `json:"en_hot_name"`
	Content     []hotdayDataContent `json:"content"`
	CrawlerTime int64               `json:"crawler_time"`
}

type hotdayDataContent struct {
	Category  *string `json:"category,omitempty"`
	Href      string  `json:"href"`
	Title     string  `json:"title"`
	TopicFlag int     `json:"topic_flag,omitempty"`
	Num       int     `json:"num,omitempty"`
}
