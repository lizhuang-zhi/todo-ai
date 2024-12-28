package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitProfileRouter(Router *gin.RouterGroup) {
	router := Router.Group("profile").Use(api.AuthCheck())
	{
		router.GET("data", WrapperHandler(api.ProfileData)) // 个人信息数据
	}
}
