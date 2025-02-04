package dns

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/orion-tec/oriondns/internal/stats"
	"go.uber.org/fx"
)

type DNS struct {
	m     sync.Map
	cache sync.Map
	stats stats.DB
}

func (d *DNS) handleRequest(c *dns.Client, blockedNames []string) dns.HandlerFunc {
	return func(rw dns.ResponseWriter, msg *dns.Msg) {

		for _, q := range msg.Question {
			err := d.stats.Insert(context.Background(), time.Now(), q.Name)
			if err != nil {
				log.Printf("Failed to insert stats: %s", err.Error())
			}
		}

		respFromCache, loaded := d.cache.Load(msg.String())
		if loaded {
			rw.WriteMsg(respFromCache.(*dns.Msg))
			return
		}

		for _, q := range msg.Question {
			for _, v := range blockedNames {
				if strings.Contains(q.Name, v) {
					m := new(dns.Msg)
					m.Answer = append(m.Answer, &dns.A{
						A: net.ParseIP("127.0.0.1"),
					})
					m.RecursionAvailable = true
					m.SetReply(msg)

					rw.WriteMsg(m)
					return
				}
			}
		}

		resp, _, err := c.Exchange(msg, "8.8.8.8:53")
		if err != nil {
			log.Printf("Failed to exchange: %s", err.Error())
			return
		}

		d.cache.Store(msg.String(), resp)
		rw.WriteMsg(resp)
	}
}

func New(lc fx.Lifecycle, stats stats.DB) *DNS {
	c := new(dns.Client)

	blockedNames := []string{
		"googleads",
		"cleverwebserver.com",
		"googleadservices",
		"doubleclick",
		"googlesyndication",
		"googletagservices",
		"googletagmanager",
		"googletagmanager",
	}

	dnsStruct := DNS{
		stats: stats,
	}

	srv := &dns.Server{Addr: ":53", Net: "udp"}

	srv.Handler = dns.HandlerFunc(dnsStruct.handleRequest(c, blockedNames))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Println("Listening on :53")
				if err := srv.ListenAndServe(); err != nil {
					log.Fatalf("Failed to set udp listener: %s", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return &dnsStruct
}
