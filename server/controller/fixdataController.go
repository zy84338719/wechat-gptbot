package controller

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"wechat-gptbot/core/fix"
)

type FixResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type FixResquest struct {
	Version    int64 `json:"version"`
	FixSign    bool  `json:"fix_sign"`
	FixMessage bool  `json:"fix_message"`
}

// GetFriends 获取当前微信所有朋友群组关系
func FixData(c *gin.Context, bot *openwechat.Bot) {
	response := FixResponse{
		Code: http.StatusOK,
		Msg:  "ok",
	}
	resquest := FixResquest{}
	if err := c.ShouldBindJSON(&resquest); err != nil {
		response.Code = http.StatusBadRequest
		response.Msg = err.Error()
		return
	}

	user, err := bot.GetCurrentUser()
	if err != nil {
		logrus.Errorf("Failed to get current user: %v", err)
		return
	}
	groupUsernameMap, groupNameNewData := fix.FixLogic(c, user.Uin, resquest.Version)
	if resquest.FixMessage {
		fix.FixMessage(c, groupUsernameMap, groupNameNewData)
	}
	if resquest.FixSign {
		fix.FixSign(c, groupUsernameMap, groupNameNewData)
	}
	c.JSON(http.StatusOK, response)

}
