package stats

type MostUsedDomainResponse struct {
	Domain string
	Count  int
}

type ServerUsageByTimeRangeResponse struct {
	TimeRange string
	Count     int
}
