package main

import (
	"github.com/orion-tec/oriondns/config"
	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/stats"
	"github.com/orion-tec/oriondns/server/dns"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.New),
		fx.Provide(db.New),
		fx.Provide(stats.New),
		fx.Provide(dns.New),
		fx.Invoke(func(s *dns.DNS) {}),
	).Run()
}
