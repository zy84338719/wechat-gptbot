package cache

import (
	"context"
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"github.com/yi-nology/common/utils/xjson"
	"github.com/yi-nology/sdk/conf"
	"time"
	"wechat-gptbot/model"
)

type Friend struct {
}

func (f *Friend) GetFriends(ctx context.Context, bot *openwechat.Bot, user *openwechat.Self) (*model.Data, error) {
	datas := model.Data{}
	dataStr, err := conf.RedisClient.Get(bot.Context(), "bot_"+user.UserName).Result()
	if err == nil {
		if err = xjson.UnmarshalFromString(dataStr, &datas); err != nil {
			logrus.Errorf("Failed to unmarshal data: %v", err)
			return nil, err
		}
		return &datas, nil
	}
	// 获取当前用户朋友
	friends, err := user.Friends(true)
	if err != nil {
		logrus.Errorf("Failed to get friends: %v", err)
		return nil, errors.New("Failed to get friends")
	}
	// 收集昵称
	for _, f := range friends {
		datas.Users = append(datas.Users, model.FriendsInfo{
			Username:   f.UserName,
			Nickname:   f.NickName,
			RemarkName: f.RemarkName,
			HeadImgUrl: f.HeadImgUrl,
		})
	}
	// 获取群聊
	groups, err := user.Groups(true)
	if err != nil {
		logrus.Errorf("Failed to get groups: %v", err)
		return nil, errors.New("Failed to get groups")
	}
	for _, gs := range groups {
		datas.Groups = append(datas.Groups, model.GroupsInfo{
			Username:   gs.UserName,
			Nickname:   gs.NickName,
			RemarkName: gs.RemarkName,
			HeadImgUrl: gs.HeadImgUrl,
		})
	}
	toString, _ := xjson.MarshalToString(datas)
	conf.RedisClient.Set(bot.Context(), "bot_"+user.UserName, toString, 24*time.Hour)
	return &datas, nil
}
