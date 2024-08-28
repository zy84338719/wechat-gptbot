package logic

import (
	"context"
	"wechat-gptbot/data"
	"wechat-gptbot/data/proto"
)

type Logic struct {
}

// 获得最新的group信息
func (Logic) GetGroupNewInfo(ctx context.Context, gs []proto.Group) proto.Group {
	group := proto.Group{}
	for _, g := range gs {
		if g.CreatedAt.Sub(group.CreatedAt) > 0 {
			group = g
		}
	}
	return group
}

func (Logic) GetNewGroupUserInfoByUserName(ctx context.Context, group proto.Group, userInfo []proto.GroupUserinfo) *proto.GroupUserinfo {
	for _, u := range userInfo {
		if u.GroupUsername == group.GroupUsername {
			return &u
		}
	}
	return nil
}

func (l Logic) GetGroupsInfo(ctx context.Context, botId string, version int64) map[string][]GroupInfo {
	gs, err := data.GroupSingleton.GetByBotId(ctx, botId)
	if err != nil {
		return nil
	}
	if gs == nil {
		return nil
	}
	m := map[string][]GroupInfo{}
	flag := false
	for _, g := range gs {
		usernameList, err := data.GroupUserInfoSingleton.GetByGroupUsername(ctx, g.GroupUsername)
		if err != nil {
			return nil
		}
		if g.Version == version {
			flag = true
		} else {
			flag = false
		}
		m[g.GroupNickname] = append(m[g.GroupNickname], GroupInfo{
			Groupname:    g.GroupNickname,
			Username:     g.GroupUsername,
			UserInfoList: usernameList,
			IsNew:        flag,
		})
	}
	return m
}

// 获得群组名称和群组username
func (Logic) GetGroupUsernameMap(ctx context.Context, gi map[string][]GroupInfo) map[string]string {
	m := map[string]string{}
	for name, group := range gi {
		for _, n := range group {
			m[n.Username] = name
		}
	}
	return m
}

type NewUserData struct {
	NewGroupUsername string
	NewUsername      string
}

func (l Logic) GetOldUsername2NewUsername(ctx context.Context, groups []GroupInfo) map[string]NewUserData {
	if groups == nil {
		return nil
	}
	newInfo := GroupInfo{}
	for _, g := range groups {
		if g.IsNew {
			newInfo = g
		}
	}
	m := map[string]NewUserData{}
	for _, info := range groups {
		for _, uOld := range info.UserInfoList {
			for _, uNew := range newInfo.UserInfoList {
				if uOld.Nickname == uNew.Nickname {
					m[uOld.Username] = NewUserData{
						NewGroupUsername: uNew.GroupUsername,
						NewUsername:      uNew.Username,
					}
					break
				}
			}
		}

	}

	return m
}

func (Logic) GetGroupUserInfoByUserNameMap(ctx context.Context) map[string]proto.GroupUserinfo {
	usernameList, err := data.GroupUserInfoSingleton.List(ctx, 0, 100000)
	if err != nil {
		return nil
	}
	m := map[string]proto.GroupUserinfo{}
	for _, u := range usernameList {
		m[u.Username] = u
	}
	return m
}

type GroupInfo struct {
	Groupname    string
	Username     string
	UserInfoList []proto.GroupUserinfo
	IsNew        bool
}
