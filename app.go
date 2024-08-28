package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"wechat-gptbot/core"

	"wechat-gptbot/core/handler"
	"wechat-gptbot/server"
)

func main() {
	// 初始化核心配置
	core.Initialize()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 启动监听端口
	go server.NewApiServer(bot).Run()
	// 定义消息处理函数
	// 获取消息处理器
	dispatcher := handler.NewMessageMatchDispatcher()
	bot.MessageHandler = dispatcher.AsMessageHandler()
	bot.UUIDCallback = consoleQrCode // 注册登陆二维码回调
	// 登录回调
	//bot.SyncCheckCallback = nil
	reloadStorage := openwechat.NewFileHotReloadStorage("token/config.json")
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); nil != err {
		panic(err)
	}
	// 获取当前登录的用户
	self, err := bot.GetCurrentUser()
	if nil != err {
		panic(err)
	}
	logrus.Infof("login success %s,%s,%s", self.User, self.City, self.Province)

	go handler.SpciallyActionFunc(self)
	go logrus.Infof("login success %s,%s,%s", self.User, self.City, self.Province)
	bot.Block()
}

func consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	fmt.Println(q.ToSmallString(false))
}
