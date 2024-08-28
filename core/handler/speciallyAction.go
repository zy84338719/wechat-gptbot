package handler

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
	"wechat-gptbot/core/global"
	"wechat-gptbot/data"
	"wechat-gptbot/data/logic"
	"wechat-gptbot/data/proto"
	"wechat-gptbot/utils"
)

func getTextByUsername(getUsername string) global.TextDefine {
	userGroup, ok := global.AllowGroupUsernameMap[getUsername]
	if ok {
		return *userGroup.TextDefine
	}
	return defText
}

func updateGroupUserInfoDB(ctx context.Context, bot *openwechat.Self, membersMap map[string]global.GroupMember) {
	for _, member := range membersMap {
		dataList, err := data.GroupUserInfoSingleton.GetByGroupUsername(context.Background(), member.Group.UserName)
		if err != nil {
			continue
		}
		if dataList == nil {
			dataList = []proto.GroupUserinfo{}
			for _, groupMember := range member.Members {
				dataList = append(dataList, proto.GroupUserinfo{
					BotID:         bot.Uin,
					GroupUsername: member.Group.UserName,
					Username:      groupMember.UserName,
					DisplayName:   groupMember.DisplayName,
					Nickname:      groupMember.NickName,
				})
			}

			if err = data.GroupUserInfoSingleton.BatchCreate(context.Background(), dataList); err != nil {
				logrus.Error(err)
			}
			continue
		}
		for _, userinfo := range dataList {
			newData, ok := member.Members[userinfo.Username]
			if ok {
				if userinfo.DisplayName != newData.DisplayName || userinfo.Nickname != newData.NickName {
					userinfo.DisplayName = newData.DisplayName
					userinfo.Nickname = newData.NickName
					err := data.GroupUserInfoSingleton.Update(context.Background(), userinfo)
					if err != nil {
						logrus.Error(err.Error())
					}
				}
			}
		}
	}
}

func updateGroupDB(ctx context.Context, bot *openwechat.Self) {
	version := time.Now().Unix()
	for _, g := range global.AllowGroupUsernameMap {
		group, err := data.GroupSingleton.GetByUsername(ctx, bot.Uin, g.Group.UserName)
		if err != nil {
			logrus.Errorf("get group by username error:%v", err)
			continue
		}
		if group != nil {
			if group.GroupNickname != g.Group.NickName {
				group.GroupNickname = g.Group.NickName
				err := data.GroupSingleton.Update(ctx, *group)
				if err != nil {
					logrus.Error(err)
				}
			}
			continue
		}
		err = data.GroupSingleton.Create(ctx, proto.Group{
			BotID:         bot.Uin,
			BotName:       bot.NickName,
			GroupUsername: g.Group.UserName,
			GroupNickname: g.Group.NickName,
			Version:       version,
		})
		if err != nil {
			logrus.Errorf("create group error:%v", err)
		}
	}
}

func SpciallyActionFunc(bot *openwechat.Self) {
	global.AllowGroupUsernameMap = groupAction{}.getGroupMembersByUserName(bot)
	updateGroupDB(context.Background(), bot)
	updateGroupUserInfoDB(context.Background(), bot, global.AllowGroupUsernameMap)

	ticker := time.NewTicker(time.Second * 123)
	defer ticker.Stop()
	go sendTextToGroup(bot)

	userMap := map[string]map[string]UserInfo{}
	for range ticker.C {
		heartBeat(bot)
		groupActionMap := groupAction{}
		members := groupActionMap.getGroupMembersByUserName(bot)
		for _, member := range members {
			if member.AllowGroup.Redisplay {
				groupActionMap.redisplayCheck(member, userMap[member.Username])
			}
			if member.AllowGroup.ExitCheck {
				groupActionMap.exitGroupCheck(member, userMap[member.Username])
			}
			userMap[member.Username] = groupActionMap.updateUsernameBaseData(member.Username, member.Members)
		}
		if len(members) != len(global.AllowGroupUsernameMap) {
			updateGroupUserInfoDB(context.Background(), bot, members)
		}
		updateGroupDB(context.Background(), bot)
		global.AllowGroupUsernameMap = members
	}
}

func sendTextToGroup(bot *openwechat.Self) {
	for obejct := range global.SendObejctChannel {
		switch obejct.UserType {
		case global.IsFriend:
			bot.SendTextToFriend(obejct.SendFriend, obejct.Content)
		case global.IsGroup:
			bot.SendTextToGroup(obejct.SendGroup, obejct.Content)
		case global.IsFriendHelper:
			bot.SendTextToFriend(openwechat.NewFriendHelper(bot), obejct.Content)
		}
		time.Sleep(time.Second * 1)
	}
}

