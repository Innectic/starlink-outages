package uptime

import (
	"time"
)

type LastData struct {
	PopPingDropRate int
	OutageStart time.Time
	OutageEnd time.Time
	Cause DishyStatus

	Collecting bool
}

func (l LastData) FriendlyStartTime() string {
	return l.OutageStart.Format("01/02/2006 15:04:05 MST")
}

func (l LastData) Duration() int64 {
	return l.OutageEnd.Unix() - l.OutageStart.Unix()
}

func (l *LastData) Started() {
	l.OutageStart = time.Now()
}

func (l *LastData) Ended() {
	l.OutageEnd = time.Now()
}
