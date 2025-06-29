package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/blockeddomains"
	"github.com/orion-tec/oriondns/internal/domains"
	"github.com/orion-tec/oriondns/internal/stats"
	"github.com/orion-tec/oriondns/internal/testutil"
)

func TestIntegration_FullWorkflow(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := stats.New(database)
	domainsDB := domains.New(database)
	blockedDomainsDB := blockeddomains.New(database)

	ctx := context.Background()

	err := blockedDomainsDB.Insert(ctx, "malware.com", false)
	require.NoError(t, err)

	err = blockedDomainsDB.Insert(ctx, ".ads.example.com", true)
	require.NoError(t, err)

	blockedDomainsList, err := blockedDomainsDB.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, blockedDomainsList, 2)

	err = domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)

	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)

	domainsList, err := domainsDB.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, domainsList, 2)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search'),
		('facebook.com', 'social')
	`)
	require.NoError(t, err)

	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	err = statsDB.Insert(ctx, testTime, "google.com", "A")
	require.NoError(t, err)

	err = statsDB.Insert(ctx, testTime, "facebook.com", "A")
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	mostUsed, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{}, 10)
	require.NoError(t, err)
	assert.Len(t, mostUsed, 2)

	usage, err := statsDB.GetServerUsageByTimeRange(ctx, from, to, []string{})
	require.NoError(t, err)
	assert.Len(t, usage, 1)
	assert.Equal(t, int64(2), usage[0].Count)
}

func TestIntegration_StatsAggregation(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := stats.New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search')
	`)
	require.NoError(t, err)

	baseTime := time.Date(2023, 1, 1, 12, 5, 0, 0, time.UTC)

	err = statsDB.Insert(ctx, baseTime, "google.com", "A")
	require.NoError(t, err)

	sameWindowTime := time.Date(2023, 1, 1, 12, 7, 30, 0, time.UTC)
	err = statsDB.Insert(ctx, sameWindowTime, "google.com", "A")
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	results, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{}, 10)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "google.com", results[0].Domain)
	assert.Equal(t, int64(2), results[0].Count)
}

func TestIntegration_CategoryFiltering(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := stats.New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "youtube.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search'),
		('facebook.com', 'social'),
		('youtube.com', 'social')
	`)
	require.NoError(t, err)

	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type) VALUES 
		($1, 'google.com', 10, 'A'),
		($1, 'facebook.com', 5, 'A'),
		($1, 'youtube.com', 3, 'A')
	`, baseTime)
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	searchResults, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{"search"}, 10)
	require.NoError(t, err)
	require.Len(t, searchResults, 1)
	assert.Equal(t, "google.com", searchResults[0].Domain)

	socialResults, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{"social"}, 10)
	require.NoError(t, err)
	require.Len(t, socialResults, 2)
	assert.Equal(t, "facebook.com", socialResults[0].Domain)
	assert.Equal(t, "youtube.com", socialResults[1].Domain)

	allResults, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{}, 10)
	require.NoError(t, err)
	require.Len(t, allResults, 3)
}
