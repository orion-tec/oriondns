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
