package ability

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"regexp"
	"strings"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type weather struct {
}

func newWeather() *weather {
	return &weather{}
}
func (w *weather) TextFunc(text string) string {
	getWeather := w.getWeather(context.Background(), text)
	//<text>天气：多云</text><br><text>气温：31℃</text><br><text>体感温度：32℃</text><br><text>风向：西南风</text><br><text>风力：2级 风速：6km/h</text><br><text>湿度：52% 大气压强：992hPa</text><br><text>本小时降水量：0.0mm</text><br><text>能见度：20km</text><br>
	// 将上述网页标签剔除，转化为普通文本

	return stripHtmlTagsAndFormat(text + "\n" + getWeather)
}

func stripHtmlTagsAndFormat(input string) string {
	// 定义正则表达式来匹配 HTML 标签
	re := regexp.MustCompile(`<[^>]*>`)
	// 替换 <br> 标签为换行符
	input = strings.ReplaceAll(input, "<br>", "\n")
	// 用空字符串替换其他 HTML 标签
	output := re.ReplaceAllString(input, "")
	return output
}

func (w weather) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "天气",
		Long:    "输入城市名即可获得对应城市的天气",
		Keyword: "天气",
	}
}

func (weather) getWeather(ctx context.Context, city string) string {
	w := weatherReq{}
	w.Uri = fmt.Sprintf("/qweather/now/%s", city)
	wResp := weatherResp{}
	err := req.C().DevMode().SetBaseURL(config.C.CrawlerDomain).Post("/info/common").SetBody(w).Do().Into(&wResp)
	if err != nil {
		return "晴，但是我不知道是哪里的天气"
	}
	if wResp.Code != 0 {
		return "晴，但是我不知道是哪里的天气"
	}
	return wResp.Data.Items[0].Description
}

type weatherReq struct {
	Uri     string `json:"uri"`
	Foreign bool   `json:"foreign"`
}

type weatherResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		Link          string    `json:"link"`
		FeedLink      string    `json:"feedLink"`
		Links         []string  `json:"links"`
		Updated       string    `json:"updated"`
		UpdatedParsed time.Time `json:"updatedParsed"`
		Author        struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Authors []struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"authors"`
		Language   string `json:"language"`
		Generator  string `json:"generator"`
		Extensions struct {
			Atom struct {
				Link []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
					Attrs struct {
						Href string `json:"href"`
						Rel  string `json:"rel"`
						Type string `json:"type"`
					} `json:"attrs"`
					Children struct {
					} `json:"children"`
				} `json:"link"`
			} `json:"atom"`
		} `json:"extensions"`
		Items []struct {
			Title           string    `json:"title"`
			Description     string    `json:"description"`
			Link            string    `json:"link"`
			Links           []string  `json:"links"`
			Published       string    `json:"published"`
			PublishedParsed time.Time `json:"publishedParsed"`
			Guid            string    `json:"guid"`
		} `json:"items"`
		FeedType    string `json:"feedType"`
		FeedVersion string `json:"feedVersion"`
	} `json:"data"`
}
