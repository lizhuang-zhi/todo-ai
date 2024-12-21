package utils

import "testing"

func TestGetLastWeekDateRange(t *testing.T) {
	t.Log(GetLastWeekDateRange())
}

func TestGetLastWeekDateRangeFromMondayToFriday(t *testing.T) {
	t.Log(GetLastWeekDateRangeFromMondayToFriday())
}

func TestGetCurWeekDateRangeFormat(t *testing.T) {
	t.Log(GetCurWeekDateRangeFormat())
}
