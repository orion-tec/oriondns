package stats

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/domains"
	"github.com/orion-tec/oriondns/internal/testutil"
)

func TestStatsDB_Insert(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)

	ctx := context.Background()
	testTime := time.Now()
	domain := "example.com"
	domainType := "A"

	err := statsDB.Insert(ctx, testTime, domain, domainType)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM stats_aggregated WHERE domain = $1", domain).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestStatsDB_Insert_Aggregation(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)

	ctx := context.Background()
	baseTime := time.Date(2023, 1, 1, 12, 5, 0, 0, time.UTC)
	domain := "example.com"
	domainType := "A"

	err := statsDB.Insert(ctx, baseTime, domain, domainType)
	require.NoError(t, err)

	sameMinuteTime := time.Date(2023, 1, 1, 12, 7, 30, 0, time.UTC)
	err = statsDB.Insert(ctx, sameMinuteTime, domain, domainType)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT count FROM stats_aggregated WHERE domain = $1", domain).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestStatsDB_GetMostUsedDomains(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search'),
		('facebook.com', 'social')
	`)
	require.NoError(t, err)

	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type) VALUES 
		($1, 'google.com', 10, 'A'),
		($1, 'facebook.com', 5, 'A')
	`, baseTime)
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	results, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{}, 10)
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.Equal(t, "google.com", results[0].Domain)
	assert.Equal(t, int64(10), results[0].Count)
	assert.Equal(t, "facebook.com", results[1].Domain)
	assert.Equal(t, int64(5), results[1].Count)
}

func TestStatsDB_GetMostUsedDomains_WithCategories(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search'),
		('facebook.com', 'social')
	`)
	require.NoError(t, err)

	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type) VALUES 
		($1, 'google.com', 10, 'A'),
		($1, 'facebook.com', 5, 'A')
	`, baseTime)
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	results, err := statsDB.GetMostUsedDomains(ctx, from, to, []string{"search"}, 10)
	require.NoError(t, err)

	require.Len(t, results, 1)
	assert.Equal(t, "google.com", results[0].Domain)
	assert.Equal(t, int64(10), results[0].Count)
}

func TestStatsDB_GetServerUsageByTimeRange(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('google.com', 'search'),
		('facebook.com', 'social')
	`)
	require.NoError(t, err)

	baseTime1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	baseTime2 := time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type) VALUES 
		($1, 'google.com', 10, 'A'),
		($1, 'facebook.com', 5, 'A'),
		($2, 'google.com', 8, 'A')
	`, baseTime1, baseTime2)
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	results, err := statsDB.GetServerUsageByTimeRange(ctx, from, to, []string{})
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.Equal(t, int64(15), results[0].Count)
	assert.Equal(t, int64(8), results[1].Count)
}

func TestStatsDB_GetUsedDomainsByTimeAggregation(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.NewWithPool(pool)
	statsDB := New(database)
	domainsDB := domains.New(database)

	ctx := context.Background()
	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	err := domainsDB.Insert(ctx, "google.com")
	require.NoError(t, err)
	err = domainsDB.Insert(ctx, "facebook.com")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type) VALUES 
		($1, 'google.com', 10, 'A'),
		($1, 'facebook.com', 5, 'A')
	`, baseTime)
	require.NoError(t, err)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	domains := []string{"google.com", "facebook.com"}

	results, err := statsDB.GetUsedDomainsByTimeAggregation(ctx, from, to, domains)
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.Equal(t, "google.com", results[0].Domain)
	assert.Equal(t, int64(10), results[0].Count)
	assert.Equal(t, "facebook.com", results[1].Domain)
	assert.Equal(t, int64(5), results[1].Count)
}
