package ability

import (
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type gasoline struct {
}

func newGasoline() *gasoline {
	return &gasoline{}
}

func (gasoline) TextFunc(str string) string {
	var data gasolineResp
	err := req.C().SetBaseURL(config.C.CrawlerDomain).Post("/info/gasoline").Do().Into(&data)
	if err != nil {
		return "😭！查询遇到错误"
	}
	now := time.Now()
	s := now.Format("今天是01月02日")

	if data.Data.Time != nil {
		s += "\n调价消息" + data.Data.Time.Format("时间：01月02日晚上")
	}
	for _, d := range data.Data.Prices {
		if strings.Contains(str, d.Area) {
			s += "\n"
			s += fmt.Sprintf("92号: %.2f 元", float64(d.Price92)*0.01)
			s += "\n"
			s += fmt.Sprintf("95号: %.2f 元", float64(d.Price95)*0.01)
			s += "\n"
			s += fmt.Sprintf("98号: %.2f 元", float64(d.Price98)*0.01)
			s += "\n"
			s += fmt.Sprintf("0号: %.2f 元", float64(d.Price0)*0.01)
		}
	}
	return s
}

func (gasoline) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "油价 北京",
		Long:    "【省】\n输入即可获得当天油价信息",
		Keyword: "油价",
	}
}

type gasolineResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Msg    string     `json:"msg"`
		Time   *time.Time `json:"time"`
		Prices []struct {
			Area    string `json:"area"`
			Price92 int    `json:"price_92"`
			Price95 int    `json:"price_95"`
			Price98 int    `json:"price_98"`
			Price0  int    `json:"price_0"`
		} `json:"prices"`
	} `json:"data"`
}
