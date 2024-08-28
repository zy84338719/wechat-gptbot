package ability

import (
	"github.com/imroc/req/v3"
	"math/rand"
	"time"
	"wechat-gptbot/core/ability/proto"
)

type kfcApiResponse struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

type kfc struct {
}

func newKfc() *kfc {
	return &kfc{}
}

var kfcList = []string{}

func init() {
	kfcList = getCrazyKFCSentence()
}

func getCrazyKFCSentence() []string {
	var data []kfcApiResponse
	api := "https://raw.gitmirror.com/whitescent/KFC-Crazy-Thursday/main/kfc.json"
	if err := req.C().SetTimeout(time.Minute).Get(api).Do().Into(&data); err != nil {
		return nil
	}
	sentence := make([]string, 0)
	for i := range data {
		if len(data[i].Text) != 0 {
			sentence = append(sentence, data[i].Text)
		}
	}
	return sentence
}

func (a kfc) TextFunc(text string) string {
	if len(kfcList) == 0 {
		return "😭！今天肯德基初始化失败了，明天再来吧"
	}
	rand.Seed(time.Now().Unix())
	return kfcList[rand.Int()%len(kfcList)]
}

func (k kfc) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "疯狂星期四",
		Long:    "输入疯狂星期四即可获得当前的肯德基疯狂星期四的句子",
		Keyword: "疯狂星期四",
	}
}
