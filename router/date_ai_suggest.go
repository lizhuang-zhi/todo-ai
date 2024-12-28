package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitDateAiSuggestRouter(Router *gin.RouterGroup) {
	router := Router.Group("date_ai_suggest").Use(api.AuthCheck())
	{
		router.GET("data", WrapperHandler(api.DateAiSuggestion))        // 获取每日AI合理化建议
		router.POST("click", WrapperHandler(api.ClickDateAiSuggestion)) // 点击每日AI合理化建议(消除红点)
		// router.POST("apply", WrapperHandler(api.ApplyAiSuggestion)) // 应用每日AI合理化建议
	}
}
