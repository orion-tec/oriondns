package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var m sync.Map
var cache sync.Map

func handleRequest(c *dns.Client, blockedNames []string) dns.HandlerFunc {
	return func(rw dns.ResponseWriter, msg *dns.Msg) {

		for _, q := range msg.Question {
			actual, loaded := m.LoadOrStore(q.Name, 1)
			if loaded {
				count := actual.(int)
				m.Store(q.Name, count+1)
			}
		}

		respFromCache, loaded := cache.Load(msg.String())
		if loaded {
			fmt.Println("Getting from cache")
			rw.WriteMsg(respFromCache.(*dns.Msg))
			return
		} else {
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

				fmt.Printf("Question: %s\n", q.Name)
			}

		}

		resp, _, err := c.Exchange(msg, "8.8.8.8:53")
		if err != nil {
			log.Printf("Failed to exchange: %s", err.Error())
			return
		}

		cache.Store(msg.String(), resp)
		rw.WriteMsg(resp)
	}
}

func printMap() {
	for {
		m.Range(func(key, value interface{}) bool {
			fmt.Printf("%s: %d\n", key, value)
			return true
		})

		fmt.Println("=====================================")

		time.Sleep(5 * time.Second)
	}
}

func main() {
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

	srv := &dns.Server{Addr: ":53", Net: "udp"}
	fmt.Println("Listening on :53")
	srv.Handler = dns.HandlerFunc(handleRequest(c, blockedNames))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener: %s", err.Error())
	}

}
