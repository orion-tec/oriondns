package dto

import "time"

type GetMostUsedDomainsRequest struct {
	From       int64    `json:"from"`
	To         int64    `json:"to"`
	Categories []string `json:"categories"`
}

type GetMostUsedDomainsResponse struct {
	Domain string `json:"domain"`
	Count  int64  `json:"count"`
}

type GetServerUsageByTimeRangeRequest struct {
	From       int64    `json:"from"`
	To         int64    `json:"to"`
	Categories []string `json:"categories"`
}

type GetServerUsageByTimeRangeResponse struct {
	TimeRange time.Time `json:"timeRange"`
	Count     int64     `json:"count"`
}
