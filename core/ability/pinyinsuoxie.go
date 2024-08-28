package ability

import (
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"strings"
	"wechat-gptbot/core/ability/proto"
)

type suoxie struct {
}

func newSuoxie() *suoxie {
	return &suoxie{}
}

func (suoxie) transPinYinSuoXie(text string) (string, error) {
	api := "https://lab.magiconch.com/api/nbnhhsh/guess"
	resp := req.C().Post(api).SetFormData(map[string]string{"text": text}).Do()
	var ret []string
	gjson.Get(resp.String(), "0.trans").ForEach(func(key, val gjson.Result) bool {
		ret = append(ret, val.String())
		return true
	})
	return strings.Join(ret, "；"), nil
}

func (suoxie) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "缩写 wtf",
		Long:    "输入缩写即可获得对应的缩写",
		Keyword: "缩写",
	}
}

func (a suoxie) TextFunc(keyword string) string {
	xie, err := a.transPinYinSuoXie(keyword)
	if err != nil {
		return "啊，我在缩写字典中没有找到对应的数据"
	}
	if len(xie) == 0 {
		return "很抱歉我没有找到热门的网络缩写 😭"
	}
	return xie
}
