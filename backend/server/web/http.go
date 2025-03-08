package web

import (
	"context"
	"net/http"

	"go.uber.org/fx"

	"github.com/orion-tec/oriondns/internal/stats"
)

type HTTP struct {
	stats stats.DB

	s *http.Server
}

func New(lc fx.Lifecycle, stats stats.DB) *HTTP {
	httpStruct := HTTP{
		stats: stats,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				httpStruct.setupRoutes()

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
