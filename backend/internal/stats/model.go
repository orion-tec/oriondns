package stats

import "time"

type MostUsedDomainResponse struct {
	Domain string
	Count  int
}

type ServerUsageByTimeRangeResponse struct {
	TimeRange time.Time
	Count     int
}
