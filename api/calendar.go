package api

import (
	"errors"
	"fmt"
	"time"
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

// 法定节假日：持续补充
var legalHolidays = map[string]string{
	"2024-10-01": "国庆节",
	"2024-10-02": "国庆节",
	"2024-10-03": "国庆节",
	"2024-10-04": "国庆节",
	"2024-10-05": "国庆节",
	"2024-10-06": "国庆节",
	"2024-10-07": "国庆节",

	// 2025年上半年法定节假日
	"2025-01-01": "元旦",
	"2025-01-28": "春节",
	"2025-01-29": "春节",
	"2025-01-30": "春节",
	"2025-01-31": "春节",
	"2025-02-01": "春节",
	"2025-02-02": "春节",
	"2025-02-03": "春节",
	"2025-02-04": "春节",
	"2025-04-04": "清明节",
	"2025-04-05": "清明节",
	"2025-04-06": "清明节",
	"2025-05-01": "劳动节",
	"2025-05-02": "劳动节",
	"2025-05-03": "劳动节",
	"2025-06-29": "端午节",
	"2025-06-30": "端午节",
}

// 普通节日
var normalHolidays = map[string]string{
	"02-14": "情人节",
	"03-08": "妇女节",
	"03-12": "植树节",
	"04-01": "愚人节",
	"05-04": "青年节",
	"06-01": "儿童节",
	"07-01": "建党节",
	"08-01": "建军节",
	"09-10": "教师节",
	"12-24": "平安夜",
	"12-25": "圣诞节",
}

// GetHolidays 获取指定年份的节假日信息
func GetHolidays(year int) []Holiday {
	holidays := make([]Holiday, 0)

	// 添加法定节假日
	for dateStr, name := range legalHolidays {
		holidays = append(holidays, Holiday{
			Date:      dateStr,
			Name:      name,
			IsHoliday: true,
			Type:      "legal",
		})
	}

	// 添加普通节日
	for monthDay, name := range normalHolidays {
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
	if name, ok := legalHolidays[dateStr]; ok {
		return true, name
	}

	monthDay := date.Format("01-02")
	if name, ok := normalHolidays[monthDay]; ok {
		return false, name
	}

	return false, ""
}
