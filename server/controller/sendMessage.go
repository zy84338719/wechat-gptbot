package controller

import (
	"github.com/eatmoreapple/openwechat"
	"wechat-gptbot/core/global"

	"github.com/gin-gonic/gin"
	"net/http"
)

type SendMessageResponse struct {
	BaseResponse
}

type SendMessageRequest struct {
	ToUsername string          `json:"to_username"`
	ToNickname string          `json:"to_nickname"`
	Content    []string        `json:"content"`
	UserType   global.UserType `json:"user_type"`
}

func SendMessage(c *gin.Context, bot *openwechat.Bot) {
	request := SendMessageRequest{}
	err := c.ShouldBindBodyWithJSON(&request)
	response := SendMessageResponse{
		BaseResponse: BaseResponse{
			Code: 0,
			Msg:  "ok",
		},
	}
	defer c.JSON(http.StatusOK, response)
	if nil != err {
		response.Code = http.StatusNetworkAuthenticationRequired
		response.Msg = err.Error()
		return
	}
	data, err := getFriends(c, bot)
	if err != nil {
		response.Code = http.StatusNetworkAuthenticationRequired
		response.Msg = err.Error()
		return
	}
	user, err := bot.GetCurrentUser()
	if nil != err {
		response.Code = http.StatusNetworkAuthenticationRequired
		response.Msg = err.Error()
		return
	}
	switch request.UserType {
	case global.IsGroup:
		if request.ToUsername == "" {
			for _, g := range data.Groups {
				if g.Nickname == request.ToNickname {
					request.ToUsername = g.Username
					break
				}
			}
		}
		group, ok := global.AllowGroupUsernameMap[request.ToUsername]
		if !ok {
			response.Code = http.StatusNotFound
			response.Msg = "not found"
			return
		}
		for _, text := range request.Content {
			global.SendObejctChannel <- global.BuildSendObjectByGroup(text, group.Group)

		}
	case global.IsFriend:
		if request.ToUsername == "" {
			for _, u := range data.Users {
				if u.Nickname == request.ToNickname {
					request.ToUsername = u.Username
					break
				}
			}
		}
		friends, err := user.Friends()
		if err != nil {
			response.Code = http.StatusNetworkAuthenticationRequired
			response.Msg = err.Error()
			return
		}
		u := friends.GetByUsername(request.ToUsername)
		if u == nil {
			response.Code = http.StatusNotFound
			response.Msg = "not found"
			return
		}
		for _, text := range request.Content {
			global.SendObejctChannel <- global.BuildSendObjectByFriend(text, u)
		}
	case global.IsFriendHelper:
		for _, text := range request.Content {
			global.SendObejctChannel <- global.BuildSendObjectByFriendHelper(text)
		}
	default:
		response.Code = http.StatusNetworkAuthenticationRequired
		response.Msg = err.Error()
	}
}
