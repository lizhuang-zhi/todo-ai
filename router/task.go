package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitTaskRouter(Router *gin.RouterGroup) {
	router := Router.Group("task").Use(api.AuthCheck())
	{
		router.POST("create", WrapperHandler(api.CreateTask))     // 创建任务
		router.POST("update", WrapperHandler(api.UpdateTask))     // 修改任务
		router.POST("delete", WrapperHandler(api.DeleteTask))     // 删除任务
		router.POST("finished", WrapperHandler(api.FinishedTask)) // 完成任务
		router.GET("list", WrapperHandler(api.ListTask))          // 任务列表
	}
}
