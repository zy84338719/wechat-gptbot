package logic

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"time"
)

type Message struct {
}
type MessageData struct {
	GroupUsername string `json:"group_username" gorm:"column:group_username"`
	GroupNickname string `json:"group_nickname" gorm:"column:group_nickname"`
	Username      string `json:"username" gorm:"column:username"`
	DisplayName   string `json:"display_name" gorm:"column:display_name"`
	Nickname      string `json:"nickname" gorm:"column:nickname"`
	MessageCount  int64  `json:"message_count" gorm:"column:message_count"`
}

func (Message) Top(ctx context.Context, startTime, endTime time.Time) map[string][]MessageData {
	sql := "WITH RankedMessages AS (SELECT gm.group_username, gm.username, COUNT(*) AS message_count, ROW_NUMBER() OVER (PARTITION BY gm.group_username ORDER BY COUNT(*) DESC) AS rank FROM group_message gm WHERE gm.deleted_at IS NULL AND gm.created_at BETWEEN ? AND ? GROUP BY gm.group_username, gm.username) SELECT rm.group_username, g.group_nickname, rm.username, gu.display_name, gu.nickname, rm.message_count FROM RankedMessages rm JOIN `group` g ON rm.group_username = g.group_username LEFT JOIN group_userinfo gu ON rm.group_username = gu.group_username AND rm.username = gu.username ORDER BY rm.group_username, rm.rank;"

	data := []MessageData{}
	if err := conf.MysqlClient.WithContext(ctx).Raw(sql, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05")).Scan(&data).Error; err != nil {
		return nil
	}

	result := map[string][]MessageData{}
	for _, v := range data {
		result[v.GroupUsername] = append(result[v.GroupUsername], v)
	}
	return result
}
