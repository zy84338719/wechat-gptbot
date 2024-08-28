package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/yi-nology/common/utils/xstrings"
	"github.com/yi-nology/sdk/conf"
	"gorm.io/gorm"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/consts"
	"wechat-gptbot/core/ability"
	"wechat-gptbot/core/global"
	"wechat-gptbot/core/self_large_model/gpt"
	"wechat-gptbot/core/self_large_model/keyword"
	"wechat-gptbot/core/svc"
	"wechat-gptbot/data"
	"wechat-gptbot/data/proto"
	"wechat-gptbot/utils"
)

var Context *svc.ServiceContext

type MessageMatchDispatcher struct {
	*openwechat.MessageMatchDispatcher
	ctx *svc.ServiceContext
}

func NewMessageMatchDispatcher() *MessageMatchDispatcher {
	windows := []time.Duration{
		time.Minute,
		time.Hour,
		24 * time.Hour,
	}
	limits := []int{
		2,  // 1分钟请求2次
		5,  // 1小时请求10次
		20, // 1天请求40次
	}
	dispatcher := openwechat.NewMessageMatchDispatcher()
	self := &MessageMatchDispatcher{
		dispatcher,
		svc.NewServiceContext(),
	}
	limit := utils.NewRateLimiter("wb:limit", windows, limits)

	// 存储
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		isGroup, member := global.JugleGroup(message.FromUserName)
		if isGroup && member.AllowGroup.SaveMessage {
			return true
		}
		return false
	}, self.saveMessage)

	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		needReply, isImage := checkMessageType(message)
		if !message.IsSendByGroup() {
			return false
		}
		isGroup, member := global.JugleGroup(message.FromUserName)
		if !isGroup {
			return false
		}
		user, err := message.Bot().GetCurrentUser()
		if err != nil {
			return false
		}
		sender, err := message.SenderInGroup()
		if nil != err {
			logrus.Error(err.Error())
			return false
		}
		if strings.HasPrefix(message.Content, "签到") && member.AllowGroup.AllowSignIn {

			if sender.DisplayName == "" {
				global.BuildSendObjectByGroup("请尽快完善群昵称哦", member.Group)
				sender.DisplayName = sender.NickName
			}
			text := SignLogic(user.Uin, sender.UserName, sender.DisplayName, member)
			message.ReplyText(text)

			return false
		}
		if !message.IsSendBySelf() && needReply && member.AllowGroup.Reply {
			allow, err := limit.Allow(context.Background(), sender.UserName)
			if err != nil {
				logrus.Errorf("Allow Error: %+v", err)
			}
			if !allow {
				name := sender.DisplayName
				if name == "" {
					name = sender.NickName
				}
				message.ReplyText("@" + name + " QAQ 你的请求过于频繁哦，如果您需要更多服务，可以单独添加我好友，享受更多更好更有意思的功能哦")
				return false
			}
		}
		return !message.IsSendBySelf() && needReply && !isImage && member.AllowGroup.Reply
	}, self.groupText)

	// 注册文本函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		needReply, isImage := checkMessageType(message)
		return !message.IsSendByGroup() && !message.IsSendBySelf() && needReply && !isImage
	}, self.text)

	// 注册图片函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		needReply, isImage := checkMessageType(message)
		return needReply && isImage
	}, self.image)
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsTickledMe()
	}, self.trickMe)

	// 注册新人加群函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		if !message.IsJoinGroup() {
			return false
		}

		flag, _ := global.JugleGroup(message.FromUserName)
		if !flag {
			return false
		}
		//user, err := message.Bot().GetCurrentUser()
		//if err != nil {
		//	return false
		//}
		//if message.IsSystem() {
		//
		//}
		//if err = data.GroupUserInfoSingleton.Create(context.Background(), proto.GroupUserinfo{
		//	BotID:         user.Uin,
		//	GroupUsername: member.Group.UserName,
		//	Username:      sender.UserName,
		//	DisplayName:   sender.DisplayName,
		//	Nickname:      sender.NickName,
		//}); err != nil {
		//	logrus.Error(err)
		//}
		return true
	}, self.joinGroup)
	// 注册发送红包函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		if !message.IsSendRedPacket() {
			return false
		}

		flag, _ := global.JugleGroup(message.FromUserName)
		if !flag {
			return false
		}
		return true
	}, self.sendRedPacket)
	dispatcher.SetAsync(true)
	return self
}

