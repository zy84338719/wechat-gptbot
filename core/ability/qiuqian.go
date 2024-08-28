package ability

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
	"wechat-gptbot/utils"
)

type qiuqian struct {
}

func newQiuqian() *qiuqian {
	return &qiuqian{}
}

type GuaOpenRequest struct {
	GuaID    uint   `json:"gua_id"`
	Question string `json:"question"`
}

type GuaOpenResponse struct {
	ChatId string `json:"chat_id"`
}

type GuaDetailResponse struct {
	Gua Gua `json:"gua"`
}
type Gua struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	GuaType      int64  `json:"gua_type"`
	GuaTypeImage string `json:"gua_type_image"`
	GuaLevel     int64  `json:"gua_level"`
	GuaLevelStr  string `json:"gua_level_str"`
	GuaDesc      string `json:"gua_desc"`
}

func (qiuqian) TextFunc(text string) string {
	response := GuaOpenResponse{}
	i := utils.GetRandInt64(63) + 1
	err := req.C().DevMode().SetTimeout(3 * time.Minute).SetBaseURL(config.C.MoyuDomain).Post("/gua/talk/open").
		SetBody(GuaOpenRequest{
			GuaID:    uint(i),
			Question: text,
		}).Do().Into(&response)
	if err != nil {
		return "算卦服务异常"
	}
	response2 := GuaDetailResponse{}
	err = req.C().DevMode().SetTimeout(3 * time.Minute).SetBaseURL(config.C.MoyuDomain).Get(fmt.Sprintf("/gua/detail?gua_id=%d", i)).Do().Into(&response2)
	if err != nil {
		return "算卦服务异常"
	}
	resp := ShortResp{}
	_ = req.C().DevMode().SetTimeout(time.Minute).SetBaseURL(config.C.ShortUrlDomain).Post("/temp/shorturl/create").SetBody(
		map[string]interface{}{
			"expireInMinutes": 30,
			"url":             "https://games.murphyyi.com/suangua/chat/" + response.ChatId}).Do(context.Background()).Into(&resp)
	return "卦象：" + response2.Gua.Title + "\n卦等级：" + response2.Gua.GuaLevelStr + "\n卦内容: " + response2.Gua.GuaDesc + "\n" + "算卦详情：" + "https://d.murphyyi.com/t/" + resp.Data.Code
}

func (qiuqian) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "求签",
		Long:    "求签",
		Keyword: "求签",
	}
}
