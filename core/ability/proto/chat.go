package proto

type ChatCacheInfo struct {
	AbilityName     string `json:"ability_name"`
	LastTime        int64  `json:"last_time"`
	CacheExpiredSec int    `json:"expired_sec"`
}

type AbilityHelpInfo struct {
	Short   string `json:"short"`
	Long    string `json:"long"`
	Keyword string `json:"keyword"`
}
