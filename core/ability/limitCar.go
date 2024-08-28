package ability

import (
	"github.com/imroc/req/v3"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type limitCar struct {
}

func newLimitCar() *limitCar {
	return &limitCar{}
}

func (limitCar) TextFunc(string string) string {
	var data limitCarResp
	err := req.C().SetBaseURL(config.C.CrawlerDomain).Post("/info/limitCarWeek/" + string).Do().Into(&data)
	if err != nil {
		return "😭！没有查到限号信息"
	}
	var s = ""
	now := time.Now()
	for _, d := range data.Data {
		if d.TimeInfo.Day() == now.Day() {
			s += "\n今天是 " + d.LimitedTime + "\n【" + d.LimitedWeek + "】 \n限行:" + d.LimitedNumber
		}
	}
	return s
}

func (limitCar) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "限号 北京",
		Long:    "【省】\n输入即可获得当天限号信息 目前仅支持北京",
		Keyword: "限号",
	}
}

type limitCarResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		TimeInfo      time.Time `json:"timeInfo"`
		LimitedTime   string    `json:"limitedTime"`
		LimitedWeek   string    `json:"limitedWeek"`
		LimitedNumber string    `json:"limitedNumber"`
	} `json:"data"`
}
