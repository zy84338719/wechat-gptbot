package config

import (
	"encoding/json"
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/yi-nology/sdk/conf"
	"log"
	"os"
	"sync"
)

var configFileYml = flag.String("f", "./config/config-dev.yaml", "the config file")

var (
	C      = &Config{RWMutex: new(sync.RWMutex)}
	Prompt string
)

const (
	defaultNewCron     = "0 30 7 1/1 * ?"
	defaultWeatherCron = "0 0 8 1/1 * ?"
	defaultSportCron   = "0 26 17 ? * ?"
	cronConfigFile     = "./config/cron.json"
	promptConfigFile   = "./config/prompt.conf"
	configFile         = "./config/config.json"
)

type Config struct {
	*sync.RWMutex
	Base           BaseConfig `json:"base"` // 基础配置
	Cron           CronConfig `json:"cron"` // 定时任务配置
	SelfGpt        SelfGpt    `json:"-"`
	CrawlerDomain  string     `json:"-"`
	FilmDomain     string     `json:"-"`
	MoyuDomain     string     `json:"-"`
	ShortUrlDomain string     `json:"-"`
}

type DefConfig struct {
	conf.Config
	SelfGpt        SelfGpt `json:"self_gpt" yaml:"self_gpt"`
	CrawlerDomain  string  `json:"crawler_domain" yaml:"crawler_domain"`
	FilmDomain     string  `json:"film_domain" yaml:"film_domain"`
	ShortUrlDomain string  `json:"shorturl_domain" yaml:"shorturl_domain"`
	MoyuDomain     string  `json:"moyu_domain" yaml:"moyu_domain"`
}

type SelfGpt struct {
	BaseUrl              string `json:"base_url" yaml:"base_url"`
	KeyWordAuthorization string `json:"keyword_authorization" yaml:"keyword_authorization"`
	GptAuthorization     string `json:"gpt_authorization" yaml:"gpt_authorization"`
}

// ResetBase 修改基础配置文件
func (c *Config) ResetBase(hand func(cfg *BaseConfig)) {
	c.Lock()
	defer func() {
		c.Base.BackUp()
		c.Unlock()
	}()
	hand(&c.Base)
}

func (c *Config) ResetCron(hand func(cfg *CronConfig)) {
	c.Lock()
	defer func() {
		c.Cron.BackUp()
		c.Unlock()
	}()
	hand(&c.Cron)
}

type BaseConfig struct {
	Gpt struct {
		TextConfig  AuthConfig `json:"text_config"`
		ImageConfig AuthConfig `json:"image_config"`
	} `json:"gpt"`
	ContextStatus  bool       `json:"context_status"`
	BaseModel      string     `json:"base_model"`
	KeepaliveRobot string     `json:"keepalive_robot"`
	CronConfig     CronConfig `json:"cron_config"`
}

func (b *BaseConfig) IsValid() bool {

	authConfigs := []AuthConfig{
		b.Gpt.TextConfig,
		b.Gpt.ImageConfig,
	}

	for _, authConfig := range authConfigs {
		if authConfig.AuthToken == "" || authConfig.TriggerPrefix == "" {
			return false
		}
	}
	return true
}

// BackUp 备份基础文件
func (b *BaseConfig) BackUp() {
	err := writeFile(configFile, b)
	if err != nil {
		logrus.Errorf("备份定时任务文件失败 ,err = %s", err.Error())
	}
}

type CronConfig struct {
	WeatherConfig WeatherCronConfig     `json:"weather_config"`
	NewsConfig    NewsCronConfig        `json:"news_config"`
	SportConfig   WechatSportCronConfig `json:"sport_config"`
}

func (c *CronConfig) IsValid() bool {
	if c.WeatherConfig.Spec == "" {
		c.WeatherConfig.Spec = defaultWeatherCron
	}

	if c.NewsConfig.Spec == "" {
		c.NewsConfig.Spec = defaultNewCron
	}
	if c.SportConfig.Spec == "" {
		c.SportConfig.Spec = defaultSportCron
	}
	return true
}

// BackUp 备份定时任务文件
func (c *CronConfig) BackUp() {
	err := writeFile(cronConfigFile, c)
	if err != nil {
		logrus.Errorf("备份定时任务文件失败 ,err = %s", err.Error())
	}
}

