package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitTaskRouter(Router *gin.RouterGroup) {
	router := Router.Group("task").Use(api.AuthCheck())
	{
		router.POST("create", WrapperHandler(api.CreateTask)) // 创建任务
	}
}
