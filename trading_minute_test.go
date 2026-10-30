package thsdk

import (
	"testing"
	"time"
)

func TestTradingMinuteCoversTheFullSession(t *testing.T) {
	session := time.Date(2026, 10, 30, 9, 30, 0, 0, Shanghai)

	cases := []struct {
		index int
		want  string
	}{
		{0, "09:30"},
		{1, "09:31"},
		{119, "11:29"}, // last morning minute
		{120, "13:00"}, // first afternoon minute, lunch break skipped
		{121, "13:01"},
		{239, "14:59"},
	}
	for _, item := range cases {
		if got := tradingMinute(session, item.index).Format("15:04"); got != item.want {
			t.Errorf("index %d: got %s want %s", item.index, got, item.want)
		}
	}
}

func TestTradingMinuteNeverFallsInTheLunchBreak(t *testing.T) {
	session := time.Date(2026, 10, 30, 9, 30, 0, 0, Shanghai)
	for index := 0; index < 240; index++ {
		clock := tradingMinute(session, index)
		hour, minute, _ := clock.Clock()
		if hour == 11 && minute >= 30 {
			t.Fatalf("index %d landed at %02d:%02d", index, hour, minute)
		}
		if hour == 12 {
			t.Fatalf("index %d landed in the 12 o'clock hour", index)
		}
		if hour > 14 || (hour == 14 && minute > 59) {
			if clock.After(time.Date(2026, 10, 30, 15, 0, 0, 0, Shanghai)) {
				t.Fatalf("index %d passed the close: %v", index, clock)
			}
		}
	}
}

func TestTradingMinuteIsStrictlyIncreasing(t *testing.T) {
	session := time.Date(2026, 10, 30, 9, 30, 0, 0, Shanghai)
	previous := tradingMinute(session, 0)
	for index := 1; index < 240; index++ {
		current := tradingMinute(session, index)
		if !current.After(previous) {
			t.Fatalf("index %d did not advance: %v after %v", index, current, previous)
		}
		previous = current
	}
}
