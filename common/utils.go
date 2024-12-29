package common

import "time"

// GetYearByDate 获取日期(2024-12-23)的年份的年份
func GetYearByDate(date string) int {
	// 解析日期
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	return t.Year()
}
