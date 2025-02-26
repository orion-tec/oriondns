package web

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"

	"github.com/orion-tec/oriondns/internal/stats"
)

type HTTP struct {
	stats stats.DB

	s *http.Server
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("HTTP %s %s", r.Method, r.URL.Path)
}

func New(lc fx.Lifecycle, stats stats.DB) *HTTP {

	httpStruct := HTTP{
		stats: stats,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				http.HandleFunc("POST /api/v1/dashboard/most-used-domains",
					http.HandlerFunc(httpStruct.getMostUsedDomainsDashboard))
				err := http.ListenAndServe(":8080", nil)
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			return nil
		},
	})
	return &httpStruct
}
