package main

import (
	"go.uber.org/fx"

	"github.com/orion-tec/oriondns/config"
	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/stats"
	"github.com/orion-tec/oriondns/server/web"
)

func main() {
	fx.New(
		fx.Provide(config.New),
		fx.Provide(db.New),
		fx.Provide(stats.New),
		fx.Provide(web.New),
		fx.Invoke(func(s *web.HTTP) {}),
	).Run()
}
