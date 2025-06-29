package testutil

import (
	"time"

	"github.com/orion-tec/oriondns/internal/blockeddomains"
	"github.com/orion-tec/oriondns/internal/categories"
	"github.com/orion-tec/oriondns/internal/domains"
	"github.com/orion-tec/oriondns/internal/stats"
)

func NewTestBlockedDomain(domain string, recursive bool) blockeddomains.BlockedDomain {
	return blockeddomains.BlockedDomain{
		ID:        1,
		Domain:    domain,
		Recursive: recursive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestCategory(name string) categories.Category {
	return categories.Category{
		ID:        1,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestDomain(domain string, categoryID int) domains.Domain {
	return domains.Domain{
		ID:         1,
		Domain:     domain,
		CategoryID: categoryID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewTestStat(domain string, blocked bool) stats.Stat {
	return stats.Stat{
		ID:        1,
		Domain:    domain,
		Blocked:   blocked,
		CreatedAt: time.Now(),
	}
}

func CreateTestBlockedDomains() []blockeddomains.BlockedDomain {
	return []blockeddomains.BlockedDomain{
		NewTestBlockedDomain("malware.com", false),
		NewTestBlockedDomain(".ads.example.com", true),
		NewTestBlockedDomain("phishing.net", false),
	}
}

func CreateTestCategories() []categories.Category {
	return []categories.Category{
		NewTestCategory("malware"),
		NewTestCategory("ads"),
		NewTestCategory("social"),
	}
}

func CreateTestDomains() []domains.Domain {
	return []domains.Domain{
		NewTestDomain("google.com", 1),
		NewTestDomain("facebook.com", 3),
		NewTestDomain("youtube.com", 3),
	}
}

func CreateTestStats() []stats.Stat {
	return []stats.Stat{
		NewTestStat("google.com", false),
		NewTestStat("malware.com", true),
		NewTestStat("facebook.com", false),
	}
}