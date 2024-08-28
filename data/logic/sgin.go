package logic

import (
	"context"
	"fmt"
	"github.com/yi-nology/common/utils/xnow"
	"github.com/yi-nology/sdk/conf"
	"time"
)

type Sign struct {
}

var SignSingleton = Sign{}

type SignData struct {
	GroupUsername string `json:"group_username" gorm:"column:group_username"`
	GroupNickname string `json:"group_nickname" gorm:"column:group_nickname"`
	Username      string `json:"username" gorm:"column:username"`
	DisplayName   string `json:"display_name" gorm:"column:display_name"`
	SignCount     int    `json:"sign_count" gorm:"column:sign_count"`
	Nickname      string `json:"nickname" gorm:"column:nickname"`
}

func (receiver Sign) Top(ctx context.Context, count int, startTime, endTime time.Time) map[string][]SignData {
	sql := "WITH RankedSignIns AS (SELECT gs.group_username, gs.username, COUNT(*) AS sign_count, ROW_NUMBER() OVER (PARTITION BY gs.group_username ORDER BY COUNT(*) DESC) AS rank FROM `group_sign` gs WHERE gs.sign_time BETWEEN ? AND ? GROUP BY gs.group_username, gs.username) SELECT rs.group_username, g.group_nickname, rs.username, gu.display_name, gu.nickname, rs.sign_count FROM RankedSignIns rs LEFT JOIN `group_userinfo` gu ON rs.group_username = gu.group_username AND rs.username = gu.username LEFT JOIN `group` g ON rs.group_username = g.group_username WHERE rs.rank <= ?;"
	result := map[string][]SignData{}
	data := []SignData{}
	conf.MysqlClient.WithContext(ctx).Raw(sql, startTime.Format("2006-01-02 15:04:06"), endTime.Format("2006-01-02 15:04:06"), count).Scan(&data)
	for _, v := range data {
		result[v.GroupUsername] = append(result[v.GroupUsername], v)
	}
	return result
}

type TodyData struct {
	Username      string `json:"username" gorm:"column:username"`
	GroupUsername string `json:"group_username" gorm:"column:group_username"`
	SignCount     int64  `json:"sign_count" gorm:"column:sign_count"`
	Rank          int64  `json:"rank" gorm:"column:rank"`
}

func TodaySignTop(ctx context.Context, group_username, username string) (TodyData, error) {
	sql := "WITH SignInRanks AS (SELECT gs.group_username, gs.username, COUNT(*) AS sign_count, ROW_NUMBER() OVER (PARTITION BY gs.group_username ORDER BY COUNT(*) DESC) AS rank FROM `group_sign` gs WHERE gs.deleted_at IS NULL AND gs.sign_time BETWEEN '%s' AND '%s' GROUP BY gs.group_username, gs.username) SELECT sr.group_username, sr.username, sr.sign_count, sr.rank FROM SignInRanks sr WHERE sr.group_username = '%s' AND sr.username = '%s' ORDER BY sr.rank;"
	data := TodyData{}
	query := fmt.Sprintf(sql, xnow.BeginningOfMonth().Format("2006-01-02 15:04:05"), xnow.EndOfMonth().Format("2006-01-02 15:04:05"), group_username, username)
	tx := conf.MysqlClient.WithContext(ctx).Raw(query).Scan(&data)
	if tx.Error != nil {
		return data, tx.Error
	}
	return data, nil
}