type WechatSportCronConfig struct {
	Users []SportAccount `json:"users"`
	Spec  string         `json:"spec"` // cron 表达式
	Desc  string
}
type SportAccount struct {
	Name    string `json:"name"`    // 绑定微信名
	Account string `json:"account"` // 账号
	Pwd     string `json:"pwd"`     // 密码
	Min     int    `json:"min"`     //  最少步数
	Max     int    `json:"max"`     // 最多步数
}

// WeatherCronConfig 天气预报定时任务配置
type WeatherCronConfig struct {
	Users []struct {
		Name string `json:"name"` // 用户名
		City string `json:"city"` // 城市
	} `json:"users"`
	Spec string `json:"spec"` // cron 表达式
	Desc string
}

type NewsCronConfig struct {
	Users     []string // 用户名
	Groups    []string // 群名称
	Spec      string   `json:"spec"`       // cron 表达式
	RssSource string   `json:"rss_source"` // rss 推送源
	TopN      int      `json:"top_n"`      // 限制接收量
	Desc      string
}

func writeFile(filePath string, data interface{}) error {
	// 创建或打开文件 进行文件覆盖
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// 确保文件在函数结束时被关闭
	defer file.Close()

	// 创建一个 JSON 编码器，并设置编码的目标为文件
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置 JSON 格式化输出

	// 编码数据到 JSON，并写入文件
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) GetBaseModel() string {
	c.RLock()
	defer c.RUnlock()
	return c.Base.BaseModel
}

type AuthConfig struct {
	ProxyUrl      string `json:"proxy_url"` //代理地址，不填使用官方地址
	AuthToken     string `json:"auth_token"`
	TriggerPrefix string `json:"trigger_prefix"`
}

func (c *Config) IsValid() bool {

	authConfigs := []AuthConfig{
		c.Base.Gpt.TextConfig,
		c.Base.Gpt.ImageConfig,
	}

	for _, authConfig := range authConfigs {
		if authConfig.AuthToken == "" || authConfig.TriggerPrefix == "" {
			return false
		}
	}
	return true
}

func (c *Config) CheckCronValid() {
	if c.Cron.WeatherConfig.Spec == "" {
		c.Cron.WeatherConfig.Spec = defaultWeatherCron
	}

	if c.Cron.NewsConfig.Spec == "" {
		c.Cron.NewsConfig.Spec = defaultNewCron
	}
	if c.Cron.SportConfig.Spec == "" {
		c.Cron.SportConfig.Spec = defaultSportCron
	}
}

func InitConfig() {
	// 1. 读取 `config.json`
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `/config/config.json` 的配置, 错误信息: %+v\n", err)
	}
	var baseConfig BaseConfig
	if err = json.Unmarshal(data, &baseConfig); err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `config.json` 的格式, 错误信息: %+v\n", err)
	}
	if !baseConfig.IsValid() {
		log.Fatal("配置文件校验失败，请检查 `config.json`")
	}

	defConfig := DefConfig{}
	flag.Parse()
	println(*configFileYml)
	conf.MustLoad(*configFileYml, &defConfig)

	err = defConfig.Init()
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `config.yaml` 的配置, 错误信息: %+v\n", err)
	}
	C.Base = baseConfig
	// 2. 读取 prompt.txt
	prompt, err := os.ReadFile(promptConfigFile)
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `prompt.conf` 的配置, 错误信息: %+v\n", err)
	}
	Prompt = string(prompt)

	// 读取定时任务配置表
	cronConfigData, err := os.ReadFile(cronConfigFile)
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `cron.json` 的配置, 错误信息: %+v\n", err)
	}
	var cronConfig CronConfig
	err = json.Unmarshal(cronConfigData, &cronConfig)
	if err != nil {
		log.Fatalf("读取定时任务配置文件失败， 错误信息 %+v\n", err)
	}
	cronConfig.IsValid()
	C.Cron = cronConfig
	C.SelfGpt = defConfig.SelfGpt
	C.CrawlerDomain = defConfig.CrawlerDomain
	C.FilmDomain = defConfig.FilmDomain
	C.ShortUrlDomain = defConfig.ShortUrlDomain
	C.MoyuDomain = defConfig.MoyuDomain
}
