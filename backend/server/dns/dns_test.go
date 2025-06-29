package dns

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/internal/ai"
	"github.com/orion-tec/oriondns/internal/blockeddomains"
	"github.com/orion-tec/oriondns/internal/domains"
	"github.com/orion-tec/oriondns/internal/stats"
)

type MockAI struct {
	mock.Mock
}

func (m *MockAI) CategorizeURL(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

type MockBlockedDomains struct {
	mock.Mock
}

func (m *MockBlockedDomains) Insert(ctx context.Context, domain string, recursive bool) error {
	args := m.Called(ctx, domain, recursive)
	return args.Error(0)
}

func (m *MockBlockedDomains) GetAll(ctx context.Context) ([]blockeddomains.BlockedDomain, error) {
	args := m.Called(ctx)
	return args.Get(0).([]blockeddomains.BlockedDomain), args.Error(1)
}

type MockDomains struct {
	mock.Mock
}

func (m *MockDomains) Insert(ctx context.Context, domain string) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *MockDomains) GetAll(ctx context.Context) ([]domains.Domain, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domains.Domain), args.Error(1)
}

func (m *MockDomains) GetByDomain(ctx context.Context, domain string) (*domains.Domain, error) {
	args := m.Called(ctx, domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domains.Domain), args.Error(1)
}

func (m *MockDomains) GetDomainsWithoutCategory(ctx context.Context) ([]domains.Domain, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domains.Domain), args.Error(1)
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Insert(ctx context.Context, t time.Time, domain, domainType string) error {
	args := m.Called(ctx, t, domain, domainType)
	return args.Error(0)
}

func (m *MockStats) GetMostUsedDomains(ctx context.Context, from, to time.Time, categories []string, limit int) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, categories, limit)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStats) GetUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, domains []string) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, domains)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStats) GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, categories []string) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, categories)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStats) GetServerUsageByTimeRange(ctx context.Context, from, to time.Time, categories []string) ([]stats.ServerUsageByTimeRangeResponse, error) {
	args := m.Called(ctx, from, to, categories)
	return args.Get(0).([]stats.ServerUsageByTimeRangeResponse), args.Error(1)
}

func createTestDNS() *DNS {
	return &DNS{
		cacheMap:             sync.Map{},
		blockedDomainsMap:    make(map[string][]blockeddomains.BlockedDomain),
		blockedDomainsMutext: sync.Mutex{},
		blockedDomains:       &MockBlockedDomains{},
		domain:               &MockDomains{},
		stats:                &MockStats{},
		ai:                   &MockAI{},
	}
}

func TestDNS_updateBlockedDomainsMap(t *testing.T) {
	dnsHandler := createTestDNS()
	
	testDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware.com", Recursive: false},
		{ID: 2, Domain: ".ads.example.com", Recursive: true},
	}

	dnsHandler.updateBlockedDomainsMap(testDomains)

	assert.Len(t, dnsHandler.blockedDomainsMap, 2)
	assert.Contains(t, dnsHandler.blockedDomainsMap, "1")
	assert.Contains(t, dnsHandler.blockedDomainsMap, "2")
	
	assert.Equal(t, "malware.com", dnsHandler.blockedDomainsMap["1"][0].Domain)
	assert.False(t, dnsHandler.blockedDomainsMap["1"][0].Recursive)
	
	assert.Equal(t, ".ads.example.com", dnsHandler.blockedDomainsMap["2"][0].Domain)
	assert.True(t, dnsHandler.blockedDomainsMap["2"][0].Recursive)
}

func TestDNS_updateBlockedDomainsMap_ClearsExisting(t *testing.T) {
	dnsHandler := createTestDNS()
	
	oldDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "old.com", Recursive: false},
	}
	dnsHandler.updateBlockedDomainsMap(oldDomains)
	assert.Len(t, dnsHandler.blockedDomainsMap, 1)

	newDomains := []blockeddomains.BlockedDomain{
		{ID: 2, Domain: "new.com", Recursive: false},
	}
	dnsHandler.updateBlockedDomainsMap(newDomains)
	
	assert.Len(t, dnsHandler.blockedDomainsMap, 1)
	assert.Contains(t, dnsHandler.blockedDomainsMap, "2")
	assert.NotContains(t, dnsHandler.blockedDomainsMap, "1")
}

func TestDNS_DomainBlocking_ExactMatch(t *testing.T) {
	dnsHandler := createTestDNS()
	
	blockedDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware.com.", Recursive: false},
	}
	dnsHandler.updateBlockedDomainsMap(blockedDomains)

	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "malware.com.", Qtype: dns.TypeA},
	}

	isBlocked := false
	for _, bds := range dnsHandler.blockedDomainsMap {
		for _, q := range msg.Question {
			for _, bd := range bds {
				if bd.Recursive && strings.HasSuffix(q.Name, bd.Domain) {
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

	assert.True(t, isBlocked)
}

func TestDNS_DomainBlocking_RecursiveMatch(t *testing.T) {
	dnsHandler := createTestDNS()
	
	blockedDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: ".ads.example.com", Recursive: true},
	}
	dnsHandler.updateBlockedDomainsMap(blockedDomains)

	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "tracker.ads.example.com", Qtype: dns.TypeA},
	}

	isBlocked := false
	for _, bds := range dnsHandler.blockedDomainsMap {
		for _, q := range msg.Question {
			for _, bd := range bds {
				if bd.Recursive && strings.HasSuffix(q.Name, bd.Domain) {
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

	assert.True(t, isBlocked)
}

func TestDNS_DomainBlocking_NoMatch(t *testing.T) {
	dnsHandler := createTestDNS()
	
	blockedDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: "malware.com.", Recursive: false},
	}
	dnsHandler.updateBlockedDomainsMap(blockedDomains)

	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA},
	}

	isBlocked := false
	for _, bds := range dnsHandler.blockedDomainsMap {
		for _, q := range msg.Question {
			for _, bd := range bds {
				if bd.Recursive && strings.HasSuffix(q.Name, bd.Domain) {
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

	assert.False(t, isBlocked)
}

func TestDNS_DomainBlocking_RecursiveNoMatch(t *testing.T) {
	dnsHandler := createTestDNS()
	
	blockedDomains := []blockeddomains.BlockedDomain{
		{ID: 1, Domain: ".ads.example.com", Recursive: true},
	}
	dnsHandler.updateBlockedDomainsMap(blockedDomains)

	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA},
	}

	isBlocked := false
	for _, bds := range dnsHandler.blockedDomainsMap {
		for _, q := range msg.Question {
			for _, bd := range bds {
				if bd.Recursive && strings.HasSuffix(q.Name, bd.Domain) {
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

	assert.False(t, isBlocked)
}

func TestDNS_Cache_Store_Load(t *testing.T) {
	dnsHandler := createTestDNS()
	
	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA},
	}

	respMsg := &dns.Msg{}
	respMsg.SetReply(msg)

	dnsHandler.cacheMap.Store(msg.String(), respMsg)

	cachedResp, loaded := dnsHandler.cacheMap.Load(msg.String())
	require.True(t, loaded)
	assert.Equal(t, respMsg, cachedResp.(*dns.Msg))
}

func TestDNS_Cache_NotFound(t *testing.T) {
	dnsHandler := createTestDNS()
	
	msg := &dns.Msg{}
	msg.Question = []dns.Question{
		{Name: "google.com.", Qtype: dns.TypeA},
	}

	_, loaded := dnsHandler.cacheMap.Load(msg.String())
	assert.False(t, loaded)
}