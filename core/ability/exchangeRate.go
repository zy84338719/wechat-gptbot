package ability

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/yi-nology/common/utils/xstrings"
	"strings"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type exchangeRate struct {
}

var searchList = []string{
	"英镑", "港币", "美元", "瑞士法郎", "德国马克", "法国法郎", "新加坡元", "瑞典克朗", "丹麦克朗", "挪威克朗",
	"日元", "加拿大元", "澳大利亚元", "欧元", "澳门元", "菲律宾比索", "泰国铢", "新西兰元", "韩国元", "卢布", "林吉特",
	"新台币", "西班牙比塞塔", "意大利里拉", "荷兰盾", "比利时法郎", "芬兰马克", "印度卢比", "印尼卢比", "巴西里亚尔",
	"阿联酋迪拉姆", "印度卢比", "南非兰特", "沙特里亚尔", "土耳其里拉",
}

//
// 请为上面数据添加双引号
//

type exchangeRateResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		SpotBuyingRate  float64 `json:"SpotBuyingRate"`
		CashBuyingRate  float64 `json:"CashBuyingRate"`
		SpotSellingRate float64 `json:"SpotSellingRate"`
		CashSellingRate float64 `json:"CashSellingRate"`
		ConversionRate  float64 `json:"ConversionRate"`
		ReleaseTime     string  `json:"ReleaseTime"`
		ReleaseTimeUnix int     `json:"ReleaseTimeUnix"`
	} `json:"data"`
}

func (e exchangeRate) TextFunc(text string) string {
	search := ""
	for _, s := range searchList {
		if strings.Contains(text, s) {
			search = s
			break
		}
	}
	resp := exchangeRateResp{}
	err := req.C().SetBaseURL(config.C.CrawlerDomain).Post("/info/exchangeRate").SetBody(map[string]interface{}{
		"currency": search,
		"limit":    1,
	}).Do().Into(&resp)
	if err != nil {
		return "😭！查询遇到错误"
	}
	if resp.Code != 0 {
		return "😭！查询遇到错误"
	}
	if len(resp.Data) == 0 {
		return "😭！查询遇到错误"
	}
	data := resp.Data[0]
	return "1" + search + "可以购买 人民币" +
		"数据时间：" + data.ReleaseTime + "\n" +
		"现汇买入价：" + fmt.Sprintf("%.5f", data.SpotBuyingRate/100.0) + "\n" +
		"现钞买入价：" + fmt.Sprintf("%.5f", data.CashBuyingRate/100.0) + "\n" +
		"现汇卖出价：" + fmt.Sprintf("%.5f", data.SpotSellingRate/100.0) + "\n" +
		"现钞卖出价：" + fmt.Sprintf("%.5f", data.CashSellingRate/100.0) + "\n" +
		"中行折算价：" + fmt.Sprintf("%.5f", data.ConversionRate/100.0)
}

func (e exchangeRate) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "汇率 美元",
		Long:    "美元\n输入即可获得汇率信息 \n\n" + xstrings.Join('\n', searchList...),
		Keyword: "汇率",
	}
}

func (e exchangeRate) buildChatCache() *proto.ChatCacheInfo {
	return &proto.ChatCacheInfo{
		AbilityName:     e.help().Keyword,
		LastTime:        time.Now().Unix(),
		CacheExpiredSec: 60,
	}
}

func newExchangeRate() *exchangeRate {
	return &exchangeRate{}
}
