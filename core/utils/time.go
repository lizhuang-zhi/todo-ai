package utils

import "time"

// 设置默认时区
func SetDefaultTimeZone(timezone string) error {
	if timezone == "" || timezone == "Local" {
		// 采用本地时区
		return nil
	}

	if timezone == "UTC" {
		// 采用 UTC 时区
		time.Local = time.UTC
		return nil
	}

	/*
		除开 Local 和 UTC 时区外，其他时区还有如下几种：
		- "America/New_York"
		- "Asia/Shanghai"
		- "Europe/London"
		- "Europe/Berlin"
		- "Pacific/Auckland"
		- "Australia/Sydney"
		- ...
	*/

	// 读取时区信息
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	time.Local = loc

	return nil
}

// 当前毫秒级时间戳
func NowUnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// 当前日期时间
func NowDateTime() string {
	return FormatTime(time.Now())
}

// 格式化时间
func FormatTime(t time.Time) string {
	return t.Format(time.RFC822Z)
}

// 时间戳转化为字符串形式显示
func TimestampToString(t int64) string {
	return time.Unix(t, 0).Format(time.RFC3339)
}

// 检查小时时间是否正常
func IsValidHour(hour int) bool {
	return hour >= 0 && hour <= 23
}

// 检查分钟时间是否正常
func IsValidMinute(minute int) bool {
	return minute >= 0 && minute <= 59
}

// 检查秒时间是否正常
func IsValidSecond(second int) bool {
	return second >= 0 && second <= 59
}

// 是否为同一天
func IsSameDay(left time.Time, right time.Time) bool {
	return (left.Year() == right.Year()) && (left.YearDay() == right.YearDay())
}

// 是否为同一周
func IsSameWeek(left time.Time, right time.Time) bool {
	yearLeft, Left := left.ISOWeek()
	yearRight, weekRight := right.ISOWeek()
	return (yearLeft == yearRight) && (Left == weekRight)
}

// 是否为同一月
func IsSameMonth(left time.Time, right time.Time) bool {
	return (left.Year() == right.Year()) && (left.Month() == right.Month())
}

// 判读某个时间是否在开始时间和结束时间内（包含开始时间和结束时间）
func InTime(t time.Time, start time.Time, end time.Time) bool {
	// 如果与开始时间或者结束时间相等，则返回True
	if t.Equal(start) || t.Equal(end) {
		return true
	}

	// 如果开始时间不为空，且时间比开始时间早，则返回False
	if start.Unix() != 0 && t.Before(start) {
		return false
	}

	// 如果结束时间不为空，且时间比开始时间晚，则返回False
	if end.Unix() != 0 && t.After(end) {
		return false
	}

	return true
}

// 获取指定时间周一的日期
func GetFirstDateOfWeekWithParam(t time.Time) (weekMonday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekMonday = weekStartDate.Format("2006-01-02")
	return
}

// 获取指定时间周一的日期(自定义格式)
func GetFirstDateOfWeekWithParamByFormat(t time.Time, format string) (weekMonday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekMonday = weekStartDate.Format(format)
	return
}

// 获取指定时间周日的日期
func GetLastDateOfWeekWithParam(t time.Time) (weekSunday string) {
	offset := int(time.Sunday - t.Weekday())
	if offset < 0 {
		offset += 7
	}

	weekEndDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekSunday = weekEndDate.Format("2006-01-02")
	return
}

// 获取指定时间周日的日期(自定义格式)
func GetLastDateOfWeekWithParamByFormat(t time.Time, format string) (weekSunday string) {
	offset := int(time.Sunday - t.Weekday())
	if offset < 0 {
		offset += 7
	}

	weekEndDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekSunday = weekEndDate.Format(format)
	return
}

// 获取上周周一的日期
func GetLastWeekFirstDate() (weekMonday string) {
	return GetLastWeekFirstDateWithParam(time.Now())
}

// 获取指定时间上周一的日期
func GetLastWeekFirstDateWithParam(t time.Time) (weekMonday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset-7)
	weekMonday = weekStartDate.Format("2006-01-02")
	return
}

// 获取指定时间上周一的日期(自定义格式)
func GetLastWeekFirstDateWithParamByFormat(t time.Time, format string) (weekMonday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset-7)
	weekMonday = weekStartDate.Format(format)
	return
}

// 获取指定时间上周五的日期
func GetLastWeekFiveDateWithParam(t time.Time) (lastFriday string) {
	offset := int(time.Friday - t.Weekday())
	if t.Weekday() == 0 { // t为周日
		offset -= 7
	}

	offset -= 7

	lastFridayDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastFriday = lastFridayDate.Format("2006-01-02")
	return
}

// 获取指定时间上周五的日期(自定义格式)
func GetLastWeekFiveDateWithParamByFormat(t time.Time, format string) (lastFriday string) {
	offset := int(time.Friday - t.Weekday())
	if t.Weekday() == 0 { // t为周日
		offset -= 7
	}

	offset -= 7

	lastFridayDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastFriday = lastFridayDate.Format(format)
	return
}

// 获取上周周日的日期
func GetLastWeekLastDate() (weekSunday string) {
	return GetLastWeekLastDateWithParam(time.Now())
}

// 获取指定时间上周日的日期
func GetLastWeekLastDateWithParam(t time.Time) (weekSunday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekEndDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset-1)
	weekSunday = weekEndDate.Format("2006-01-02")
	return
}

// 获取指定时间上周日的日期(自定义格式)
func GetLastWeekLastDateWithParamByFormat(t time.Time, format string) (weekSunday string) {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekEndDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset-1)
	weekSunday = weekEndDate.Format(format)
	return
}

// 计算两个时间戳（秒级）之间的时间差
func TimeDiff(start int64, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return end - start
}

// 判断时间差（秒级）是否在x天内
func IsTimeDiffInDayCnt(diff int64, dayCnt int) bool {
	return diff <= int64(dayCnt*24*60*60)
}

// 获取指定日期（格式2006/01/01）下一天的日期，格式为"2006/01/02"
func GetNextDay(date string) (string, error) {
	t, err := time.Parse("2006/01/02", date)
	if err != nil {
		return "", err
	}
	t = t.AddDate(0, 0, 1)
	return t.Format("2006/01/02"), nil
}

// 秒级时间戳转时间格式(2006/01/02 15:04:05)
func UnixToDateTime(unix int64) string {
	return time.Unix(unix, 0).Format("2006/01/02 15:04:05")
}
