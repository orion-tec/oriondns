package stats

type MostUsedDomainResponse struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}