func (dispatcher *MessageMatchDispatcher) saveMessage(message *openwechat.MessageContext) {
	//message.ToUserName
	//message.FromUserName
	//message.Content
	//message.CreateTime
	//打印上面几个字段
	sender, err := message.SenderInGroup()
	if nil != err {
		logrus.Error(err.Error())
		return
	}
	user, err := message.Bot().GetCurrentUser()
	if err != nil {
		logrus.Errorf("%+v", err)
		return
	}
	if user == nil {
		logrus.Error("user is nil")
		return
	}
	err = data.GroupMessageSingleton.Create(message.Context(), proto.GroupMessage{
		Model:         gorm.Model{CreatedAt: time.Now()},
		BotID:         user.Uin,
		GroupUsername: message.FromUserName,
		Username:      sender.UserName,
		MsgId:         message.MsgId,
		Content:       message.Content,
		MsgType:       message.MsgType,
		MsgCreateTime: message.CreateTime,
	})
	if err != nil {
		logrus.Error(err)
	}
}

func (dispatcher *MessageMatchDispatcher) trickMe(message *openwechat.MessageContext) {
	if !message.IsComeFromGroup() {
		text := getTextByUsername(defText.TrickMe)
		message.ReplyText(text.TrickMe)
		return
	}
	sender, _ := message.Sender()
	text := getTextByUsername(sender.UserName)
	message.ReplyText(text.TrickMe)
	return
}

func (dispatcher *MessageMatchDispatcher) joinGroup(message *openwechat.MessageContext) {
	sender, err := message.Sender()
	if err != nil {
		logrus.Errorf("加入群聊数据为空%+v", err)
		text := getTextByUsername("")
		message.ReplyText(text.JoinGroup)
		return
	}
	if sender == nil {
		text := getTextByUsername("")
		message.ReplyText(text.JoinGroup)
	}
	text := getTextByUsername(sender.UserName)
	message.ReplyText(text.JoinGroup)
	return
}

func (dispatcher *MessageMatchDispatcher) groupText(message *openwechat.MessageContext) {
	sender, err := message.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}
	if sender == nil {
		message.ReplyText(" QAQ 敬请期待 更多有意思的功能")
		return
	}
	isGroup, member := global.JugleGroup(message.FromUserName)
	if !isGroup {
		return
	}

	sender, err = message.SenderInGroup()
	if err != nil {
		logrus.Error(err.Error())
		message.ReplyText(" QAQ 敬请期待 更多有意思的功能 ")
	}
	name := sender.DisplayName
	if name == "" {
		name = sender.NickName
	}

	keyWord, err := keyword.GetKeyword(context.Background(), message.Content)
	if err != nil {
		logrus.Errorf("GetKeyword Error: %+v", err)
	}

	containAbility := ability.Contains(keyWord)
	keyList := strings.Split(keyWord, "-")
	if keyWord == "" || strings.Contains(keyWord, "不知道") || containAbility == nil || len(keyList) != 3 || keyList[0] != "技能" && member.AllowGroup.EnableListAbility == nil || !xstrings.Contains(member.AllowGroup.EnableListAbility, keyList[1]) {
		if !member.AllowGroup.AiReply {
			message.ReplyText("@" + sender.NickName + "\nQAQ\n该群聊能力暂未开放\n如需开启请联系群主\nPS❗️❗️❗️：加我私聊可以体验自然对话能力哦")
			return
		}
		key, err := conf.RedisClient.Get(context.Background(), "gpt_"+sender.UserName).Result()
		if err != redis.Nil && err != nil {
			message.ReplyText("QAQ服务出点问题，请稍后重试  \n问题反馈：https://txc.qq.com/products/658405")
			return
		}
		once.Do(func() {
			gpt.Init()
		})

		text, key, err := gpt.GPT(context.Background(), key, message.Content)
		if err != redis.Nil && err != nil {
			message.ReplyText("QAQ服务出点问题，请稍后重试  \n问题反馈：https://txc.qq.com/products/658405")
			return
		}

		_ = conf.RedisClient.Set(context.Background(), "gpt_"+sender.UserName, key, 3*24*time.Hour).Err()
		message.ReplyText("@" + name + "\n" + text)
		return
	}
	if member.AllowGroup.EnableAbility {
		message.ReplyText(containAbility.TextFunc(keyList[2]))
	}
	return
}

