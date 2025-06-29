package testutil

import (
	"time"
)

type TestBlockedDomain struct {
	ID        int
	Domain    string
	Recursive bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestCategory struct {
	Category string
}

type TestDomain struct {
	Domain    string
	UsedCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTestBlockedDomain(domain string, recursive bool) TestBlockedDomain {
	return TestBlockedDomain{
		ID:        1,
		Domain:    domain,
		Recursive: recursive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestCategory(name string) TestCategory {
	return TestCategory{
		Category: name,
	}
}

func NewTestDomain(domain string) TestDomain {
	return TestDomain{
		Domain:    domain,
		UsedCount: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateTestBlockedDomains() []TestBlockedDomain {
	return []TestBlockedDomain{
		NewTestBlockedDomain("malware.com", false),
		NewTestBlockedDomain(".ads.example.com", true),
		NewTestBlockedDomain("phishing.net", false),
	}
}

func CreateTestCategories() []TestCategory {
	return []TestCategory{
		NewTestCategory("malware"),
		NewTestCategory("ads"),
		NewTestCategory("social"),
	}
}

func CreateTestDomains() []TestDomain {
	return []TestDomain{
		NewTestDomain("google.com"),
		NewTestDomain("facebook.com"),
		NewTestDomain("youtube.com"),
	}
}