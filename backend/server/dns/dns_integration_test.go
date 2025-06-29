package dns

import (
	"context"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/internal/blockeddomains"
	"github.com/orion-tec/oriondns/internal/testutil"
)

func TestDNS_Integration_BlockedDomainHandling(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	
	dnsHandler := createTestDNS()
	
	testDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware.com.", Recursive: false},
		{ID: 2, Domain: ".ads.example.com", Recursive: true},
	}
	dnsHandler.updateBlockedDomainsMap(testDomains)

	testCases := []struct {
		name        string
		queryDomain string
		queryType   uint16
		shouldBlock bool
		description string
	}{
		{
			name:        "Block exact match",
			queryDomain: "malware.com.",
			queryType:   dns.TypeA,
			shouldBlock: true,
			description: "Should block exact domain match",
		},
		{
			name:        "Block recursive match",
			queryDomain: "tracker.ads.example.com",
			queryType:   dns.TypeA,
			shouldBlock: true,
			description: "Should block subdomain with recursive rule",
		},
		{
			name:        "Allow non-blocked domain",
			queryDomain: "google.com.",
			queryType:   dns.TypeA,
			shouldBlock: false,
			description: "Should allow non-blocked domains",
		},
		{
			name:        "Allow non-matching recursive",
			queryDomain: "safe.example.com",
			queryType:   dns.TypeA,
			shouldBlock: false,
			description: "Should allow domains that don't match recursive rules",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := &dns.Msg{}
			msg.Question = []dns.Question{
				{Name: tc.queryDomain, Qtype: tc.queryType},
			}

			isBlocked := false
			for _, bds := range dnsHandler.blockedDomainsMap {
				for _, q := range msg.Question {
					for _, bd := range bds {
						if bd.Recursive && dns.IsSubDomain(bd.Domain, q.Name) {
							isBlocked = true
							break
						}
						if q.Name == bd.Domain {
							isBlocked = true
							break
						}
					}
				}
			}

			assert.Equal(t, tc.shouldBlock, isBlocked, tc.description)
		})
	}
}

func TestDNS_Integration_CacheEviction(t *testing.T) {
	dnsHandler := createTestDNS()

	msg1 := &dns.Msg{}
	msg1.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA},
	}

	msg2 := &dns.Msg{}
	msg2.Question = []dns.Question{
		{Name: "facebook.com.", Qtype: dns.TypeA},
	}

	resp1 := &dns.Msg{}
	resp1.SetReply(msg1)
	
	resp2 := &dns.Msg{}
	resp2.SetReply(msg2)

	dnsHandler.cacheMap.Store(msg1.String(), resp1)
	dnsHandler.cacheMap.Store(msg2.String(), resp2)

	cachedResp1, loaded1 := dnsHandler.cacheMap.Load(msg1.String())
	require.True(t, loaded1)
	assert.Equal(t, resp1, cachedResp1.(*dns.Msg))

	cachedResp2, loaded2 := dnsHandler.cacheMap.Load(msg2.String())
	require.True(t, loaded2)
	assert.Equal(t, resp2, cachedResp2.(*dns.Msg))

	dnsHandler.cacheMap.Delete(msg1.String())

	_, loaded1After := dnsHandler.cacheMap.Load(msg1.String())
	assert.False(t, loaded1After)

	cachedResp2After, loaded2After := dnsHandler.cacheMap.Load(msg2.String())
	require.True(t, loaded2After)
	assert.Equal(t, resp2, cachedResp2After.(*dns.Msg))
}

func TestDNS_Integration_ConcurrentBlockedDomainsUpdate(t *testing.T) {
	dnsHandler := createTestDNS()

	initialDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware1.com.", Recursive: false},
		{ID: 2, Domain: "malware2.com.", Recursive: false},
	}
	dnsHandler.updateBlockedDomainsMap(initialDomains)

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			msg := &dns.Msg{}
			msg.Question = []dns.Question{
				{Name: "malware1.com.", Qtype: dns.TypeA},
			}

			isBlocked := false
			for _, bds := range dnsHandler.blockedDomainsMap {
				for _, q := range msg.Question {
					for _, bd := range bds {
						if q.Name == bd.Domain {
							isBlocked = true
							break
						}
					}
				}
			}
			_ = isBlocked
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			newDomains := []blockeddomains.BlockedDomain{
				{ID: 3, Domain: "newmalware.com.", Recursive: false},
			}
			dnsHandler.updateBlockedDomainsMap(newDomains)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	<-done
	<-done

	assert.Contains(t, dnsHandler.blockedDomainsMap, "3")
	assert.NotContains(t, dnsHandler.blockedDomainsMap, "1")
	assert.NotContains(t, dnsHandler.blockedDomainsMap, "2")
}

func TestDNS_Integration_BlockedDomainsMapConsistency(t *testing.T) {
	dnsHandler := createTestDNS()

	domains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware.com.", Recursive: false},
		{ID: 1, Domain: "phishing.com.", Recursive: false},
		{ID: 2, Domain: ".ads.example.com", Recursive: true},
	}

	dnsHandler.updateBlockedDomainsMap(domains)

	assert.Len(t, dnsHandler.blockedDomainsMap, 2)
	
	group1, exists1 := dnsHandler.blockedDomainsMap["1"]
	require.True(t, exists1)
	assert.Len(t, group1, 2)
	
	domainNames := make([]string, len(group1))
	for i, bd := range group1 {
		domainNames[i] = bd.Domain
	}
	assert.Contains(t, domainNames, "malware.com.")
	assert.Contains(t, domainNames, "phishing.com.")
	
	group2, exists2 := dnsHandler.blockedDomainsMap["2"]
	require.True(t, exists2)
	assert.Len(t, group2, 1)
	assert.Equal(t, ".ads.example.com", group2[0].Domain)
	assert.True(t, group2[0].Recursive)
}

func TestDNS_Integration_MessageCacheKeyConsistency(t *testing.T) {
	dnsHandler := createTestDNS()

	msg1 := &dns.Msg{}
	msg1.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
	}

	msg2 := &dns.Msg{}
	msg2.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA, Qclass: dns.ClassINET},
	}

	key1 := msg1.String()
	key2 := msg2.String()
	assert.Equal(t, key1, key2)

	resp := &dns.Msg{}
	resp.SetReply(msg1)
	
	dnsHandler.cacheMap.Store(key1, resp)

	cachedResp, loaded := dnsHandler.cacheMap.Load(key2)
	require.True(t, loaded)
	assert.Equal(t, resp, cachedResp.(*dns.Msg))
}