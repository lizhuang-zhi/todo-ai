package utils

import (
	"testing"
	"time"
)

func TestGetDefaultTimeZone(t *testing.T) {
	t.Log("default time zone: ", time.Local)
}

func TestSetDefaultTimeZone(t *testing.T) {
	timezone := "UTC"
	err := SetDefaultTimeZone(timezone)
	if err != nil {
		panic(err)
	}

	t.Log("time zone: ", time.Local)
}

func TestNowDateTime(t *testing.T) {
	t.Log("current date time: ", NowDateTime())
}

func TestFormatTime(t *testing.T) {
	now := time.Now()
	formatted := FormatTime(now)
	if formatted == "" {
		t.Fatal("formatted time should not be empty")
	}

	t.Log("formatted time: ", formatted)
}

func TestTimestampToString(t *testing.T) {
	// 设置时区为UTC
	err := SetDefaultTimeZone("UTC")
	if err != nil {
		panic(err)
	}

	timestamp := time.Now().Unix()
	t.Log("timestamp to string: ", TimestampToString(timestamp))
}

func TestIsSameDay(t *testing.T) {
	t1 := time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 7, 4, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 8, 4, 0, 0, 0, 0, time.UTC)
	if !IsSameDay(t1, t2) {
		t.Fatal("t1 and t2 should be same day")
	}
	if IsSameDay(t1, t3) {
		t.Fatal("t1 and t3 should not be same day")
	}
}

func TestGetFirstDateOfWeekWithParam(t *testing.T) {
	today := time.Now()
	monday := GetFirstDateOfWeekWithParam(today)
	t.Log("获取本周周一: ", monday)
}

func TestGetLastDateOfWeekWithParam(t *testing.T) {
	// today := time.Now()

	// 构造指定日期的time.Time对象
	today := time.Date(2024, 9, 10, 0, 0, 0, 0, time.Local)

	sunday := GetLastDateOfWeekWithParam(today)
	t.Log("获取本周周日: ", sunday)
}

func TestGetLastDateOfWeekWithParamByFormat(t *testing.T) {
	today := time.Now()
	format := "01月02日"

	sunday := GetLastDateOfWeekWithParamByFormat(today, format)
	t.Log("获取本周周日: ", sunday)
}

func TestGetLastWeekFirstDate(t *testing.T) {
	t.Log("获取上周周一: ", GetLastWeekFirstDate())
}

func TestGetLastWeekLastDate(t *testing.T) {
	t.Log("获取上周周日: ", GetLastWeekLastDate())
}

func TestGetLastWeekLastDateWithParamByFormat(t *testing.T) {
	// 构造指定日期的time.Time对象
	day := time.Date(2024, 9, 9, 0, 0, 0, 0, time.Local)
	format := "01月02日"

	t.Log("获取上周周日: ", GetLastWeekLastDateWithParamByFormat(day, format))
}

func TestGetLastFridayWeekLastDateWithParam(t *testing.T) {
	// 构造指定日期的time.Time对象
	day := time.Date(2024, 9, 15, 0, 0, 0, 0, time.Local)

	t.Log("获取上周周五: ", GetLastWeekFiveDateWithParam(day))
}

func TestTimeDiff(t *testing.T) {
	start := 1727853642
	end := 1728717828

	t.Log("时间差: ", TimeDiff(int64(start), int64(end)))
}

func TestIsTimeDiffInDayCnt(t *testing.T) {
	start := 1728631242
	end := 1728717828

	diff := TimeDiff(int64(start), int64(end))
	t.Log("时间差是否在5天内 ", IsTimeDiffInDayCnt(diff, 5))
}

func TestNowTime(t *testing.T) {
	nowStamp := time.Now().Unix()

	t.Log("当前时间戳（秒级） ", nowStamp)
}

func TestGetNextDate(t *testing.T) {
	today := "2024/10/24"
	nextDate, err := GetNextDay(today)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("获取明天日期: ", nextDate)
}
