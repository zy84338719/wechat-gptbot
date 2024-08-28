package ability

import (
	"github.com/imroc/req/v3"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type tiangouResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type tiangou struct {
}

func newTiangou() *tiangou {
	return &tiangou{}
}

func (tiangou) TextFunc(text string) string {
	var data tiangouResp
	err := req.C().SetBaseURL(config.C.MoyuDomain).Get("/data/moyu/tiangou/randomOne").Do().Into(&data)
	if err != nil {
		return "😭！我的爱人不开心了"
	}
	return data.Data
}

func (tiangou) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "舔狗日记",
		Long:    "输入即可获得舔狗日记 一篇 可以访问 games.murphyyi.com/tiangou 查看更多",
		Keyword: "舔狗",
	}
}

func (t tiangou) buildChatCache() *proto.ChatCacheInfo {
	return &proto.ChatCacheInfo{
		AbilityName:     t.help().Keyword,
		LastTime:        time.Now().Unix(),
		CacheExpiredSec: 60,
	}
}
