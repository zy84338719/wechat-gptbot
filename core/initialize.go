package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"wechat-gptbot/config"
	"wechat-gptbot/core/handler"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/news"
	"wechat-gptbot/core/plugins/weather"
	"wechat-gptbot/core/self_large_model/keyword"
	"wechat-gptbot/core/svc"
	"wechat-gptbot/logger"
	"wechat-gptbot/streamlit_app"
)

func Initialize() {
	// 初始化日志
	logger.InitLogrus(logger.Config{
		Level:      logrus.DebugLevel,
		ObjectName: "wechat-gptbot",
		WriteFile:  false,
	})
	// 初始化配置文件
	config.InitConfig()
	// 初始化插件
	plugins.Manger.Register(weather.NewPlugin(),
		news.NewPlugin(news.SetRssSource(config.C.Cron.NewsConfig.RssSource),
			news.SetTopN(config.C.Cron.NewsConfig.TopN)))

	// 初始化会话上下文管理器
	handler.Context = svc.NewServiceContext()
	// 启动streamlit
	go streamlit_app.RunStreamlit()
	keyword.Init(context.Background())
}
