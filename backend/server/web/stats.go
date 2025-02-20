package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h *HTTP) getMostUsedDomainsDashboard(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	results, err := h.stats.GetMostUsedDomains(context.Background(), oneWeekAgo, now, 10)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(results)
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
