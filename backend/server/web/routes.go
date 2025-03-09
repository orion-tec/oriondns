package web

import "net/http"

func withCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handler.ServeHTTP(w, r)
	}
}

func (h *HTTP) setupRoutes() {
	http.HandleFunc("POST /api/v1/dashboard/most-used-domains", withCors(h.getMostUsedDomainsDashboard))
	http.HandleFunc("POST /api/v1/dashboard/server-usage-by-time-range", withCors(h.getServerUsageByTimeRangeDashboard))
}
