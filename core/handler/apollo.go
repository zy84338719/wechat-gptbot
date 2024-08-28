package handler

import (
	"context"
	"github.com/eatmoreapple/openwechat"
	"github.com/yi-nology/common/utils/xjson"
	"github.com/yi-nology/sdk/conf"
	"wechat-gptbot/core/global"
)

func getApolloConfigList(ctx context.Context) ([]global.AllowGroup, error) {
	config := conf.Apollo.GetConfig("tempConfig.json")
	configAllowGroupStr := config.GetStringValue("content", "[]")
	allowGroupList := make([]global.AllowGroup, 0)
	err := xjson.UnmarshalFromString(configAllowGroupStr, &allowGroupList)
	if err != nil {
		return nil, err
	}
	return allowGroupList, nil
}

var defSignedIn = global.SignedIn{
	AllowStartTime:     "6:30",
	AllowEndTime:       "11:00",
	SignInSuccess:      "%s 签到成功🎉🎉🎉\n你是第 %d 个签到的人\n签到时间：%s\n本月你已经累计签到 %d 天\n暂列第 %d 名",
	SignInFailSignedIn: "%s 签到失败🥹，你已经签到过了",
	SignInFailTimeout:  "%s 签到失败🥹，签到时间已经过了，签到时间：%s-%s",
	SignInError:        "%s 签到失败🥹，出现了一些问题，你可以联系群管理员",
}

var defText = global.TextDefine{
	Redisplay: "号外号外，大家注意哦⚠️ 微信名称【 %s】, 由 %s 改为 %s 快来问问Ta，为什么改名吧！",
	ExitCheck: "呜呜呜 🥹 【%s】 群昵称：%s 已经永久的离开了这个群",
	TrickMe:   "QAQ 🥹❗️ 我还很弱小，现在什么都不会。你先不要生气",
	JoinGroup: "Hi👋 欢迎新同学，大伙儿快出来🦷，有新人加入哦➕",
	Error:     "呜呜呜 🥹 出现了一些问题，我还不会处理，你可以联系群管理员",
}

func (groupAction) getGroupMembersByUserName(bot *openwechat.Self) map[string]global.GroupMember {
	groups, err := bot.Groups()
	if err != nil {
		return nil
	}
	allowGroupList, err := getApolloConfigList(context.Background())
	if err != nil {
		return nil
	}

	groupMemberMap := make(map[string]global.GroupMember, 0)
	for _, allowGroup := range allowGroupList {
		var group *openwechat.Group
		if allowGroup.Username != "" {
			group = groups.GetByUsername(allowGroup.Username)
		} else {
			group = groups.GetByNickName(allowGroup.GroupName)
		}
		if group == nil {
			continue
		}
		members, err := group.Members()
		if err != nil {
			continue
		}
		membersMap := make(map[string]*openwechat.User, 0)

		for _, member := range members {
			membersMap[member.UserName] = member
		}
		if allowGroup.IsTextDefault || allowGroup.TextDefine == nil {
			allowGroup.TextDefine = &defText
		}
		if allowGroup.SignedIn == nil {
			allowGroup.SignedIn = &defSignedIn
		}
		groupMemberMap[group.UserName] = global.GroupMember{
			Username:   group.UserName,
			GroupName:  group.NickName,
			Members:    membersMap,
			Group:      group,
			AllowGroup: allowGroup,
			TextDefine: allowGroup.TextDefine,
			SignedIn:   allowGroup.SignedIn,
		}

	}

	return groupMemberMap
}
