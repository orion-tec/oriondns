package dto

type GetMostUsedDomainsRequest struct {
	Range      string   `json:"range"`
	Categories []string `json:"categories"`
	TzOffset   int      `json:"tzOffset"`
}

type GetMostUsedDomainsResponse struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type GetServerUsageByTimeRangeRequest struct {
	Range      string   `json:"range"`
	Categories []string `json:"categories"`
	TzOffset   int      `json:"tzOffset"`
}

type GetServerUsageByTimeRangeResponse struct {
	TimeRange string `json:"timeRange"`
	Count     int    `json:"count"`
}
