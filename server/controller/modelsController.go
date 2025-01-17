package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"wechat-gptbot/config"
)

type CurrentModelResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data ModelInfo `json:"data"`
}

// CurrentModel 获取当前模型
func CurrentModel(c *gin.Context) {
	c.JSON(http.StatusOK, CurrentModelResponse{
		Code: 200,
		Msg:  "ok",
		Data: ModelInfo{config.C.GetBaseModel(), openai.CreateImageModelDallE3}})
}

type ModelInfo struct {
	TextModel    string `json:"text_model"`
	DrawingModel string `json:"drawing_model"`
}

func ResetModel(c *gin.Context) {
	info := ModelInfo{}
	c.ShouldBindBodyWithJSON(&info)
	if info.TextModel != "" {
		config.C.ResetBase(func(cfg *config.BaseConfig) {
			cfg.BaseModel = info.TextModel
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200,
		"msg": "ok"})
}