var once sync.Once

func (dispatcher *MessageMatchDispatcher) text(message *openwechat.MessageContext) {
	sender, err := message.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}
	if sender == nil {
		message.ReplyText(" QAQ 敬请期待 更多有意思的功能")
		return
	}

	//for _, reply := range Context.Session.Chat(context.WithValue(context.TODO(), "sender", sender.NickName), utils.BuildPersonalMessage(sender.NickName, message.Content)) {
	//	fmt.Printf("[text] Response: %s\n", reply) // 输出回复消息到日志
	//	_, err = message.ReplyText(reply)
	//	if err != nil {
	//		logrus.Infof("msg.ReplyText Error: %+v", err)
	//	}
	//}

	//once.Do(func() {
	//	gpt = self_large_model.NewSelfGpt(config.C.SelfGpt.BaseUrl, config.C.SelfGpt.GptAuthorization)
	//	profile, err := gpt.Profile(context.Background())
	//	if err != nil {
	//		return
	//	}
	//	Id = profile.Id
	//	token, err = gpt.ChatOpen(context.Background(), Id)
	//	if err != nil {
	//		return
	//	}
	//})

	//chatMessage, err := gpt.ChatMessage(context.Background(), token, self_large_model.ChatMessageRequest{
	//	Message: message.Content,
	//	ReChat:  true,
	//	Stream:  false,
	//})
	keyWord, err := keyword.GetKeyword(context.Background(), message.Content)
	if err != nil {
		logrus.Errorf("GetKeyword Error: %+v", err)
	}

	containAbility := ability.Contains(keyWord)
	keyList := strings.Split(keyWord, "-")
	if keyWord == "" || strings.Contains(keyWord, "不知道") || containAbility == nil || len(keyList) != 3 || keyList[0] != "技能" {
		key, err := conf.RedisClient.Get(context.Background(), "gpt_"+sender.UserName).Result()
		if err != redis.Nil && err != nil {
			message.ReplyText("QAQ服务出点问题，请稍后重试  \n问题反馈：https://txc.qq.com/products/658405")
			return
		}
		once.Do(func() {
			gpt.Init()
		})

		text, key, err := gpt.GPT(context.Background(), key, message.Content)
		if err != redis.Nil && err != nil {
			message.ReplyText("QAQ服务出点问题，请稍后重试  \n问题反馈：https://txc.qq.com/products/658405")
			return
		}
		_ = conf.RedisClient.Set(context.Background(), "gpt_"+sender.UserName, key, 3*24*time.Hour).Err()
		message.ReplyText(text)
		return
	}

	message.ReplyText(containAbility.TextFunc(keyList[2]))

}

func (dispatcher *MessageMatchDispatcher) sendRedPacket(message *openwechat.MessageContext) {
	sender, err := message.SenderInGroup()
	if err != nil {
		logrus.Error(err.Error())
		message.ReplyText(" QAQ 敬请期待 更多有意思的功能 ")
	}
	if sender == nil {
		return
	}
	name := sender.DisplayName
	if name == "" {
		name = sender.NickName
	}
	message.ReplyText("@" + name + " 发红包了🦷，大家快来抢🧧")
}

