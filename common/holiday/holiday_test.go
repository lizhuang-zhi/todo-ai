package holiday

import (
	"testing"
	"time"
)

func TestIsWeekend(t *testing.T) {
	date1 := "2024-12-28" // 周六
	date2 := "2024-12-29" // 周天
	date3 := "2024-12-30" // 周一
	date4 := "2024-12-27" // 周五

	d1, err := time.Parse("2006-01-02", date1)
	if err != nil {
		t.Errorf("time.Parse error:%s", err)
		return
	}

	d2, err := time.Parse("2006-01-02", date2)
	if err != nil {
		t.Errorf("time.Parse error:%s", err)
		return
	}

	d3, err := time.Parse("2006-01-02", date3)
	if err != nil {
		t.Errorf("time.Parse error:%s", err)
		return
	}

	d4, err := time.Parse("2006-01-02", date4)
	if err != nil {
		t.Errorf("time.Parse error:%s", err)
		return
	}

	t.Log("d1:", d1, "is weekend:", IsWeekend(d1))
	t.Log("d2:", d2, "is weekend:", IsWeekend(d2))
	t.Log("d3:", d3, "is weekend:", IsWeekend(d3))
	t.Log("d4:", d4, "is weekend:", IsWeekend(d4))
}
