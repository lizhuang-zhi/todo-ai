package router

import (
	"todo-ai/api"

	"github.com/gin-gonic/gin"
)

func InitCalendarRouter(Router *gin.RouterGroup) {
	router := Router.Group("calendar").Use(api.AuthCheck())
	{
		router.GET("data", WrapperHandler(api.CalendarData)) // 获取日历数据
	}
}