func heartBeat(bot *openwechat.Self) {
	// 向文件传输助手发送消息，不要再关注公众号了
	// 生成要发送的消息
	outMessage := fmt.Sprintf("防退出登录[%d]", utils.GetRandInt64(3000))
	global.SendObejctChannel <- global.BuildSendObjectByFriendHelper(outMessage)
}

type UserInfo struct {
	NickName    string
	DisplayName string
}
type groupAction struct {
}

func (g groupAction) exitGroupCheck(member global.GroupMember, userMap map[string]UserInfo) {
	for oldUserName, oldUser := range userMap {
		if _, ok := member.Members[oldUserName]; !ok {
			if oldUser.DisplayName == "" {
				oldUser.DisplayName = oldUser.NickName
			}
			text := getTextByUsername(member.Username)
			global.SendObejctChannel <- global.BuildSendObjectByGroup(fmt.Sprintf(text.ExitCheck, oldUser.NickName, oldUser.DisplayName), member.Group)
			data.GroupUserInfoSingleton.DeleteByGroupUsername(context.Background(), member.Username, oldUserName)
		}
	}
}

func (g groupAction) updateUsernameBaseData(username string, members map[string]*openwechat.User) map[string]UserInfo {
	membersMap := map[string]UserInfo{}
	for _, groupMember := range members {
		membersMap[groupMember.UserName] = UserInfo{
			NickName:    groupMember.NickName,
			DisplayName: groupMember.DisplayName,
		}
	}
	return membersMap
}

func (g groupAction) redisplayCheck(member global.GroupMember, userMap map[string]UserInfo) {
	for oldUserName, oldUser := range userMap {
		if newMember, ok := member.Members[oldUserName]; ok {
			if newMember.DisplayName != oldUser.DisplayName {
				if oldUser.DisplayName == "" {
					oldUser.DisplayName = oldUser.NickName
				}
				text := getTextByUsername(member.Username)
				global.SendObejctChannel <- global.BuildSendObjectByGroup(fmt.Sprintf(text.Redisplay, newMember.NickName, oldUser.DisplayName, newMember.DisplayName), member.Group)
			}
		}
	}
}

func SignLogic(botId int64, username, displayName string, member global.GroupMember) string {
	isSiginTime, err := isInTimeRange(member.SignedIn.AllowStartTime, member.SignedIn.AllowEndTime)
	if err != nil {
		logrus.Error(err.Error())
		return ""
	}
	if !isSiginTime {
		return fmt.Sprintf(member.SignedIn.SignInFailTimeout, displayName, member.SignedIn.AllowStartTime, member.SignedIn.AllowEndTime)
	}

	userSign, err := data.GroupSignSingleton.GetByGroupUsernameAndNowDay(member.Group.UserName, username)
	if err != nil {
		return fmt.Sprintf(member.SignedIn.SignInError, displayName)
	}
	if userSign != nil {
		return fmt.Sprintf(member.SignedIn.SignInFailSignedIn, displayName)
	}

	dayTopCount, err := data.GroupSignSingleton.GetByGroupUsernameAndDayCount(member.Group.UserName)
	if err != nil {
		return fmt.Sprintf(member.SignedIn.SignInError, displayName)
	}
	monthCount, err := data.GroupSignSingleton.GetByGroupUsernameAndMonthCount(username)
	if err != nil {
		return fmt.Sprintf(member.SignedIn.SignInError, displayName)
	}
	err = data.GroupSignSingleton.Create(proto.GroupSign{
		BotID:         botId,
		GroupUsername: member.Username,
		Username:      username,
		SignTime:      time.Now(),
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	})
	if err != nil {
		return fmt.Sprintf(member.SignedIn.SignInError, displayName)
	}
	rank := int64(0)
	top, err := logic.TodaySignTop(context.Background(), member.Group.UserName, username)
	if err != nil {
		logrus.Errorf("%+v", err)
	}
	if top.Rank == 0 {
		rank = 1
	} else {
		rank = top.Rank
	}
	return fmt.Sprintf(member.SignedIn.SignInSuccess, displayName, dayTopCount+1, time.Now().Format("15:04:05"), monthCount+1, rank)
}