func (dispatcher *MessageMatchDispatcher) image(message *openwechat.MessageContext) {
	message.Content = strings.TrimLeft(message.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix)
	sender, err := message.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}

	if message.IsSendByGroup() {
		sender, err = message.SenderInGroup()
	}

	prompt := strings.TrimSpace(message.Content)
	uri := Context.Session.CreateImage(context.WithValue(context.TODO(), "sender", sender.NickName), prompt)
	if uri == "" {
		logrus.Infof("[image] Response: url 为空")
		message.ReplyText(consts.ErrTips)
		return
	}
	logrus.Infof("[image] Response: url = %s", uri)
	reader := bytes.Buffer{}
	err = utils.CompressImage(uri, &reader)
	if err != nil {
		logrus.Infof("[image] downloadImage err, err=%+v", err)
		message.ReplyText(consts.ErrTips)
		return
	}
	fu := message.ReplyImage
	if checkFile(uri) {
		fu = message.ReplyFile
	}
	_, err = fu(&reader)
	if err != nil {
		logrus.Infof("msg.ReplyImage Error: %+v", err)
	}
}

// 判断是否是发给我的消息
func checkMessageType(msg *openwechat.Message) (needReply bool, isImage bool) {
	// 如果包含了我的唤醒词
	msg.Content = strings.TrimLeft(msg.Content, " ")
	sender, err := msg.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}
	if !msg.IsText() {
		return false, false
	}
	if msg.IsSendBySelf() {
		return false, false
	}
	if !msg.IsSendByGroup() {
		// 私信消息
		// 私信消息不要管公众号消息
		if sender.IsMP() {
			return false, false
		}
		return true, checkCreateImage(msg)
	}
	//  如果是艾特我的消息
	if msg.IsAt() {
		prefix := fmt.Sprintf("@%s", msg.Owner().NickName)
		if strings.HasPrefix(msg.Content, prefix) {
			msg.Content = msg.Content[len(prefix):]
		}
		return true, checkCreateImage(msg)
	}

	if strings.HasPrefix(msg.Content, "星期五") {
		msg.Content = strings.TrimLeft(msg.Content, "星期五")
		return true, false
	}
	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.TextConfig.TriggerPrefix) {
		msg.Content = strings.TrimLeft(msg.Content, config.C.Base.Gpt.TextConfig.TriggerPrefix)
		return true, false
	}

	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix) {
		msg.Content = strings.TrimLeft(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix)
		return true, true
	}

	return false, false
}

// 通过语义判断是否是文生图的需求
func checkCreateImage(msg *openwechat.Message) bool {
	msg.Content = strings.TrimPrefix(msg.Content, "\u2005")
	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix) {
		return true
	}
	return false
}

func checkFile(uri string) bool {
	u, _ := url.Parse(uri)
	// 获取文件名
	name := path.Base(u.Path)
	return path.Ext(name) == ".webp"
}

func isInTimeRange(allowStartTime, allowEndTime string) (bool, error) {
	// 解析时间范围
	start, err := time.Parse("15:04", allowStartTime)
	if err != nil {
		return false, fmt.Errorf("解析开始时间出错: %v", err)
	}
	end, err := time.Parse("15:04", allowEndTime)
	if err != nil {
		return false, fmt.Errorf("解析结束时间出错: %v", err)
	}

	// 获取当前时间
	current := time.Now()
	// 只保留当前时间的小时和分钟
	currentTime := time.Date(0, 1, 1, current.Hour(), current.Minute(), 0, 0, time.UTC)

	// 判断当前时间是否在时间范围内
	if end.Before(start) {
		return currentTime.After(start) || currentTime.Before(end), nil
	}
	return currentTime.After(start) && currentTime.Before(end), nil
}
