package web

import (
	"context"
	"log"
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
			go func() {
				httpStruct.setupRoutes()

				// Create HTTP server with proper timeouts to satisfy G114
				server := &http.Server{
					Addr:         ":8080",
					Handler:      nil,
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
					IdleTimeout:  60 * time.Second,
				}

				err := server.ListenAndServe()
				if err != nil {
					log.Fatalf("HTTP server failed to start: %v", err)
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
