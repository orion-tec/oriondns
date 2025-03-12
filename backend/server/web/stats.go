package web

import (
	"context"
	"net/http"
	"time"

	"github.com/orion-tec/oriondns/internal/dto"
)

func getTimeFromFE(t int64) time.Time {
	return time.Unix(t/1000, 0).UTC()
}

func (h *HTTP) getMostUsedDomainsDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.GetMostUsedDomainsRequest
	err := readFromJSON(r, &req)
	if err != nil {
		logAndWriteError(w, err)
		return
	}

	from := getTimeFromFE(req.From)
	to := getTimeFromFE(req.To)

	results, err := h.stats.GetMostUsedDomains(context.Background(), from, to, req.Categories, 10)
	if err != nil {
		logAndWriteError(w, err)
		return
	}

	transformedResult := make([]dto.GetMostUsedDomainsResponse, len(results))
	for i, r := range results {
		transformedResult[i] = dto.GetMostUsedDomainsResponse{
			Domain: r.Domain,
			Count:  r.Count,
		}
	}

	responseWithJSON(w, transformedResult)
}

func (h *HTTP) getServerUsageByTimeRangeDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.GetServerUsageByTimeRangeRequest
	err := readFromJSON(r, &req)
	if err != nil {
		logAndWriteError(w, err)
		return
	}

	from := getTimeFromFE(req.From)
	to := getTimeFromFE(req.To)

	results, err := h.stats.GetServerUsageByTimeRange(context.Background(), from, to, req.Categories)
	if err != nil {
		logAndWriteError(w, err)
		return
	}

	transformedResult := make([]dto.GetServerUsageByTimeRangeResponse, len(results))
	for i, r := range results {
		transformedResult[i] = dto.GetServerUsageByTimeRangeResponse{
			TimeRange: r.TimeRange,
			Count:     r.Count,
		}
	}

	responseWithJSON(w, transformedResult)
}
