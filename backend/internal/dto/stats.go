package dto

type GetMostUsedDomainsRequest struct {
	From       int64    `json:"from"`
	To         int64    `json:"to"`
	Categories []string `json:"categories"`
}

type GetMostUsedDomainsResponse struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type GetServerUsageByTimeRangeRequest struct {
	From       int64    `json:"from"`
	To         int64    `json:"to"`
	Categories []string `json:"categories"`
}

type GetServerUsageByTimeRangeResponse struct {
	TimeRange string `json:"timeRange"`
	Count     int    `json:"count"`
}
