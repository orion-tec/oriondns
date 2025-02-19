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
	"github.com/orion-tec/oriondns/internal/blockeddomains"
	"github.com/orion-tec/oriondns/internal/domains"
	"github.com/orion-tec/oriondns/internal/stats"
	"go.uber.org/fx"
)

type DNS struct {
	cacheMap             sync.Map
	blockedDomainsMap    map[string][]blockeddomains.BlockedDomain
	blockedDomainsMutext sync.Mutex

	blockedDomains blockeddomains.DB
	domain         domains.DB
	stats          stats.DB
}

func (d *DNS) updateBlockedDomainsMap(blockedDomains []blockeddomains.BlockedDomain) {
	d.blockedDomainsMutext.Lock()
	defer d.blockedDomainsMutext.Unlock()

	for k := range d.blockedDomainsMap {
		delete(d.blockedDomainsMap, k)
	}

	for _, bd := range blockedDomains {
		id := fmt.Sprintf("%d", bd.ID)
		d.blockedDomainsMap[id] = append(d.blockedDomainsMap[id], bd)
	}
}

func (d *DNS) updateBlockedDomains() {
	for {
		fmt.Println("Updating blocked domains")
		bds, err := d.blockedDomains.GetAll(context.Background())
		if err != nil {
			fmt.Println(err)
		}

		d.updateBlockedDomainsMap(bds)
		time.Sleep(1 * time.Minute)
	}
}

func (d *DNS) handleRequest(c *dns.Client) dns.HandlerFunc {
	return func(rw dns.ResponseWriter, msg *dns.Msg) {

		// Store stats for the request
		go func() {
			for _, q := range msg.Question {
				name := strings.TrimSuffix(q.Name, ".")
				err := d.domain.Insert(context.Background(), name)
				if err != nil {
					log.Printf("Failed to insert domain: %s", err.Error())
				}

				err = d.stats.Insert(context.Background(), time.Now(), q.Name)
				if err != nil {
					log.Printf("Failed to insert stats: %s", err.Error())
				}
			}
		}()

		// Validate if it's blocked
		isBlocked := false
		for _, bds := range d.blockedDomainsMap {
			for _, q := range msg.Question {
				for _, bd := range bds {
					if bd.Recursive && strings.HasSuffix(q.Name, bd.Domain) {
						fmt.Println("Blocked recursive domain: ", q.Name)
						isBlocked = true
						break
					}

					if q.Name == bd.Domain {
						fmt.Println("Blocked domain: ", q.Name)
						isBlocked = true
						break
					}
				}
			}
		}

		if isBlocked {
			m := new(dns.Msg)
			m.Answer = append(m.Answer, &dns.A{
				A: net.ParseIP("127.0.0.1"),
			})
			m.RecursionAvailable = true
			m.SetReply(msg)
			rw.WriteMsg(m)
			return
		}

		respFromCache, loaded := d.cacheMap.Load(msg.String())
		if loaded {
			rw.WriteMsg(respFromCache.(*dns.Msg))
			return
		}

		resp, _, err := c.Exchange(msg, "8.8.8.8:53")
		if err != nil {
			log.Printf("Failed to exchange: %s", err.Error())
			return
		}

		d.cacheMap.Store(msg.String(), resp)
		rw.WriteMsg(resp)
	}
}

func New(lc fx.Lifecycle, stats stats.DB, blockedDomains blockeddomains.DB, domain domains.DB) *DNS {
	c := new(dns.Client)

	dnsStruct := DNS{
		stats:             stats,
		blockedDomains:    blockedDomains,
		domain:            domain,
		blockedDomainsMap: make(map[string][]blockeddomains.BlockedDomain),
	}

	go dnsStruct.updateBlockedDomains()

	srv := &dns.Server{Addr: ":53", Net: "udp"}

	srv.Handler = dns.HandlerFunc(dnsStruct.handleRequest(c))

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
