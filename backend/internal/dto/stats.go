package dto

type GetMostUsedDomainsRequest struct {
	Range      string   `json:"range"`
	Categories []string `json:"categories"`
}

type GetMostUsedDomainsResponse struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}
