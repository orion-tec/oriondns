package domains

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/testutil"
)

func TestDomainsDB_Insert(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()
	domain := "google.com"

	err := domainsDB.Insert(ctx, domain)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM domains WHERE domain = $1", domain).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDomainsDB_Insert_Duplicate_UpdatesUsedCount(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()
	domain := "google.com"

	err := domainsDB.Insert(ctx, domain)
	require.NoError(t, err)

	err = domainsDB.Insert(ctx, domain)
	require.NoError(t, err)

	var usedCount int
	err = pool.QueryRow(ctx, "SELECT used_count FROM domains WHERE domain = $1", domain).Scan(&usedCount)
	require.NoError(t, err)
	assert.Equal(t, 2, usedCount)
}

func TestDomainsDB_GetAll(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()

	domains := []string{"google.com", "facebook.com", "youtube.com"}
	for _, domain := range domains {
		err := domainsDB.Insert(ctx, domain)
		require.NoError(t, err)
	}

	results, err := domainsDB.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, results, 3)

	resultDomains := make([]string, len(results))
	for i, result := range results {
		resultDomains[i] = result.Domain
	}

	assert.Contains(t, resultDomains, "google.com")
	assert.Contains(t, resultDomains, "facebook.com")
	assert.Contains(t, resultDomains, "youtube.com")
}

func TestDomainsDB_GetByDomain(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()
	domain := "google.com"

	err := domainsDB.Insert(ctx, domain)
	require.NoError(t, err)

	result, err := domainsDB.GetByDomain(ctx, domain)
	require.NoError(t, err)
	assert.Equal(t, domain, result.Domain)
}

func TestDomainsDB_GetByDomain_NotFound(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()
	domain := "nonexistent.com"

	result, err := domainsDB.GetByDomain(ctx, domain)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDomainsDB_GetDomainsWithoutCategory(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	domainsDB := New(database)

	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		INSERT INTO domains (domain, used_count) VALUES 
		('google.com', 100),
		('facebook.com', 50),
		('categorized.com', 75)
	`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO domain_categories (domain, category) VALUES 
		('categorized.com', 'social')
	`)
	require.NoError(t, err)

	results, err := domainsDB.GetDomainsWithoutCategory(ctx)
	require.NoError(t, err)
	require.Len(t, results, 2)

	domainNames := make([]string, len(results))
	for i, result := range results {
		domainNames[i] = result.Domain
	}

	assert.Contains(t, domainNames, "google.com")
	assert.Contains(t, domainNames, "facebook.com")
	assert.NotContains(t, domainNames, "categorized.com")

	assert.Equal(t, "google.com", results[0].Domain)
	assert.Equal(t, "facebook.com", results[1].Domain)
}