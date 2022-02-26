package datetime

import "time"

func FormatLongWeekDay(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 Monday")
}
