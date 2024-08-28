package routes

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"wechat-gptbot/server/controller"
)

func InitRoute(bot *openwechat.Bot) *gin.Engine {
	router := gin.New()
	r := router.Group("/wechat-gptbot")
	{
		r.GET("/checklogin", func(context *gin.Context) {
			controller.CheckLogin(context, bot)
		})
		r.GET("/current-model", controller.CurrentModel)
		r.POST("/reset-model", controller.ResetModel)
		r.GET("/friends", func(context *gin.Context) {
			controller.GetFriends(context, bot)
		})

		r.POST("/send-message", func(context *gin.Context) {
			controller.SendMessage(context, bot)
		})
		r.POST("/fixdata", func(context *gin.Context) {
			controller.FixData(context, bot)
		})
	}
	return router
}
