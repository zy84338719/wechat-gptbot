package fix

import (
	"context"
	"github.com/sirupsen/logrus"
	"strconv"
	"wechat-gptbot/data"
	"wechat-gptbot/data/logic"
)

func FixLogic(ctx context.Context, id, version int64) (map[string]string, map[string]map[string]logic.NewUserData) {
	l := logic.Logic{}
	botId := strconv.Itoa(int(id))
	info := l.GetGroupsInfo(ctx, botId, version)
	groupUsernameMap := l.GetGroupUsernameMap(ctx, info)
	groupNameNewData := map[string]map[string]logic.NewUserData{}
	for name, group := range info {
		groupNameNewData[name] = l.GetOldUsername2NewUsername(ctx, group)
	}

	return groupUsernameMap, groupNameNewData
}

func FixMessage(ctx context.Context, groupUsernameMap map[string]string, groupNameNewData map[string]map[string]logic.NewUserData) {
	count, err := data.GroupMessageSingleton.Count(ctx)
	if err != nil {
		logrus.Error(err)
	}
	// 修复消息
	for i := 0; i < int(count/100+1); i++ {
		list, err := data.GroupMessageSingleton.List(ctx, i, 100)
		if err != nil {
			logrus.Error(err)
		}
		for _, m := range list {
			if groupUsername, ok := groupUsernameMap[m.GroupUsername]; ok {
				if userData, ok := groupNameNewData[groupUsername][m.Username]; ok {
					if m.Username == userData.NewUsername {
						continue
					}
					m.GroupUsername = userData.NewGroupUsername
					m.Username = userData.NewUsername
					err := data.GroupMessageSingleton.Update(ctx, m)
					if err != nil {
						logrus.Error(err)
					}

				}
			}
		}
	}
}

func FixSign(ctx context.Context, groupUsernameMap map[string]string, groupNameNewData map[string]map[string]logic.NewUserData) {
	// 修复签到
	count := data.GroupSignSingleton.Count(ctx)
	for i := 0; i < int(count/100+1); i++ {
		list, err := data.GroupSignSingleton.List(ctx, i, 100)
		if err != nil {
			logrus.Error(err)
		}
		for _, m := range list {
			if groupUsername, ok := groupUsernameMap[m.GroupUsername]; ok {
				if userData, ok := groupNameNewData[groupUsername][m.Username]; ok {
					if m.Username == userData.NewUsername {
						continue
					}
					m.GroupUsername = userData.NewGroupUsername
					m.Username = userData.NewUsername
					err := data.GroupSignSingleton.Update(ctx, m)
					if err != nil {
						logrus.Error(err)
					}
				}
			}
		}
	}
}
