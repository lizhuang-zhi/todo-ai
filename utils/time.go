package utils

import (
	"fmt"
	"strings"
	"time"
	"todo-ai/core/utils"
)

// 推算上周一到周日的日期范围（格式：01月01日 - 01月07日）
func GetLastWeekDateRange() string {
	t := time.Now() // 获取当前时间
	startStr := utils.GetLastWeekFirstDateWithParamByFormat(t, "01月02日")
	endStr := utils.GetLastWeekLastDateWithParamByFormat(t, "01月02日")
	return fmt.Sprintf("%s - %s", startStr, endStr)
}

// 推算上周一到周五的日期范围（格式：01月01日 - 01月05日）
func GetLastWeekDateRangeFromMondayToFriday() string {
	t := time.Now() // 获取当前时间
	startStr := utils.GetLastWeekFirstDateWithParamByFormat(t, "01月02日")
	endStr := utils.GetLastWeekFiveDateWithParamByFormat(t, "01月02日")
	return fmt.Sprintf("%s - %s", startStr, endStr)
}

// 推算当前时间的本周一到周日的日期范围（格式：8月12日 ~ 8月18日）
func GetCurWeekDateRangeFormat() string {
	t := time.Now() // 获取当前时间
	startStr := utils.GetFirstDateOfWeekWithParamByFormat(t, "1月2日")
	endStr := utils.GetLastDateOfWeekWithParamByFormat(t, "1月2日")

	return fmt.Sprintf("%s ~ %s", startStr, endStr)
}

// 获取当天的日期（格式：2024-08-26）
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// 获取当天的日期（格式：/）
func GetTodayDateFormat(split string) string {
	return time.Now().Format("2006" + split + "01" + split + "02")
}

// 提取日期范围的开始和结束日期
func ExtractDateRange(name string) []string {
	parts := strings.Split(name, " - ")
	if len(parts) == 2 {
		return parts
	}
	return nil
}

// 检查给定的日期是否在指定的范围内（包括临界日期，截止时间点的天数+2）
func IsDateWithinRange(date, startDate, endDate string) bool {
	// 1月2日格式 允许日期不含前导0
	currentTime, _ := time.Parse("1月2日", date)
	startTime, _ := time.Parse("1月2日", startDate)
	endTime, _ := time.Parse("1月2日", endDate)
	return (currentTime.Equal(startTime) || currentTime.After(startTime)) && (currentTime.Equal(endTime) || currentTime.Before(endTime))
}
