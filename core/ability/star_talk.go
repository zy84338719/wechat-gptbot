package ability

import (
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"wechat-gptbot/config"
	"wechat-gptbot/core/ability/proto"
)

type starTalkResp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data XingZuoYunShi `json:"data"`
}

type XingZuoYunShi struct {
	TodayYunShi    YunShi `json:"today_yun_shi"`
	TomorrowYunShi YunShi `json:"tomorrow_yun_shi"`
	WeeklyYunShi   YunShi `json:"weekly_yun_shi"`
	MonthlyYunShi  YunShi `json:"monthly_yun_shi"`
	YearlyYunShi   YunShi `json:"yearly_yun_shi"`

	Phase   YsInfo `json:"phase"`
	Love    YsInfo `json:"love"`
	Cause   YsInfo `json:"cause"`
	Weather YsInfo `json:"weather"`
}

type YsInfo struct {
	Title            string             `json:"title"`
	Score            int                `json:"score"`
	PointList        []int              `json:"point_list"`
	YsContent        string             `json:"ys_content"`
	HouseContentList []HouseContentList `json:"house_content_list"`
}

type HouseContentList struct {
	PlanetIcon   []interface{} `json:"planet_icon"`
	Planets      string        `json:"planets"`
	YsTitle      string        `json:"ys_title"`
	YsTag        string        `json:"ys_tag"`
	YsContent    string        `json:"ys_content"`
	YsTotalDay   int           `json:"ys_total_day"`
	YsSurplusDay int           `json:"ys_surplus_day"`
	StartTime    string        `json:"start_time"`
	EndTime      string        `json:"end_time"`
	Power        float64       `json:"power"`
}

type starTalk struct {
}

func newStarTalk() *starTalk {
	return &starTalk{}
}

func (starTalk) TextFunc(text string) string {
	text = strings.Replace(text, "座", "", -1)
	var data starTalkResp
	err := req.C().DevMode().SetBaseURL(config.C.CrawlerDomain).Get(fmt.Sprintf("/info/star/talk?star=%s", text)).Do().Into(&data)
	if err != nil {
		return "结果解析异常"
	}

	if data.Code != 0 {
		return "QaQ 你的提问我不是很能理解，请换个方式再试试吧！"
	}

	shi := data.Data
	todayYunShi := shi.TodayYunShi
	yiList := strings.Join(todayYunShi.YiList, ", ")
	jiList := strings.Join(todayYunShi.JiList, ", ")

	return fmt.Sprintf("✨ %s运势（%s） ✨\n\n"+
		"📊 综合运势：%d 分\n%s"+
		"\n\n⚠️ 重点关注：%d 分\n%s"+
		"\n\n❤️ 爱情运势：%d 分\n%s"+
		"\n\n💼 事业财运：%d 分\n%s"+
		"\n\n🤔 其他运势：%d 分\n%s"+
		"\n\n🍀 幸运指南：\n"+
		"幸运颜色🌈：%s\n"+
		"幸运方位🛰️：%s\n"+
		"幸运数字🔢：%s\n"+
		"幸运食物🍜：%s\n"+
		"幸运饮品🥤：%s\n"+
		"幸运穿着👔：%s\n"+
		"幸运时间⌚️：%s\n\n"+
		"不宜星座👎：%s\n"+
		"贵人星座🎉：%s\n"+
		"桃花星座😍：%s\n"+
		"\n"+
		"📅 宜：%s\n"+
		"⛔ 忌：%s\n", todayYunShi.Xingzuo, todayYunShi.YunShiTime, todayYunShi.GeneralScore, todayYunShi.General, shi.Phase.Score, shi.Phase.YsContent, shi.Love.Score, shi.Love.YsContent, shi.Cause.Score, shi.Cause.YsContent, shi.Weather.Score, shi.Weather.YsContent,
		todayYunShi.LuckyColor, todayYunShi.LuckyPosition, todayYunShi.LuckyNum, todayYunShi.LuckyFood, todayYunShi.LuckyDrink, todayYunShi.LuckyDress, todayYunShi.LuckyTime, todayYunShi.XiaoRenXingZuo, todayYunShi.GuiRenXingZuo, todayYunShi.TaoHuaXingZuo, yiList, jiList)

}

func (starTalk) help() *proto.AbilityHelpInfo {
	return &proto.AbilityHelpInfo{
		Short:   "星座",
		Long:    "星座 输入水瓶座 返回星座说",
		Keyword: "星座",
	}
}

type YunShi struct {
	Xingzuo            string        `json:"xingzuo"`
	ConstellationIndex int           `json:"constellation_index"`
	XingzuoTime        string        `json:"xingzuo_time"`
	YunShiTime         string        `json:"yun_shi_time"`
	GeneralStar        int           `json:"general_star"`
	GeneralComment     string        `json:"general_comment"`
	EmotionStar        int           `json:"emotion_star"`
	CareerStar         int           `json:"career_star"`
	WealthStar         int           `json:"wealth_star"`
	HealthStar         int           `json:"health_star"`
	BargainingStar     int           `json:"bargaining_star"`
	GeneralScore       int           `json:"general_score"`
	EmotionScore       int           `json:"emotion_score"`
	CareerScore        int           `json:"career_score"`
	WealthScore        int           `json:"wealth_score"`
	HealthScore        int           `json:"health_score"`
	BargainingScore    int           `json:"bargaining_score"`
	XiaoRenXingZuo     string        `json:"xiao_ren_xing_zuo"`
	GuiRenXingZuo      string        `json:"gui_ren_xing_zuo"`
	TaoHuaXingZuo      string        `json:"tao_hua_xing_zuo"`
	LuckyColor         string        `json:"lucky_color"`
	LuckyPosition      string        `json:"lucky_position"`
	LuckyNum           string        `json:"lucky_num"`
	LuckyFood          string        `json:"lucky_food"`
	LuckyDrink         string        `json:"lucky_drink"`
	LuckyDress         string        `json:"lucky_dress"`
	LuckyTime          string        `json:"lucky_time"`
	HealthAdvice       string        `json:"health_advice"`
	DailyAdvice        string        `json:"daily_advice"`
	General            string        `json:"general"`
	Emotion            interface{}   `json:"emotion"`
	Career             interface{}   `json:"career"`
	Wealth             interface{}   `json:"wealth"`
	Yi                 string        `json:"yi"`
	Ji                 string        `json:"ji"`
	YiList             []string      `json:"yi_list"`
	JiList             []string      `json:"ji_list"`
	Tag                string        `json:"tag"`
	YunShi             []*YunShiInfo `json:"yun_shi"`
}

type YunShiInfo struct {
	Title string   `json:"title"`
	Dec   []string `json:"dec"`
}
