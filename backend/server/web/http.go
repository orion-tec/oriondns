package web

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	"github.com/orion-tec/oriondns/internal/stats"
)

type HTTP struct {
	stats stats.DB

	s *http.Server
}

type HttpDeps struct {
	fx.In
	Stats stats.DB
}

func New(lc fx.Lifecycle, deps HttpDeps) *HTTP {
	httpStruct := HTTP{
		stats: deps.Stats,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			httpStruct.setupRoutes()

			httpStruct.s = &http.Server{
				Addr:         ":8080",
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
			}

			go func() {
				err := httpStruct.s.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return httpStruct.s.Shutdown(ctx)
		},
	})
	return &httpStruct
}
