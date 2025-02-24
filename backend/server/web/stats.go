package web

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type MostUsedDomainsRequest struct {
	Range      string   `json:"range"`
	Categories []string `json:"categories"`
}

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
	enableCors(&w)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req MostUsedDomainsRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	to := getTo(req.Range)
	from := getFrom(req.Range)

	results, err := h.stats.GetMostUsedDomains(context.Background(), from, to, req.Categories, 10)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(results)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}

}
