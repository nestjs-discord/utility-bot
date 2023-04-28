package uptime

import (
	"github.com/dustin/go-humanize"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func Uptime() string {
	return humanize.Time(startTime)
}

//func UptimeDuration() time.Duration {
//	return time.Since(startTime)
//}
