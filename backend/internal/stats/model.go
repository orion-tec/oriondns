package stats

import "time"

type MostUsedDomainResponse struct {
	Domain string
	Count  int64
}

type ServerUsageByTimeRangeResponse struct {
	TimeRange time.Time
	Count     int64
}
