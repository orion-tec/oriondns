package web

import (
	"context"
	"net/http"
	"time"

	"github.com/orion-tec/oriondns/internal/dto"
)

func getTo(rng string) time.Time {
	now := time.Now()
	switch rng {
	case "Yesterday":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	default:
		return now
	}
}

// 'Last month', 'Last 2 weeks', 'Last week', 'Last 3 days', 'Yesterday', 'Today'
func getFrom(rng string) time.Time {
	now := time.Now()
	lastMidNight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch rng {
	case "Last month":
		return lastMidNight.AddDate(0, -1, 0)
	case "Last 2 weeks":
		return lastMidNight.AddDate(0, 0, -14)
	case "Last week":
		return lastMidNight.AddDate(0, 0, -7)
	case "Last 3 days":
		return lastMidNight.AddDate(0, 0, -3)
	case "Yesterday":
		return lastMidNight.AddDate(0, 0, -1)
	default:
		return lastMidNight
	}
}

func (h *HTTP) getMostUsedDomainsDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.GetMostUsedDomainsRequest
	err := readFromJSON(r, &req)
	if err != nil {
		logAndWriteError(w, err)
		return
	}

	to := getTo(req.Range)
	from := getFrom(req.Range)

	to = to.UTC()
	from = from.UTC()

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

	to := getTo(req.Range)
	from := getFrom(req.Range)

	to = to.UTC()
	from = from.UTC()

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
