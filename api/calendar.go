package api

import (
	"errors"
	"fmt"
	"time"
	"todo-ai/common/holiday"
	"todo-ai/core/logger"

	"github.com/gin-gonic/gin"
)

type CalendarDataRequest struct {
	Year int `json:"year"` // 普通节日的年份
}

// 获取日历数据
func CalendarData(ctx *gin.Context) (interface{}, error) {
	var request ListTaskRequest
	if err := ctx.Bind(&request); err != nil {
		logger.Errorf("ListTask BindJSON error:%s", err)
		return nil, err
	}

	if request.Year == 0 {
		return nil, errors.New("year is empty")
	}

	return GetHolidays(request.Year), nil
}

// Holiday 节日结构体
type Holiday struct {
	Date      string `json:"date"`
	Name      string `json:"name"`
	IsHoliday bool   `json:"isHoliday"`
	Type      string `json:"type"` // legal: 法定节假日, normal: 普通节日
}

// GetHolidays 获取指定年份的节假日信息
func GetHolidays(year int) []Holiday {
	holidays := make([]Holiday, 0)

	// 添加法定节假日
	for dateStr, name := range holiday.LegalHolidays {
		holidays = append(holidays, Holiday{
			Date:      dateStr,
			Name:      name,
			IsHoliday: true,
			Type:      "legal",
		})
	}

	// 添加普通节日
	for monthDay, name := range holiday.NormalHolidays {
		dateStr := fmt.Sprintf("%d-%s", year, monthDay)
		holidays = append(holidays, Holiday{
			Date:      dateStr,
			Name:      name,
			IsHoliday: false,
			Type:      "normal",
		})
	}

	return holidays
}

// 判断是否是节假日
func IsHoliday(date time.Time) (bool, string) {
	dateStr := date.Format("2006-01-02")
	if name, ok := holiday.LegalHolidays[dateStr]; ok {
		return true, name
	}

	monthDay := date.Format("01-02")
	if name, ok := holiday.NormalHolidays[monthDay]; ok {
		return false, name
	}

	return false, ""
}
