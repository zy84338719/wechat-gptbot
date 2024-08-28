package ability

import (
	"context"
	"github.com/imroc/req/v3"
	"net/url"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type film struct {
}

func newFilm() *film {
	return &film{}
}

type filmResppnse struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Id       int    `json:"id"`
			Cid      int    `json:"cid"`
			Pid      int    `json:"pid"`
			Name     string `json:"name"`
			SubTitle string `json:"subTitle"`
			CName    string `json:"cName"`
			State    string `json:"state"`
			Picture  string `json:"picture"`
			Actor    string `json:"actor"`
			Director string `json:"director"`
			Blurb    string `json:"blurb"`
			Remarks  string `json:"remarks"`
			Area     string `json:"area"`
			Year     string `json:"year"`
		} `json:"list"`
		Page struct {
			PageSize  int `json:"pageSize"`
			Current   int `json:"current"`
			PageCount int `json:"pageCount"`
			Total     int `json:"total"`
		} `json:"page"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func (a film) TextFunc(text string) string {
	resppnse := filmResppnse{}
	err := req.C().DevMode().SetTimeout(time.Minute).SetBaseURL(config.C.FilmDomain).Get("searchFilm?keyword=", url.QueryEscape(text)).Do().Into(&resppnse)
	if err != nil {
		return "电影查询失败"
	}
	if resppnse.Code != 0 {
		return resppnse.Msg
	}
	if len(resppnse.Data.List) == 0 {
		return "没有找到啊，🥹 \n可以试试剧棒v0.2  \nhttps://drama.murphyyi.com/"
	}
	resp := ShortResp{}
	err = req.C().SetTimeout(time.Minute).SetBaseURL(config.C.ShortUrlDomain).Post("/temp/shorturl/create").SetBody(map[string]interface{}{
		"expireInMinutes": 120,
		"url":             "https://projector.murphyyi.online/#/search?search=" + url.QueryEscape(text),
	}).Do(context.Background()).Into(&resp)
	if err != nil {
		return "没有找到啊，🥹 \n可以试试剧棒v0.2  \nhttps://drama.murphyyi.com/"
	}

	return "我找到的电影：\nhttps://d.murphyyi.com/t/" + resp.Data.Code + " \n如果你不知道想看什么，可以试试 剧棒🎉v0.2 https://drama.murphyyi.com/"
}

func (k film) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "影片",
		Long:    "输入影片名称",
		Keyword: "影片",
	}
}

type ShortResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Code string `json:"code"`
	} `json:"data"`
}
