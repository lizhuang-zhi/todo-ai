package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitIMPlanRouter(Router *gin.RouterGroup) {
	router := Router.Group("im_plan").Use(api.AuthCheck())
	{
		router.POST("chat", WrapperHandler(api.ChatMessage))           // 聊天
		router.GET("messages", WrapperHandler(api.ChatHistoryMessage)) // 获取聊天历史
		router.POST("apply", WrapperHandler(api.DayApplyAiPlan))       // 应用Ai生成的计划
	}
}
