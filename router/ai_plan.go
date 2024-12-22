package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitAiPlanRouter(Router *gin.RouterGroup) {
	router := Router.Group("ai_plan").Use(api.AuthCheck())
	{
		router.POST("day_gen", WrapperHandler(api.DayGenerateAiPlan)) // 生成AI当日规划
		router.POST("day_apply", WrapperHandler(api.DayApplyAiPlan))  // 应用AI当日规划
	}
}
