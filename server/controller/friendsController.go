package controller

import (
	"context"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"wechat-gptbot/cache"
	"wechat-gptbot/model"
)

type FriendsResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data model.Data `json:"data"`
}

// GetFriends 获取当前微信所有朋友群组关系
func GetFriends(c *gin.Context, bot *openwechat.Bot) {
	response := FriendsResponse{
		Code: http.StatusOK,
		Msg:  "ok",
	}
	data, err := getFriends(c, bot)
	// 还没有登录
	if nil != err {
		// 如果未登录，返回错误信息
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{
			"code": http.StatusNetworkAuthenticationRequired,
			"msg":  "User not authenticated",
		})
		return
	}
	response.Data = *data
	c.JSON(http.StatusOK, response)

}

func getFriends(ctx context.Context, bot *openwechat.Bot) (*model.Data, error) {
	// 获取当前用户
	user, err := bot.GetCurrentUser()
	if err != nil {
		logrus.Errorf("Failed to get current user: %v", err)
		return nil, err
	}
	friend := cache.Friend{}
	friends, err := friend.GetFriends(ctx, bot, user)
	if err != nil {
		return nil, err
	}
	return friends, nil
}
